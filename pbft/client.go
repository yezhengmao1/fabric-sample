package pbft

import (
	"bytes"
	"encoding/json"
	"net/http"
)

// 发送请求
func Send(url string, msg []byte) error {
	buff := bytes.NewBuffer(msg)

	if _, err := http.Post("http://" + url, "application/json", buff); err != nil {
		logger.Infof("POST ERROR %s", err)
		return err
	}

	return nil
}

func SendReq(url string, req *RequestMsg) error {
	msg, err := json.Marshal(req)
	if err != nil {
		logger.Info(err)
		return err
	}
	logger.Infof("Client Send RequestMsg To %s", "http://" + url)
	return Send(url + URL_REQUEST, msg)
}

func SendPrePrepare(url string, preprepare *PrePrepareMsg) error {
	msg, err := json.Marshal(preprepare)
	if err != nil {
		logger.Info(err)
		return err
	}
	logger.Infof("Client Send PrePrepareMsg To %s", "http://" + url)
	return Send(url + URL_PREPREPARE, msg)
}

func SendPrepare(url string, prepare *PrepareMsg) error {
	msg, err := json.Marshal(prepare)
	if err != nil {
		logger.Info(err)
		return err
	}
	logger.Infof("Client Send PrepareMsg To %s", "http://" + url)
	return Send(url + URL_PREPARE, msg)
}

func SendCommit(url string, commit *CommitMsg) error {
	msg, err := json.Marshal(commit)
	if err != nil {
		logger.Info(err)
		return err
	}
	logger.Infof("Client Send CommitMsg To %s", "http://" + url)
	return Send(url + URL_COMMIT, msg)
}

func SendReply(url string, reply *ReplyMsg) error {
	msg, err := json.Marshal(reply)
	if err != nil {
		logger.Info(err)
		return err
	}
	logger.Infof("Client Send Reply To %s", "http://" + url)
	return Send(url + URL_REPLAY, msg)
}
