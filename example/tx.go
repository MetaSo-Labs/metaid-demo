package example

import (
	"context"
	"encoding/hex"
	"errors"
	"github.com/mvc-labs/mvc-lib-go"
	"github.com/mvc-labs/mvc-lib-go/bscript"
	"github.com/mvc-labs/mvc-lib-go/keys/wif"
	"github.com/mvc-labs/mvc-lib-go/sighash"
	"github.com/mvc-labs/mvc-lib-go/unlocker"
	"math"
)

type Tx struct {
	inputs        []Input
	outputs       []Output
	opData        [][]byte
	changeAddress string
	fee           float64
	isBuild       bool
	outUtxos      []Utxo
}

func NewTx() *Tx {
	return &Tx{fee: 1}
}

func (t *Tx) SetInputs(inputs []Input) {
	t.inputs = inputs
}

func (t *Tx) SetOutputs(outputs []Output) {
	t.outputs = outputs
}

func (t *Tx) SetChangeAddress(changeAddress string) {
	t.changeAddress = changeAddress
}

func (t *Tx) SetOpData(opData [][]byte) {
	t.opData = opData
}

func (t *Tx) SetFee(fee float64) {
	t.fee = fee
}


//Build tx
func (t *Tx) BuildTx() (*bt.Tx, error) {
	var (
		inputAmount  uint64 = 0
		outputAmount uint64 = 0
		newTx                  = bt.NewTx()
		utxos               = make([]*bt.UTXO, 0)
		outUtxos               = make([]Utxo, 0)
		err error
		tmpChangeSatoshi uint64 = 10000000000
	)
	if len(t.inputs) == 0 {
		return nil, errors.New("inputs are empty")
	}
	if len(t.changeAddress) == 0 {
		t.changeAddress = t.inputs[0].Address
	}
	newTx.Version = 10

	//Assemble inputs in newTx
	for _, input := range t.inputs {
		scriptString, err := bscript.NewP2PKHFromAddress(input.Address)
		if err != nil {
			return nil, err
		}
		pti, err := hex.DecodeString(input.TxID)
		if err != nil {
			return nil, err
		}
		utxos = append(utxos, &bt.UTXO{
			TxID:          pti,
			Vout:          input.Index,
			LockingScript: scriptString,
			Satoshis:      input.Value,
		})
		inputAmount += input.Value
	}
	err = newTx.FromUTXOs(utxos...)
	if err != nil {
		return nil, err
	}

	//Assemble outputs in newTx
	_ = newTx.PayToAddress(t.changeAddress, tmpChangeSatoshi)
	for _, output := range t.outputs {
		err := newTx.PayToAddress(output.Address, output.Value)
		if err != nil {
			return nil, err
		}
		outputAmount += output.Value
	}

	//Assemble opData in newTx
	if t.opData != nil {
		_ = newTx.AddOpReturnPartsOutput(t.opData)
	}

	//Calculate transaction size and return the remaining gas to the change address
	size := newTx.Size()
	satoshi := math.Ceil(float64(size) * t.fee)
	change := (inputAmount - outputAmount) - uint64(satoshi)
	newTx.Outputs = newTx.Outputs[1:]
	if change >= t.calDust() {
		_ = newTx.PayToAddress(t.changeAddress, change)
	}


	//Sign and unlock the script
	for index, input := range t.inputs {
		wifKey, err := wif.DecodeWIF(input.Wif)
		if err != nil {
			return nil, err
		}
		err = newTx.FillInput(context.Background(), &unlocker.Simple{PrivateKey: wifKey.PrivKey}, bt.UnlockerParams{
			InputIdx:     uint32(index),
			SigHashFlags: sighash.AllForkID,
		})
		if err != nil {
			return nil, err
		}
	}
	t.isBuild = true

	for index, output := range newTx.Outputs {
		if output.LockingScript.IsP2PKH() {
			addresses, _ := output.LockingScript.Addresses()
			outUtxos = append(outUtxos, Utxo{
				Address:addresses[0],
				Value: output.Satoshis,
				Index: uint64(index),
				Tx:    newTx.TxID(),
			})
		}
	}
	t.outUtxos = outUtxos
	return newTx, nil
}

func (t *Tx) GetOutUtxos() ([]Utxo, error) {
	if !t.isBuild {
		return nil, errors.New("Tx had not built. ")
	}
	return t.outUtxos, nil
}

//Calculate dust
func (t *Tx) calDust() uint64 {
	dust := 3 * math.Ceil(t.fee * 182)
	return uint64(dust)
}