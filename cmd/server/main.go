package main

import (
	"bufio"
	"bytes"
	"encoding/binary"
	"golang.org/x/crypto/xtea"
	"log"
	"math/big"
	"net"
	"osrs-cache-parser/pkg/cachestore"
	"rsps-comm-test/pkg/utils"
)

var store = cachestore.NewStore()

func main() {

	listener, err := net.Listen("tcp", ":43594")
	if err != nil {
		log.Fatal(err)
		return
	}

	defer listener.Close()

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
			if gameVersion == 188 {
				connection.Write([]byte{0})
			}
			go func() {
				for handleJS5Requests(connection, reader) {
				}
			}()
		}

		if requestType == 14 {
			connection.Write([]byte{0}) // proceed
			serverSeed := uint64(0)
			binary.Write(connection, binary.BigEndian, serverSeed)

			var loginRequest LoginRequest
			binary.Read(reader, binary.BigEndian, &loginRequest)
			log.Printf("loginReq: %+v", loginRequest)

			secureBuf := make([]byte, loginRequest.RsaBlockSize)
			binary.Read(reader, binary.BigEndian, &secureBuf)

			// I got these values by running rsmod in debug and capturing the base 10 values of rsaExponent and rsaModulus
			e, _ := new(big.Int).SetString("148348469911079630378699587877554458887651856588926334927128326541191536140607970867465114083263709801890018879632658533196772725229915565152894411841475261416271281321582236105823416623786172712681530468166899094383930121268191893149065500392407752356169955528478757272431078241814609679978380952225376033562806433996728228109660312694769828344489814071015621494917392645701473740177821234108724899044366017790970020113106533060374669574194646107810847943186160131430203973239263865745553613943960329586454182123305166944314979813407192463286026550272667116112099818151107537342141656429241609974349952075315876533", 10)
			m, _ := new(big.Int).SetString("19170145713320489912151859072022254500259096926363991455668617611800271878918623428931512076795588840179453791702324170323331541120739007900995321847618625369808627488899464052761074041318183896353323609269960963555462408592083963811680638618219044683455866301621158049562441644108777919886027277997031442568533720342131888150818546453045613743049280955991189701825453826255427243627707953117170862683363604957024513356381405276233545512575360666085799937115726863049363483383120397664940232273940367851122809637163789793552566567683877420043802465899191878441559550339000760740180073343500614152433988960973619429001", 10)

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

			binary.Read(xteaBuffer, binary.BigEndian, make([]byte, 53)) // NetSocket.platformInfo

			var clientType byte
			binary.Read(xteaBuffer, binary.BigEndian, &clientType)
			if clientType != loginRequest.ClientType {
				log.Printf("client types did not match")
				connection.Write([]byte{10})
				continue
			}

			binary.Read(xteaBuffer, binary.BigEndian, new(uint32)) // 0
			// should be at Client.java:2899 - var31.packetBuffer.writeInt(GrandExchangeEvent.archive0.hash);

			connection.Write([]byte{10}) // bad session id, speeds up debugging
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
