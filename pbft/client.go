package pbft

import (
	"bytes"
	"encoding/json"
	"net/http"
)

// 发送消息请求
func Send(url string, msg []byte) error {
	buff := bytes.NewBuffer(msg)
	if _, err := http.Post("http://" + url, "application/json", buff); err != nil {
		logger.Infof("POST ERROR %s", err)
		return err
	}
	return nil
}

func (n *Node) SendRequest(url string, req *RequestMsg) error {
	msg, err := json.Marshal(req)
	if err != nil {
		logger.Info(err)
		return err
	}
	logger.Infof("[PBFT Client] send request to %s", "http://" + url)
	return Send(url + URL_REQUEST, msg)
}

func (n *Node) SendPrePrepare(url string, prePrepare *PrePrepareMsg)  {
	msg, err := json.Marshal(prePrepare)
	if err != nil {
		logger.Info(err)
	}
	logger.Infof("[PBFT Client] send pre-prepare to %s", "http://" + url)
	_ = Send(url + URL_PREPREPARE, msg)
}

func (n *Node) SendPrepare(url string, prepare *PrepareMsg)  {
	msg, err := json.Marshal(prepare)
	if err != nil {
		logger.Info(err)
	}
	logger.Infof("[PBFT Client] send prepare to %s", "http://" + url)
	_ = Send(url + URL_PREPARE, msg)
}

func (n *Node) SendCommit(url string, commit *CommitMsg)  {
	msg, err := json.Marshal(commit)
	if err != nil {
		logger.Info(err)
	}
	logger.Infof("[PBFT Client] send commit to %s", "http://" + url)
	_ = Send(url + URL_COMMIT, msg)
}

func (n *Node) SendReply(url string, reply *ReplyMsg) {
	msg, err := json.Marshal(reply)
	if err != nil {
		logger.Info(err)
	}
	logger.Infof("[PBFT Client] send reply to %s", "http://" + url)
	_ = Send(url + URL_REPLAY, msg)
}
