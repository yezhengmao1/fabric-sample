package pbft

import (
	"errors"
	"os"
	"strconv"
	"strings"
)

// 获取配置
func GetConfigureID() (id int, err error){
	rawID  := os.Getenv("PBFT_NODE_ID")
	if id, err = strconv.Atoi(rawID); err != nil {
		logger.Errorf("[PBFT CONFIGURE] get node id error %s", rawID)
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
		logger.Errorf("[PBFT CONFIGURE] get node num - %d error, not 3f + 1", len(tables))
		return nil, errors.New("")
	}
	return nodeTable, nil
}

func GetConfigurePort() (port int, err error){
	rawPort := os.Getenv("PBFT_LISTEN_PORT")
	if port, err = strconv.Atoi(rawPort); err != nil {
		logger.Errorf("[PBFT CONFIGURE] get server port  error %s", rawPort)
		return
	}
	return
}

func GetConfigureView() (int, error) {
	const ViewID = 100000
	return ViewID, nil
}


// 判断当前节点是否为主节点
func (n *Node) IsPrimary() bool {
	index := n.GetPrimary()
	if index == n.Id {
		return true
	}
	return false
}

// 获取Primary URL
func (n *Node) GetPrimaryUrl() string {
	return n.Table[n.GetPrimary()]
}

// 获取主节点
func (n *Node) GetPrimary() int {
	return n.View % len(n.Table)
}

// 生成一个Sequence
func (n *Node) GenerateSequence() int64 {
	return n.LastSequence + 1
}

type PrePrepareBuffer []*PrePrepareMsg

func (p PrePrepareBuffer) Len() int {
	return len(p)
}

func (p PrePrepareBuffer) Less(i, j int) bool {
	return p[i].Sequence < p[j].Sequence
}

func (p PrePrepareBuffer) Swap(i, j int) {
	p[i], p[j] = p[j], p[i]
}

type RequestMsgBuffer []*RequestMsg

func (p RequestMsgBuffer) Len() int {
	return len(p)
}

func (p RequestMsgBuffer) Less(i, j int) bool {
	return p[i].TimeStamp < p[j].TimeStamp
}

func (p RequestMsgBuffer) Swap(i, j int) {
	p[i], p[j] = p[j], p[i]
}