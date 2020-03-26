package net

import (
	"bufio"
	"bytes"
	"encoding/binary"
	"github.com/tpetrychyn/isaac"
	"github.com/tpetrychyn/osrs-cache-parser/pkg/cachestore"
	utils2 "github.com/tpetrychyn/osrs-cache-parser/pkg/utils"
	"github.com/tpetrychyn/rsps-comm-test/pkg/utils"
	"log"
	"math/big"
	"math/rand"
	"net"
	"unsafe"
)

type LoginHandler struct {
	CacheStore *cachestore.Store
}

func (l *LoginHandler) HandleRequest(connection net.Conn, reader *bufio.Reader) (*isaac.ISAAC, *isaac.ISAAC) {
	connection.Write([]byte{0}) // proceed
	serverSeed := rand.Uint64()
	binary.Write(connection, binary.BigEndian, serverSeed)

	var loginRequest LoginRequest
	binary.Read(reader, binary.BigEndian, &loginRequest)
	log.Printf("loginReq: %+v", loginRequest)

	secureBuf := make([]byte, loginRequest.RsaBlockSize)
	binary.Read(reader, binary.BigEndian, &secureBuf)

	// I got these values by running rsmod in debug and capturing the base 10 values of rsaExponent and rsaModulus
	//e, _ := new(big.Int).SetString("647938820811174564501994386359961232566329715148620231085023210484768344331719057997549872356145944358302879928019776209939522512099493887127188885541797604448573171793150698254030927402524326168972623509694505212035910746949393405712308360102892779655139499817708596651580539817342635450991423595932014585588510405014175673731955960405352231021121560962415148107396591983917910196047875771120213776287015906226300640663223264593459002543216867384883826063220469591493399234965151408618988915160375744659408215352143147359834065449825250330535673388894408931234365732741021749839941491681984688983902956775387923223", 10)
	//m, _ := new(big.Int).SetString("18845951212332454083776854822768074502513669601639934607217255540895979767143769429570326954651016470938067579197871264181052173601410791958445525182581678695148072893875671021101345647440460126736516835548197274786392923021699589121012703818192613010100627154659830597006336146207911878070145791564189258473567299280325151901023854730803093639677760621428074422539345082901833520173983428905869285592015645143627391159674409357147242143770639015644409263814106735096728020275635955800783094923547218134070829959486241186122853067651743831358738523668508661897945035138609772757816654638520047949322016581053188813993", 10)

	// new values for runelite client
	// base16 modulus: a04f258b9032e39f22abc6f43535e8b5246ca099e5dc921affaf9b42f0d67cf612577dda55872e6b9b5d89ae98f7626957005b9482f614ca46620902dd425d61b38f072ec775f6b7ce613402766a16547455b66234c45ebefae8a4e8c67cf3e6c1714ba8050fe0d54db1c2fa008af8c491260b4312f35f73daf5de4c5dc0cc22b22bfaf3975065e5d3f2581120c391084d07881be9989235db734ca4931409e30a0f4f5022133e98d256060f389f2a7a372ed39d33dbaeda00587f9294e7e56aa90f80a4c6a772154232df72da90b84494b43445d3d152199fc536d0b9a299aba1d78f90acbc2e903b215b9d6daa4b7f25f80466688866b8ad4f4a352dbf8281
	e, _ := new(big.Int).SetString("282697066719888242017392573796327827203250422092339778942246726904685141674102211944168166603103690254289820189637133824150980079916270829290747243173908137730751884073623045839072138462741019014014222172539359142330133691328550772237428640325003259043534090004215916089722033291883150000561606013315364867654848270661252326340240989413251482705148170847331690180617565715393327309533722296359694530980252728971128627603013425313795233077084468655900617031732908931891072315486762876530421501471838438285474695139733971808492171119735425077542515219377312979814415383439261706977041945372093670461625336103437351425", 10)
	m, _ := new(big.Int).SetString("20237157467636609193985644029371858887405158834151471428222854987604970103654436552905460551248068321349417745241123800582613633530827571096917205981309030718143403852029528733102425711013280352945330506304436897991141421877224720873975271218983876120083117047084979238418473944129050793650252291965755398505469966275781531714987251487836316163399886603193140269980225283653824805239448041494114677754531878174463240806201389143733032911002273397545973151941058462805580061242840962439870213807641884588480643048115234980996159790336547429096279012777892681833696187792172729230085683599814216501083676227645430792833", 10)


	encrypted := big.NewInt(0)
	encrypted.SetBytes(secureBuf[:])
	var rs big.Int
	rs = *rs.Exp(encrypted, e, m)
	rsaBuffer := bytes.NewBuffer(rs.Bytes())

	var successfulDecrypt byte
	binary.Read(rsaBuffer, binary.BigEndian, &successfulDecrypt)
	log.Printf("successfulDecrypt %+v", successfulDecrypt)
	if successfulDecrypt != 1 {

		connection.Write([]byte{10}) // bad session id
		return nil, nil
	}

	xteaKeys := make([]int32, 4)
	binary.Read(rsaBuffer, binary.BigEndian, &xteaKeys)

	var reportedSeed uint64
	binary.Read(rsaBuffer, binary.BigEndian, &reportedSeed)

	if serverSeed != reportedSeed {
		log.Printf("serverSeed %+v", serverSeed)
		connection.Write([]byte{10})
		return nil, nil
	}

	var authType byte
	binary.Read(rsaBuffer, binary.BigEndian, &authType)
	// TODO: Authtype can be anything valid
	//if authType != 0 {
	//	log.Printf("authType %+v", authType)
	//	connection.Write([]byte{10})
	//	return nil, nil
	//}

	var skip = make([]byte, 5) // authcode 2, unkown 1, another skip
	binary.Read(rsaBuffer, binary.BigEndian, &skip)

	password := utils.ReadString(rsaBuffer)

	log.Printf("password: %+v", password)

	xteaCipher, err := utils2.XteaKeyFromIntArray(xteaKeys)
	if err != nil {
		connection.Write([]byte{10})
		return nil, nil
	}

	xteaEncryptedBytes := make([]byte, reader.Buffered())
	binary.Read(reader, binary.BigEndian, &xteaEncryptedBytes)

	xteaBytes := utils2.XteaDecrypt(xteaCipher, xteaEncryptedBytes)
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

	return encryptor, decryptor
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