package cmd

import (
	"errors"
	"flag"
	"github.com/hyperledger/fabric/orderer/consensus/pbft/message"
	"log"
	"os"
	"strconv"
	"strings"
)

type SharedConfig struct {
	ClientServer   bool
	Port           int
	Id             message.Identify
	View           message.View
	Table		   map[message.Identify]string
	FaultNum	   uint
	ExecuteMaxNum  int
	CheckPointNum  message.Sequence
	WaterL         message.Sequence
	WaterH		   message.Sequence
}

func ReadConfig() *SharedConfig {
	port, _ := GetConfigurePort()
	id, _   := GetConfigureID()
	view, _ := GetConfigureView()
	table, _ := GetConfigureTable()

	t := make(map[message.Identify]string)
	for k, v := range table {
		t[message.Identify(k)] = v
	}
	// calc the fault num
	if len(t) % 3 != 1 {
		log.Fatalf("[Config Error] the incorrent node num : %d, need 3f + 1", len(t))
		return nil
	}

	flag.Parse()
	return &SharedConfig{
		Port: 		   port,
		Id:            message.Identify(id),
		View:          message.View(view),
		Table:         t,
		FaultNum:      uint(len(t)/3),
		ExecuteMaxNum: 1,
		CheckPointNum: 200,
		WaterL:        0,
		WaterH:		   400,
	}
}

// 获取配置
func GetConfigureID() (id int, err error){
	rawID  := os.Getenv("PBFT_NODE_ID")
	if id, err = strconv.Atoi(rawID); err != nil {
		return
	}
	return
}

func GetConfigureTable() (map[int]string, error){
	rawTable  := os.Getenv("PBFT_NODE_TABLE")
	nodeTable := make(map[int]string, 0)

	tables := strings.Split(rawTable, ";")
	for index, t := range tables {
		nodeTable[index] = t
	}
	// 节点不满足 3f + 1
	if len(tables) < 3 || len(tables) % 3 != 1 {
		return nil, errors.New("")
	}
	return nodeTable, nil
}

func GetConfigurePort() (port int, err error){
	rawPort := os.Getenv("PBFT_LISTEN_PORT")
	if port, err = strconv.Atoi(rawPort); err != nil {
		return
	}
	return
}

func GetConfigureView() (int, error) {
	const ViewID = 100000
	return ViewID, nil
}
