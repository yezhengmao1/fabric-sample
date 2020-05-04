package cmd

import (
	"errors"
	"flag"
	"github.com/hyperledger/fabric/orderer/consensus/rbft/message"
	"go.dedis.ch/kyber"
	"go.dedis.ch/kyber/group/edwards25519"
	"go.dedis.ch/kyber/pairing/bn256"
	"go.dedis.ch/kyber/share"
	"go.dedis.ch/kyber/util/encoding"
	"log"
	"os"
	"strconv"
	"strings"
	"time"
)

func ReadConfig() *SharedConfig {
	port, _ := GetConfigurePort()
	id, _ := GetConfigureID()
	view, _ := GetConfigureView()
	table, _ := GetConfigureTable()

	t := make(map[message.Identify]string)
	for k, v := range table {
		t[message.Identify(k)] = v
	}
	// calc the fault num
	if len(t)%3 != 1 {
		log.Fatalf("[Config Error] the incorrent node num : %d, need 3f + 1", len(t))
		return nil
	}

	flag.Parse()
	return &SharedConfig{
		Port:              port,
		Id:                message.Identify(id),
		View:              message.View(view),
		Table:             t,
		Fault:             uint(len(t) / 3),
		ExecuteMaxNum:     1,
		CheckPointNum:     200,
		WaterL:            0,
		WaterH:            400,
		PrivateScalar:     GetPrivateScalar(),
		PublicSet:         GetPublicSet(),
		TblsPubPoly:       GetTblsPubPoly(),
		TblsPrivateScalar: GetTblsPrivateScalar(),
		ClientBatchSize:   100,
		ClientBatchTime:   time.Second,
	}
}

// 获取配置
func GetConfigureID() (id int, err error) {
	rawID := os.Getenv("PBFT_NODE_ID")
	if id, err = strconv.Atoi(rawID); err != nil {
		return
	}
	return
}

func GetConfigureTable() (map[int]string, error) {
	rawTable := os.Getenv("PBFT_NODE_TABLE")
	nodeTable := make(map[int]string, 0)

	tables := strings.Split(rawTable, ";")
	for index, t := range tables {
		nodeTable[index] = t
	}
	// 节点不满足 3f + 1
	if len(tables) < 3 || len(tables)%3 != 1 {
		return nil, errors.New("")
	}
	return nodeTable, nil
}

func GetConfigurePort() (port int, err error) {
	rawPort := os.Getenv("PBFT_LISTEN_PORT")
	if port, err = strconv.Atoi(rawPort); err != nil {
		return
	}
	return
}

func GetConfigureView() (int, error) {
	const ViewID = 0
	return ViewID, nil
}

func GetPrivateScalar() kyber.Scalar {
	raw := os.Getenv("PBFT_PRIVATE_KEY")
	ret, err := encoding.StringHexToScalar(edwards25519.NewBlakeSHA256Ed25519(), raw)
	if err != nil {
		log.Fatalf("[Config] read private key error")
	}
	return ret
}

func GetPublicSet() []kyber.Point {
	ret := make([]kyber.Point, 0)
	raw := os.Getenv("PBFT_PUBLIC_KEY")
	tables := strings.Split(raw, ";")
	for _, k := range tables {
		point, err := encoding.StringHexToPoint(edwards25519.NewBlakeSHA256Ed25519(), k)
		if err != nil {
			log.Fatalf("[Config] read public key error")
		}
		ret = append(ret, point)
	}
	return ret
}

func GetTblsPubPoly() *share.PubPoly {
	pubPoint := make([]kyber.Point, 0)
	raw := os.Getenv("PBFT_TBLS_PUBLIC_KEY")
	tables := strings.Split(raw, ",")
	for _, k := range tables {
		p, err := encoding.StringHexToPoint(bn256.NewSuite().G2(), k)
		if err != nil {
			log.Fatalf("[Config] read tbls public key error")
		}
		pubPoint = append(pubPoint, p)
	}
	return share.NewPubPoly(bn256.NewSuite().G2(), bn256.NewSuite().G2().Point().Base(), pubPoint)
}

func GetTblsPrivateScalar() kyber.Scalar {
	raw := os.Getenv("PBFT_TBLS_PRIVATE_KEY")
	p, err := encoding.StringHexToScalar(bn256.NewSuite().G2(), raw)
	if err != nil {
		log.Fatalf("[Config] read tbls private key error")
	}
	return p
}
