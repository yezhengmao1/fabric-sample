package crypto

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"go.dedis.ch/kyber"
	"go.dedis.ch/kyber/group/edwards25519"
	"testing"
)

func TestRingSign(t *testing.T) {
	n := 4
	signer := 2
	msg := []byte("helloworld")
	suite := edwards25519.NewBlakeSHA256Ed25519()

	priSet := make([]kyber.Scalar, n)
	pubSet := make([]kyber.Point,  n)
	tmp := ""
	for i := 0; i < n; i++ {
		stream := suite.RandomStream()
		priSet[i] = suite.Scalar().Pick(stream)
		pubSet[i] = suite.Point().Mul(priSet[i], nil)
		tmp = tmp + pubSet[i].String() + ";"

		fmt.Println(stream)
		fmt.Println(priSet[i])
	}
	fmt.Println(tmp)

	pri := priSet[signer]
	s := RingSign(msg, signer, pri, pubSet)
	// verify
	b := RingVerify(msg, s, pubSet)
	assert.Equal(t, true, b)
	s[0]++
	b = RingVerify(msg, s, pubSet)
	assert.Equal(t, false, b)
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

