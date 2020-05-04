package cmd

import (
	"github.com/hyperledger/fabric/orderer/consensus/rbft/message"
	"go.dedis.ch/kyber"
	"go.dedis.ch/kyber/share"
	"time"
)

type SharedConfig struct {
	ClientServer      bool
	Port              int
	Id                message.Identify
	View              message.View
	Table             map[message.Identify]string
	Fault             uint
	ExecuteMaxNum     int
	CheckPointNum     message.Sequence
	WaterL            message.Sequence
	WaterH            message.Sequence
	PrivateScalar     kyber.Scalar
	PublicSet         []kyber.Point
	TblsPubPoly       *share.PubPoly
	TblsPrivateScalar kyber.Scalar
	ClientBatchSize   int
	ClientBatchTime   time.Duration
}
