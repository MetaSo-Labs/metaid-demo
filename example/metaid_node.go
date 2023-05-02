package example

import (
	"errors"
	"fmt"
	bt "github.com/mvc-labs/mvc-lib-go"
)

type Node struct {
	chain       string
	publicKey   string
	parentChain string
	parentTxId  string
	nodeName    string
	data        []byte
	encrypt     string
	version     string
	dataType    string
	encoding    string
}

func NewNode() *Node {
	return &Node{
		chain:       "mvc",
		encrypt:     "0",
		version:     "1.0.0",
		dataType:    "application/json",
		encoding:    "UTF-8",
	}
}

func (n *Node) SetPublicKey(publicKey string) {
	n.publicKey = publicKey
}
func (n *Node) SetParentNode(parentChain, parentTxId string) {
	n.parentChain = parentChain
	n.parentTxId = parentTxId
}
func (n *Node) SetNodeData(nodeName string, data []byte, encrypt string) {
	n.nodeName = nodeName
	n.data = data
	n.encrypt = encrypt
}
func (n *Node) SetVersion(version string) {
	n.version = version
}
func (n *Node) SetDataType(dataType string) {
	n.dataType = dataType
}
func (n *Node) SetEncoding(encoding string) {
	n.encoding = encoding
}

func (n *Node)BuildNode(utxos []Utxo, priKeys map[string]string)  {

}

func (n *Node) ToOpData() [][]byte {
	opData := [][]byte{
		[]byte(n.chain),
		[]byte(n.publicKey),
		[]byte(fmt.Sprintf("%s:%s", n.parentChain, n.parentTxId)),
		[]byte(n.nodeName),
		n.data,
		[]byte(n.encrypt),
		[]byte(n.version),
		[]byte(n.dataType),
		[]byte(n.encoding),
	}
	return opData
}




//Build MetaId func: Root, Info, Protocol, name
func BuildNewMetaId(nodeNames map[string]string, nickName string, utxos []Utxo, priKeys map[string]string) ([]string, []string, error) {
	var(
		txIds []string = make([]string, 0)
		txHexes []string = make([]string, 0)
	)
	//Build Root tx

	
	
	return txIds, txHexes, nil
}




//Encapsulate the common MetaId data construction function
func BuildCommonTx(utxos []Utxo, priKeys map[string]string, opData [][]byte, outputs []Output, changeAddress string) (*bt.Tx, []Utxo, error) {
	if len(utxos) == 0 {
		return nil, nil, errors.New("The utoxs are empty. ")
	}
	if len(priKeys) == 0 {
		return nil, nil, errors.New("The priKeys are empty. ")
	}

	inputs := make([]Input, 0)
	for _, u := range utxos{
		if _,ok := priKeys[u.Address];!ok {
			return nil, nil, errors.New(fmt.Sprintf("The utxoAddress of priKeys is empty. Address:%s", u.Address))
		}
		inputs = append(inputs, Input{
			TxID:    u.Tx,
			Index:   uint32(u.Index),
			Address: u.Address,
			Value:   u.Value,
			Wif:     priKeys[u.Address],
		})
	}

	tx := NewTx()
	tx.SetInputs(inputs)
	tx.SetOutputs(outputs)
	tx.SetChangeAddress(changeAddress)
	tx.SetOpData(opData)
	txEntity, err := tx.BuildTx()
	if err != nil {
		return nil, nil, errors.New(fmt.Sprintf("MakeTx err:%s\n", err.Error()))
	}
	outUtxos, err := tx.GetOutUtxos()
	if err != nil {
		return nil, nil, err
	}
	newOutUtxos := make([]Utxo, 0)
	for _, u := range outUtxos {
		if _, ok := priKeys[u.Address]; ok {
			newOutUtxos = append(newOutUtxos, u)
		}
	}
	return txEntity, newOutUtxos, nil
}