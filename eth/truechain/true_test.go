package truechain

import (
	"testing"
	"fmt"
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	//"github.com/ethereum/go-ethereum/rlp"
	"github.com/ethereum/go-ethereum/core/types"
	"crypto/ecdsa"
)


var (
	privkeys = make([]*ecdsa.PrivateKey,0,0)
	keysCount = 6
	tx1 = types.NewTransaction(
		0,
		common.HexToAddress("095e7baea6a6c7c4c2dfeb977efac326af552d87"),
		big.NewInt(0), 0, big.NewInt(0),
		nil,
	)
	tx2 = types.NewTransaction(
		0,
		common.HexToAddress("095e7baea6a6c7c4c2dfeb977efac326af552d87"),
		big.NewInt(0), 0, big.NewInt(0),
		nil,
	)
	tx3, _ = types.NewTransaction(
		3,
		common.HexToAddress("b94f5374fce5edbc8e2a8697c15331677e6ebf0b"),
		big.NewInt(10),
		2000,
		big.NewInt(1),
		common.FromHex("5544"),
	).WithSignature(
		types.NewEIP155Signer(common.Big1),
		common.Hex2Bytes("98ff921201554726367d2be8c804a7ff89ccf285ebc57dff8ae4c44b9c19ac4a8887321be575c8095f789dd4c743dfe42c1820f9231f98a962b210e3ac2452a301"),
	)
)
func init(){
	for i:=0;i<keysCount;i++ {
		k,_ := crypto.GenerateKey()
		privkeys = append(privkeys,k)
	}
}

func MakePbftBlock(cmm *PbftCommittee) *TruePbftBlock {
	txs := make([]*types.Transaction,0,0)
	txs = append(txs,tx1,tx2)
	pbTxs := make([]*Transaction,0,0)
	for _,vv := range txs {
		to := make([]byte,0,0)
		if tt := vv.To(); tt != nil {
			to = tt.Bytes()
		}
		v,r,s := vv.RawSignatureValues()
		pbTxs = append(pbTxs,&Transaction{
			Data:       &TxData{
				AccountNonce:       vv.Nonce(),
				Price:              vv.GasPrice().Int64(),
				GasLimit:           new(big.Int).SetUint64(vv.Gas()).Int64(),
				Recipient:          to,
				Amount:             vv.Value().Int64(),
				Payload:            vv.Data(),
				V:                  v.Int64(),
				R:                  r.Int64(),
				S:                  s.Int64(),
			},
		})
	}
	// begin make pbft block
	now := time.Now().Unix()
	head := TruePbftBlockHeader{
		Number:				10,
		GasLimit:			100,
		GasUsed:			80,
		Time:				now,
	}
	block := TruePbftBlock{
		Header:			&head,
		Txs:			&Transactions{Txs:pbTxs},
	}
	msg := rlpHash(block.Txs)
	//cc := cmm.GetCmm()
	sigs := make([]string,0,0)
	// same priveatekey to sign the message
	for i:=0;i<keysCount/2;i++ {
		sig,err := crypto.Sign(msg,privkeys[i])
		if err == nil {
			sigs = append(sigs,common.ToHex(sig))
		}
	}
	return &block
}
//func MakeFirstCommittee(num int) {
//
//}
func TestCryptoMsg(t *testing.T) {

}
func TestMainMembers(t *testing.T) {
	var (
		th = New()
	)
	err := th.StartTrueChain(nil)
	if err != nil {
		fmt.Println(err)
	}
	// test make pbft block
	block := MakePbftBlock(th.Cmm)
	// verify the pbft block
	err = th.CheckBlock(block)
	if err != nil {
		fmt.Println("verify the block failed,err=",err)
		return
	}
	// test cryptomsg for candidate Member

	th.StopTrueChain()
}