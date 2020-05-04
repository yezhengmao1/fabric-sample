package crypto

import (
	"github.com/stretchr/testify/assert"
	"go.dedis.ch/kyber"
	"go.dedis.ch/kyber/group/edwards25519"
	"go.dedis.ch/kyber/pairing/bn256"
	"go.dedis.ch/kyber/share"
	"go.dedis.ch/kyber/util/encoding"
	"testing"
)

func TestRingSign(t *testing.T) {
	n := 4
	signer := 2
	msg := []byte("helloworld")
	suite := edwards25519.NewBlakeSHA256Ed25519()

	priSet := make([]kyber.Scalar, n)
	pubSet := make([]kyber.Point,  n)

	for i := 0; i < n; i++ {
		stream := suite.RandomStream()
		priSet[i] = suite.Scalar().Pick(stream)
		pubSet[i] = suite.Point().Mul(priSet[i], nil)
	}

	pri := priSet[signer]
	s := RingSign(msg, signer, pri, pubSet)
	// verify
	b := RingVerify(msg, s, pubSet)
	assert.Equal(t, true, b)
	s[0]++
	b = RingVerify(msg, s, pubSet)
	assert.Equal(t, false, b)
}

func TestTblsSign(t *testing.T) {
	n := 4
	v := 3

	msg := []byte("HelloWorld")

	suite := bn256.NewSuite()

	secret := suite.G1().Scalar().Pick(suite.RandomStream())

	priPoly := share.NewPriPoly(suite.G2(), v, secret, suite.RandomStream())
	pubPoly := priPoly.Commit(suite.G2().Point().Base())
	_, commits := pubPoly.Info()

	priKey := make([]string, 0)
	pubKey := make([]string, 0)

	for _, x := range commits {
		pub, _ := encoding.PointToStringHex(suite.G2(), x)
		pubKey = append(pubKey, pub)
	}

	for _, x := range priPoly.Shares(n) {
		pri, _ := encoding.ScalarToStringHex(suite.G2(), x.V)
		priKey = append(priKey, pri)
	}

	// 复原 key 和 pubPoly
	priScalr := make([]kyber.Scalar, 0)
	pubPoint := make([]kyber.Point, 0)

	for _, k := range pubKey {
		p, _ := encoding.StringHexToPoint(suite.G2(), k)
		pubPoint = append(pubPoint, p)
	}
	for _, k := range priKey {
		p, _ := encoding.StringHexToScalar(suite.G2(), k)
		priScalr = append(priScalr, p)
	}

	pubPoly = share.NewPubPoly(suite.G2(), suite.G2().Point().Base(), pubPoint)

	// pubPoly  - 共享公匙
	// proScalr - 私钥

	// 0 1 2 签
	sigPart := make([][]byte, 0)
	for i := 0; i < v; i++ {
		sig := TblsSign(msg, i, priScalr[i])
		sigPart = append(sigPart, sig)
	}
	content := TblsRecover(msg, sigPart, v, n, pubPoly)
	assert.NotEqual(t, len(content), 0)
	verify := TblsVerify(msg, content, pubPoly)
	assert.Equal(t, verify, true)
	// 1 2 3 签
	sigPart = make([][]byte, 0)
	for i := 1; i <= 3; i++ {
		sigPart = append(sigPart, TblsSign(msg, i, priScalr[i]))
	}
	content = TblsRecover(msg, sigPart, v, n, pubPoly)
	assert.NotEqual(t, len(content), 0)
	verify = TblsVerify(msg, content, pubPoly)
	assert.Equal(t, verify, true)
	// 0 1 签字
	sigPart = make([][]byte, 0)
	for i := 0; i <= 1; i++ {
		sigPart = append(sigPart, TblsSign(msg, i, priScalr[i]))
	}
	content = TblsRecover(msg, sigPart, v, n, pubPoly)
	assert.Equal(t, len(content), 0)
	verify = TblsVerify(msg, content, pubPoly)
	assert.Equal(t, verify, false)
	// 0 1 2 3 签字
	sigPart = make([][]byte, 0)
	for i := 0; i <= 3; i++ {
		sigPart = append(sigPart, TblsSign(msg, i, priScalr[i]))
	}
	content = TblsRecover(msg, sigPart, v, n, pubPoly)
	assert.NotEqual(t, len(content), 0)
	verify = TblsVerify(msg, content, pubPoly)
	assert.Equal(t, verify, true)
	// 0 1 1 1 签字
	sigPart = make([][]byte, 0)
	sigPart = append(sigPart, TblsSign(msg, 0, priScalr[0]))
	sigPart = append(sigPart, TblsSign(msg, 1, priScalr[1]))
	sigPart = append(sigPart, TblsSign(msg, 1, priScalr[1]))
	sigPart = append(sigPart, TblsSign(msg, 1, priScalr[1]))
	content = TblsRecover(msg, sigPart, v, n, pubPoly)
	assert.Equal(t, len(content), 0)
	verify = TblsVerify(msg, content, pubPoly)
	assert.Equal(t, verify, false)
	// 0 1 1 1 2 签字
	sigPart = make([][]byte, 0)
	sigPart = append(sigPart, TblsSign(msg, 0, priScalr[0]))
	sigPart = append(sigPart, TblsSign(msg, 1, priScalr[1]))
	sigPart = append(sigPart, TblsSign(msg, 1, priScalr[1]))
	sigPart = append(sigPart, TblsSign(msg, 1, priScalr[1]))
	sigPart = append(sigPart, TblsSign(msg, 2, priScalr[2]))
	content = TblsRecover(msg, sigPart, v, n, pubPoly)
	assert.Equal(t, len(content), 0)
	verify = TblsVerify(msg, content, pubPoly)
	assert.Equal(t, verify, false)
}

func TestSign(t *testing.T) {
	msg := []byte("helloworld")
	suite := edwards25519.NewBlakeSHA256Ed25519()

	pri := suite.Scalar().Pick(suite.RandomStream())
	pub := suite.Point().Mul(pri, nil)

	sign := Sign(msg, pub, pri)
	v    := Verify(msg, sign, pub)
	assert.Equal(t, v, true)
	pub2 := suite.Point().Mul(suite.Scalar().Pick(suite.RandomStream()), nil)
	assert.Equal(t, Verify(msg, sign, pub2), false)
}

