package keys

import (
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
)

type Case struct {
	//the inputs
	m           int          //threshold
	msg         string       //message to sign
	signers     []Address    //signers who will participate
	testMsg     string       //test message to sign, to test sign wrong msg
	testSigners []PrivateKey //test signers who will sign

	//the outputs
	errInit error  //expected errors for init function
	errAdd  error  //expected error for AddSignature function
	valid   bool   //is it a valid threshold multisig
	log     string //test logs
}

const numberOfPrivKey = 16

var testCases map[int]Case

func init() {
	var privKeys []PrivateKey = make([]PrivateKey, numberOfPrivKey)
	var addresses []Address = make([]Address, numberOfPrivKey)

	for i := 0; i < numberOfPrivKey; i++ {
		pub, priv, _ := NewKeyPairFromTendermint()
		privKeys[i] = priv

		h, _ := pub.GetHandler()
		addresses[i] = h.Address()
	}
	testCases = make(map[int]Case)
	testCases[0] = Case{1, "test0", addresses[:3], "test0", privKeys[:1], nil, nil, true, "1-of-3 multisig"}
	testCases[1] = Case{4, "test1", addresses[:3], "test1", nil, ErrInvalidThreshold, nil, false, "wrong threshold"}
	testCases[2] = Case{3, "test2", addresses[:5], "test2", privKeys[:2], nil, nil, false, "not enough signers has signed"}
	testCases[3] = Case{3, "test3", nil, "test3", nil, ErrMissSigners, nil, false, "not enough signer for init"}
	testCases[4] = Case{3, "test4", addresses[:5], "test4", privKeys[5:7], nil, ErrNotExpectedSigner, false, "wrong signers to has signed"}
	testCases[5] = Case{5, "test5", addresses[:5], "test5", privKeys[:5], nil, nil, true, "5-of-5 multisig"}
	testCases[6] = Case{3, "test6", addresses[:5], "test6", privKeys[:3], nil, nil, true, "3-of-5 multisig"}
	testCases[7] = Case{3, "test7", addresses[:5], "test7", privKeys[2:5], nil, ErrNotExpectedSigner, false, "right signer but wrong index"}
	testCases[8] = Case{3, "test8", addresses[:5], "test", privKeys[:3], nil, ErrInvalidSignedMsg, false, "wrong signed msg"}
	testCases[9] = Case{3, "test9", addresses[:5], "test9", privKeys[:2], nil, nil, false, "3-of-5 multisig"}
	testCases[10] = Case{5, "test10", addresses[:5], "test10", privKeys[:4], nil, nil, false, "5-of-5 multisig"}
}

func TestMultiSig(t *testing.T) {

	for i, item := range testCases {
		t.Run("Testing case "+strconv.Itoa(i), func(t *testing.T) {
			//test create MultiSig

			msg := []byte(item.msg)

			ms := &MultiSig{}
			err := ms.Init(msg, item.m, item.signers)
			if item.errInit != nil {
				assert.EqualError(t, err, item.errInit.Error(), "did not get expected error for [case %d]: %s", i, item.log)
			} else {
				assert.NoError(t, err, "get unexpected error for [case %d]: %s", i, item.log)
			}

			// test add signature
			if item.testSigners == nil {
				return
			}
			for i, priv := range item.testSigners {
				signMsg := []byte(item.testMsg)
				h, err := priv.GetHandler()
				assert.NoError(t, err, "unexpected error in sign", err)
				signed, err := h.Sign(signMsg)
				assert.NoError(t, err, "unexpected error in sign", err)
				signature := Signature{Index: i, PubKey: h.PubKey(), Signed: signed}
				err = ms.AddSignature(signature)
				if item.errAdd != nil {
					assert.EqualError(t, err, item.errAdd.Error(), "did not get expected error for [case %d]: %s", i, item.log)
				} else {
					assert.NoError(t, err, "get unexpected error for [case %d]: %s", i, item.log)
				}
			}

			//test valid

			result := ms.IsValid()
			assert.Equal(t, item.valid, result, "unexpected validation of MultiSig %s", item.log)
		})
	}

}

func TestMultiSig_Bytes(t *testing.T) {
	c := testCases[0]

	ms := &MultiSig{}
	err := ms.Init([]byte(c.msg), c.m, c.signers)
	assert.NoError(t, err, "unexpected failed to init")

	b := ms.Bytes()

	newms := &MultiSig{}
	err = newms.FromBytes(b)
	assert.NoError(t, err, "failed deser %s", err)

	assert.Equal(t, ms, newms, "unmatch after ser/deser %#v. %#v", ms, newms)

	signMsg := []byte(c.testMsg)
	h, err := c.testSigners[0].GetHandler()
	assert.NoError(t, err, "unexpected error in sign", err)
	signed, err := h.Sign(signMsg)
	assert.NoError(t, err, "unexpected error in sign", err)
	signature := Signature{Index: 0, PubKey: h.PubKey(), Signed: signed}
	err = ms.AddSignature(signature)
	assert.NoError(t, err, "get unexpected error for [case 0]: %s", c.log)

	b = ms.Bytes()

	newSignedMS := &MultiSig{}
	err = newSignedMS.FromBytes(b)
	assert.NoError(t, err, "failed deser %s", err)

	assert.Equal(t, ms, newSignedMS, "unmatch after ser/deser %#v. %#v", ms, newms)
}
