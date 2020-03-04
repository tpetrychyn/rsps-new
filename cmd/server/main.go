package main

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
	"rsps-comm-test/pkg/packet/outgoing"
	"rsps-comm-test/pkg/utils"
	"unsafe"
)

var store = cachestore.NewStore()

func main() {

	listener, err := net.Listen("tcp", ":43594")
	if err != nil {
		log.Fatal(err)
		return
	}

	defer listener.Close()

	log.Printf("listening on 43594")
	for {
		connection, err := listener.Accept()
		if err != nil {
			continue
		}

		reader := bufio.NewReader(connection)

		var requestType byte
		binary.Read(reader, binary.BigEndian, &requestType)

		log.Printf("request type %d", requestType)

		if requestType == 15 {
			var gameVersion int32
			binary.Read(reader, binary.BigEndian, &gameVersion)
			if gameVersion == 181 {
				connection.Write([]byte{0})
			}
			go func() {
				for handleJS5Requests(connection, reader) {
				}
			}()
		}

		if requestType == 14 {
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
				continue
			}

			xteaKeys := make([]int32, 4)
			binary.Read(rsaBuffer, binary.BigEndian, &xteaKeys)

			var reportedSeed uint64
			binary.Read(rsaBuffer, binary.BigEndian, &reportedSeed)

			if serverSeed != reportedSeed {
				connection.Write([]byte{10})
				break
			}

			var authType byte
			binary.Read(rsaBuffer, binary.BigEndian, &authType)
			if authType != 0 {
				panic("unknown authType " + string(authType))
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
				panic(err)
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
			utils.ReadString(xteaBuffer) // some sort of hash. client hashsum? unique to pc?

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

			for k, v := range store.Indexes {
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
			log.Printf("remaining %+v", remaining)

			inC := isaac.ISAAC{}

			uXteaKeys := *(*[]uint32)(unsafe.Pointer(&xteaKeys))
			inC.Generate(uXteaKeys)
			//decryptor := &inC

			for i := 0; i < 4; i++ {
				uXteaKeys[i] += 50
			}
			outC := isaac.ISAAC{}
			outC.Generate(uXteaKeys)
			encryptor := &outC

			// TODO: START LoginEncoder This section needs to be abstracted!!
			outBuffer := new(bytes.Buffer)

			outBuffer.Write([]byte{2, 13, 0, 0, 0, 0, 0, 0, 1, 0, 1, 1}) // successful login packet

			regionBuffer := make([]byte, 0, 4710) // rebuildLoginEncoder gpi

			opcode := byte(0 + (encryptor.Rand() & 0xFF))
			regionBuffer = append(regionBuffer, opcode)

			payload := []byte{18, 102, 12, 15, 54, -92+256}
			regionBuffer = append(regionBuffer, payload...)

			pad := 4610 - len(regionBuffer)
			regionBuffer = append(regionBuffer, make([]byte, pad)...)

			regionBuffer = append(regionBuffer, []byte{0,0,0,0,0,0,0,0}...)

			outBuffer.Write(regionBuffer)

			r := &outgoing.RebuildNormalPacket{
				X: 386,
				Z: 437,
			}
			regionPacket := r.Write(nil)
			outBuffer.Write(regionPacket)

			log.Printf("length %d outbytes %+v", len(outBuffer.Bytes()), outBuffer.Bytes())
			connection.Write(outBuffer.Bytes())

			// TODO: END LoginEncoder

			log.Printf("outbuffer %+v", outBuffer.Bytes())

			b, err := reader.ReadByte()
			if err != nil {
				log.Printf("err %s", err.Error())
			}
			log.Printf("b %+v", b)
		}
	}
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

func handleJS5Requests(connection net.Conn, reader *bufio.Reader) bool {
	var opcode byte
	err := binary.Read(reader, binary.BigEndian, &opcode)
	if err != nil {
		connection.Close()
		return false
	}

	if opcode == 3 || opcode == 2 || opcode == 6 { // skip loading screen
		reader.Discard(3)
		binary.Read(reader, binary.BigEndian, &opcode)
	}

	if opcode == 0 || opcode == 1 {
		var index uint8
		var archive uint16

		binary.Read(reader, binary.BigEndian, &index)
		binary.Read(reader, binary.BigEndian, &archive)

		if index == 255 {
			data := new(bytes.Buffer)
			if archive == 255 {
				dataBuffer := new(bytes.Buffer)
				for _, v := range store.Indexes {
					binary.Write(dataBuffer, binary.BigEndian, v.Crc)
					binary.Write(dataBuffer, binary.BigEndian, v.Revision)
				}
				data.Write([]byte{0}) // no compression
				binary.Write(data, binary.BigEndian, int32(dataBuffer.Len()))
				binary.Write(data, binary.BigEndian, dataBuffer.Bytes())
			} else {
				r := store.ReadIndex(int(archive))
				data.Write(r)
			}

			outBuffer := new(bytes.Buffer)
			binary.Write(outBuffer, binary.BigEndian, index)
			binary.Write(outBuffer, binary.BigEndian, archive)
			writeIndex := 3 // 1 byte for index, 2 bytes for archive so we've written 3 so far
			for _, v := range data.Bytes() {
				if writeIndex%512 == 0 {
					binary.Write(outBuffer, binary.BigEndian, byte(255))
					writeIndex++
				}
				binary.Write(outBuffer, binary.BigEndian, v)
				writeIndex++
			}

			connection.Write(outBuffer.Bytes())
		} else {
			i := store.FindIndex(int(index))
			a := i.Archives[archive]
			data := store.LoadArchive(a)

			if a == nil || data == nil {
				log.Printf("nil data index %d", index)
			}

			if data != nil {
				compression := data[0]
				length := int(data[1])<<24 | (int(data[2])&0xFF)<<16 | (int(data[3])&0xFF)<<8 | int(data[4])&0xFF
				expectedLength := length
				if compression == 0 {
					expectedLength += 5
				} else {
					expectedLength += 9
				}
				if len(data)-expectedLength == 2 {
					data = data[:len(data)-2]
				}

				outBuffer := new(bytes.Buffer)
				binary.Write(outBuffer, binary.BigEndian, index)
				binary.Write(outBuffer, binary.BigEndian, archive)
				writeIndex := 3 // 1 byte for index, 2 bytes for archive so we've written 3 so far
				for _, v := range data {
					if writeIndex%512 == 0 {
						binary.Write(outBuffer, binary.BigEndian, byte(255))
						writeIndex++
					}
					binary.Write(outBuffer, binary.BigEndian, v)
					writeIndex++
				}
				connection.Write(outBuffer.Bytes())
			}
		}
	}

	return true
}
