package example

import (
	"fmt"
	"testing"
)

func TestBuildTx(t *testing.T) {
	//Ready for inputs
	inputs := []Input{
		{
			TxID:    "6764086ab75301a2df0eba59336a9359db630ce37a394a97647b0897613ac521",
			Index:   0,
			Address: "16dEDpxrSoHeoaoeCDKG5QtomjVUQ7gPxb",
			Value:   5000,
			Wif:     "KwEyKT8abVSJJbhuo5d7S2XQQHUcF8CPS6jtBVDxX3tsQ1QMhfJZ",
		},
	}
	//Ready for ouputs
	ouputs := []Output{
		{
			Address: "1Q4hHsc8zV2Ukr4mmPdnFpjpSyL9KM8Aej",
			Value:   2000,
		},
	}
	opData := [][]byte{
		[]byte("Tx-demo"),
	}

	//Assemble and build tx
	tx := NewTx()
	tx.SetInputs(inputs)
	tx.SetOutputs(ouputs)
	tx.SetOpData(opData)
	txEntity, err := tx.BuildTx()
	if err != nil {
		fmt.Printf("MakeTx err:%s\n", err.Error())
		return
	}
	//After broadcasting TxHex to the block node, use txId to query its record on the block browser
	fmt.Printf("TxId:%s\n", txEntity.TxID())
	fmt.Printf("TxHex:%s\n", txEntity.String())
}


func TestBuildMetaIdNodeTx(t *testing.T) {
	//Ready for inputs
	inputs := []Input{
		{
			TxID:    "6764086ab75301a2df0eba59336a9359db630ce37a394a97647b0897613ac521",
			Index:   0,
			Address: "16dEDpxrSoHeoaoeCDKG5QtomjVUQ7gPxb",
			Value:   5000,
			Wif:     "KwEyKT8abVSJJbhuo5d7S2XQQHUcF8CPS6jtBVDxX3tsQ1QMhfJZ",
		},
	}
	//Ready for ouputs
	ouputs := []Output{
		{
			Address: "1Q4hHsc8zV2Ukr4mmPdnFpjpSyL9KM8Aej",
			Value:   2000,
		},
	}

	publicKey, err := getRandPublicKey()
	if err != nil {
		fmt.Printf("Get public key err:%s\n", err.Error())
		return
	}
	//Construct MetaId node data
	opData := [][]byte{
		[]byte{0x6d, 0x76, 0x63}, //Chain-flag:mvc
		[]byte(publicKey),        //PublicKey
		[]byte("mvc:NULL"),       //Chain:ParentTxId
		[]byte("metaid"),         //Metaid-flag:metaid
		[]byte("Root"),           //NodeName
		[]byte("NULL"),           //Data
		[]byte("0"),              //Encrypt
		[]byte("1.0.0"),          //Version
		[]byte("NULL"),           //DataType
		[]byte("NULL"),           //Encoding
	}

	//Assemble and build tx
	tx := NewTx()
	tx.SetInputs(inputs)
	tx.SetOutputs(ouputs)
	tx.SetOpData(opData)
	txEntity, err := tx.BuildTx()
	if err != nil {
		fmt.Printf("MakeTx err:%s\n", err.Error())
		return
	}
	//After broadcasting TxHex to the block node, use txId to query its record on the block browser
	fmt.Printf("TxId:%s\n", txEntity.TxID())
	fmt.Printf("TxHex:%s\n", txEntity.String())

}