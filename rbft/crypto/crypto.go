package crypto

import (
	"go.dedis.ch/kyber"
	"go.dedis.ch/kyber/group/edwards25519"
	"go.dedis.ch/kyber/sign/anon"
	"go.dedis.ch/kyber/sign/eddsa"
	"log"
)

// ring signature and verify
func RingSign(msg []byte, id int, pri kyber.Scalar, pubSet []kyber.Point) []byte {
	suite := edwards25519.NewBlakeSHA256Ed25519()
	// pubset to point
	content := anon.Sign(suite, msg, pubSet, nil, id, pri)
	return content
}

func RingVerify(msg []byte, sign []byte, pubSet []kyber.Point) bool {
	suite := edwards25519.NewBlakeSHA256Ed25519()
	tag, err := anon.Verify(suite, msg, pubSet, nil, sign)
	if err != nil || tag == nil || len(tag) != 0 {
		return false
	}
	return true
}

// signature and verify
func Sign(msg []byte, pub kyber.Point, pri kyber.Scalar) []byte {
	ed := eddsa.EdDSA{
		Secret: pri,
		Public: pub,
	}
	sign, err := ed.Sign(msg)
	if err != nil {
		log.Printf("[Crypto] error to sign the message")
		return nil
	}
	return sign
}

func Verify(msg, sig []byte, pub kyber.Point) bool {
	err := eddsa.Verify(pub, msg, sig)
	if err != nil {
		return false
	}
	return true
}
