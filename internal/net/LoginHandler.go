package net

import (
	"bufio"
	"bytes"
	"encoding/binary"
	"github.com/gtank/isaac"
	"golang.org/x/crypto/xtea"
	"log"
	"math/big"
	"math/rand"
	"net"
	"osrs-cache-parser/pkg/cachestore"
	"rsps-comm-test/internal/game"
	"rsps-comm-test/pkg/models"
	"rsps-comm-test/pkg/packet/outgoing"
	"rsps-comm-test/pkg/utils"
	"unsafe"
)

type LoginHandler struct {
	CacheStore *cachestore.Store
}

func (l *LoginHandler) HandleRequest(connection net.Conn, reader *bufio.Reader) *Client {
	connection.Write([]byte{0}) // proceed
	serverSeed := rand.Uint64()
	binary.Write(connection, binary.BigEndian, serverSeed)

	var loginRequest LoginRequest
	binary.Read(reader, binary.BigEndian, &loginRequest)
	log.Printf("loginReq: %+v", loginRequest)

	secureBuf := make([]byte, loginRequest.RsaBlockSize)
	binary.Read(reader, binary.BigEndian, &secureBuf)

	// I got these values by running rsmod in debug and capturing the base 10 values of rsaExponent and rsaModulus
	e, _ := new(big.Int).SetString("647938820811174564501994386359961232566329715148620231085023210484768344331719057997549872356145944358302879928019776209939522512099493887127188885541797604448573171793150698254030927402524326168972623509694505212035910746949393405712308360102892779655139499817708596651580539817342635450991423595932014585588510405014175673731955960405352231021121560962415148107396591983917910196047875771120213776287015906226300640663223264593459002543216867384883826063220469591493399234965151408618988915160375744659408215352143147359834065449825250330535673388894408931234365732741021749839941491681984688983902956775387923223", 10)
	m, _ := new(big.Int).SetString("18845951212332454083776854822768074502513669601639934607217255540895979767143769429570326954651016470938067579197871264181052173601410791958445525182581678695148072893875671021101345647440460126736516835548197274786392923021699589121012703818192613010100627154659830597006336146207911878070145791564189258473567299280325151901023854730803093639677760621428074422539345082901833520173983428905869285592015645143627391159674409357147242143770639015644409263814106735096728020275635955800783094923547218134070829959486241186122853067651743831358738523668508661897945035138609772757816654638520047949322016581053188813993", 10)

	encrypted := big.NewInt(0)
	encrypted.SetBytes(secureBuf[:])
	var rs big.Int
	rs = *rs.Exp(encrypted, e, m)
	rsaBuffer := bytes.NewBuffer(rs.Bytes())

	var successfulDecrypt byte
	binary.Read(rsaBuffer, binary.BigEndian, &successfulDecrypt)
	if successfulDecrypt != 1 {
		connection.Write([]byte{10}) // bad session id
		return nil
	}

	xteaKeys := make([]int32, 4)
	binary.Read(rsaBuffer, binary.BigEndian, &xteaKeys)

	var reportedSeed uint64
	binary.Read(rsaBuffer, binary.BigEndian, &reportedSeed)

	if serverSeed != reportedSeed {
		connection.Write([]byte{10})
		return nil
	}

	var authType byte
	binary.Read(rsaBuffer, binary.BigEndian, &authType)
	if authType != 0 {
		connection.Write([]byte{10})
		return nil
	}

	var skip = make([]byte, 5) // authcode 2, unkown 1, another skip
	binary.Read(rsaBuffer, binary.BigEndian, &skip)

	password := utils.ReadString(rsaBuffer)

	log.Printf("password: %+v", password)

	xteaKey := make([]byte, 16)
	for i := 0; i < len(xteaKeys); i++ {
		j := i << 2
		xteaKey[j] = byte(xteaKeys[i] >> 24)
		xteaKey[j+1] = byte(xteaKeys[i] >> 16)
		xteaKey[j+2] = byte(xteaKeys[i] >> 8)
		xteaKey[j+3] = byte(xteaKeys[i])
	}

	xteaCipher, err := xtea.NewCipher(xteaKey)
	if err != nil {
		connection.Write([]byte{10})
		return nil
	}

	xteaEncryptedBytes := make([]byte, reader.Buffered())
	binary.Read(reader, binary.BigEndian, &xteaEncryptedBytes)

	xteaBytes := utils.XteaDecrypt(xteaCipher, xteaEncryptedBytes)
	xteaBuffer := bytes.NewBuffer(xteaBytes)

	username := utils.ReadString(xteaBuffer)
	log.Printf("username: %s", username)

	var clientDetails ClientDetails
	binary.Read(xteaBuffer, binary.BigEndian, &clientDetails)
	log.Printf("clientDetails %+v", clientDetails)

	// big skips
	binary.Read(xteaBuffer, binary.BigEndian, make([]byte, 24)) // skip random.dat
	utils.ReadString(xteaBuffer)                                // some sort of hash. client hashsum? unique to pc?

	binary.Read(xteaBuffer, binary.BigEndian, new(int32)) // some int

	// platformInfo
	binary.Read(xteaBuffer, binary.BigEndian, make([]byte, 18)) // 3 bytes, short, 5 bytes, short, byte, medium, short
	utils.ReadString(xteaBuffer)
	utils.ReadString(xteaBuffer)
	utils.ReadString(xteaBuffer)
	utils.ReadString(xteaBuffer)
	binary.Read(xteaBuffer, binary.BigEndian, make([]byte, 3)) // byte, short
	utils.ReadString(xteaBuffer)
	utils.ReadString(xteaBuffer)
	binary.Read(xteaBuffer, binary.BigEndian, make([]byte, 18)) // 2 bytes, 3 ints, int
	utils.ReadString(xteaBuffer)
	binary.Read(xteaBuffer, binary.BigEndian, make([]byte, 12)) // unknown 3 bytes that rsmod skips??

	// should be at Client.java:2899 - var31.packetBuffer.writeInt(GrandExchangeEvent.archive0.hash);

	for k, v := range l.CacheStore.Indexes {
		var crc uint32
		binary.Read(xteaBuffer, binary.BigEndian, &crc)
		if k == 16 || k == 20 { // 20 is being read different each time..
			continue // crc always 0
		}
		if crc != v.Crc {
			log.Printf("crc mismatch on index %d, read %+v, expected %+v", v.Id, crc, v.Crc)
			connection.Write([]byte{6}) // revision mismatch
			break
		}
	}

	remaining := make([]byte, xteaBuffer.Len())
	binary.Read(xteaBuffer, binary.BigEndian, &remaining)

	inC := isaac.ISAAC{}

	uXteaKeys := *(*[]uint32)(unsafe.Pointer(&xteaKeys))
	inC.Generate(uXteaKeys)
	decryptor := &inC

	for i := 0; i < 4; i++ {
		uXteaKeys[i] += 50
	}
	outC := isaac.ISAAC{}
	outC.Generate(uXteaKeys)
	encryptor := &outC

	connection.Write([]byte{2, 13, 0, 0, 0, 0, 0, 0, 1, 0, 1, 1})

	player := game.NewPlayer()
	client := NewClient(connection, encryptor, decryptor, player)
	player.OutgoingQueue = client.downstreamQueue

	player.Actor.Movement.Position = &models.Position{
		X:      3094,
		Z:      3497,
	}

	client.EnqueueOutgoing(&outgoing.RebuildLoginPacket{Position:player.Actor.Movement.Position})

	client.EnqueueOutgoing(&outgoing.IfOpenTopPacket{}) // main screen interface?

	client.EnqueueOutgoing(&outgoing.RebuildNormalPacket{Position:&models.Position{
		X:      player.Actor.Movement.Position.X >> 3,
		Z:      player.Actor.Movement.Position.Z >> 3,
	}})

	return client
}

type LoginRequest struct {
	Opcode       byte
	PacketSize   uint16
	Revision     uint32
	MagicValue   uint32
	ClientType   byte
	RsaBlockSize uint16
}

type ClientDetails struct {
	Flags  byte
	Width  uint16
	Height uint16
}