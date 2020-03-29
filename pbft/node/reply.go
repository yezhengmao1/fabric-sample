package node

import (
	"github.com/hyperledger/fabric/orderer/consensus/pbft/message"
	"github.com/hyperledger/fabric/orderer/consensus/pbft/server"
	cb "github.com/hyperledger/fabric/protos/common"
	"log"
)

var test_reqeust_num uint64 = 0

func (n *Node) executeAndReplyThread() {
	for {
		select {
		case <-n.executeNotify:
			// execute batch
			batchs, lastSeq := n.buffer.BatchExecute(n.sequence.GetLastSequence())
			if len(batchs) == 0 {
				log.Printf("[Reply] lost sequence now(%d)", n.sequence.GetLastSequence())
				continue
			}
			n.sequence.SetLastSequence(lastSeq)
			// check point
			if n.sequence.ReadyToCheckPoint() {
				checkSeq := n.sequence.GetCheckPoint()
				content, checkPoint := n.buffer.CheckPoint(checkSeq, n.id)
				// buffer checkpoint
				n.buffer.BufferCheckPointMsg(checkPoint, n.id)
				log.Printf("[Reply] ready to create check point to sequence(%d) msg(%s)", checkSeq, checkPoint.Digest[0:9])
				n.BroadCast(content, server.CheckPointEntry)
			}
			// map the digest to request
			requestBatchs := make([]*message.Request, 0)
			for _, b := range batchs {
				requestBatchs = append(requestBatchs, b.Message.Requests...)
			}
			test_reqeust_num = test_reqeust_num + uint64(len(requestBatchs))
			log.Printf("[Reply] set last sequence(%d) already execute request(%d)", lastSeq, test_reqeust_num)
			// pending state
			pending := make(map[string]bool)
			for _, r := range requestBatchs {
				msg		  := r.Op.Envelope
				channel   := r.Op.ChannelID
				configSeq := r.Op.ConfigSeq
				switch r.Op.Type {
				case message.TYPECONFIG:
					var err error
					seq := n.supports[channel].Sequence()
					if configSeq < seq {
						if msg, _, err = n.supports[r.Op.ChannelID].ProcessConfigMsg(r.Op.Envelope); err != nil {
							log.Println(err)
						}
					}
					batch := n.supports[channel].BlockCutter().Cut()
					if batch != nil {
						block := n.supports[channel].CreateNextBlock(batch)
						n.supports[channel].WriteBlock(block, nil)
					}
					pending[channel] = false
					// write config block
					block := n.supports[channel].CreateNextBlock([]*cb.Envelope{msg})
					n.supports[channel].WriteConfigBlock(block, nil)
				case message.TYPENORMAL:
					seq := n.supports[channel].Sequence()
					if configSeq < seq {
						if _, err := n.supports[channel].ProcessNormalMsg(msg); err != nil {
						}
					}
					batches, p := n.supports[channel].BlockCutter().Ordered(msg)
					for _, batch := range batches {
						block := n.supports[channel].CreateNextBlock(batch)
						n.supports[channel].WriteBlock(block, nil)
					}
					pending[channel] = p
				}
			}
			for k, v := range pending {
				if v {
					batch := n.supports[k].BlockCutter().Cut()
					if batch != nil {
						block := n.supports[k].CreateNextBlock(batch)
						n.supports[k].WriteBlock(block, nil)
					}
				}
			}
		}
	}
}
