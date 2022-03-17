package main

import (
	"context"
	"fmt"
	"math/big"
	"io/ioutil"
	_"testing"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/accounts/keystore"
	_"github.com/stretchr/testify/require"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	_"github.com/ethereum/go-ethereum/params"
)

var fromKeyStoreFile = "./keystore/testkey.json"
var password = ""
var toAddress = "0xfCd09085A92d2d2Af961624Fa88D22A8beBdc344"
//"0xD75924Ac4b879E15290134361884F00e57397D95"
var httpUrl = "http://192.168.56.101:8545"
//var callName = "poke(bytes32)"
//var callValue = ""

func SendTx(/*t *testing.T*/){
	// 交易发送方
	// 获取私钥方式一，通过keystore文件
	fromKeystore,err := ioutil.ReadFile(fromKeyStoreFile)
	if err != nil {
        fmt.Println(err)
    }
	//require.NoError(t,err)
	fromKey,err := keystore.DecryptKey(fromKeystore,password)
	fromPrivkey := fromKey.PrivateKey
	fromPubkey := fromPrivkey.PublicKey
	fromAddr := crypto.PubkeyToAddress(fromPubkey)

	// 获取私钥方式二，通过私钥字符串
	//privateKey, err := crypto.HexToECDSA("私钥字符串")

	toAddr := common.HexToAddress(toAddress)
	var amount *big.Int = big.NewInt(0) //params.Ether //2e18
	var gasLimit uint64 = 3000000
	var gasPrice *big.Int = big.NewInt(200)

	client, err:= ethclient.Dial(httpUrl)
	if err != nil {
        fmt.Println(err)
    }
	//require.NoError(t, err)

	nonce, err := client.PendingNonceAt(context.Background(), fromAddr)

	auth, err := bind.NewKeyedTransactorWithChainID(fromPrivkey, big.NewInt(1337))
	if err != nil {
		fmt.Println(err);
	}
	//auth,err := bind.NewTransactor(strings.NewReader(mykey),"111")
	auth.Nonce = big.NewInt(int64(nonce))
	auth.Value = amount     // in wei
	//auth.Value = big.NewInt(100000)     // in wei
	auth.GasLimit = gasLimit // in units
	//auth.GasLimit = uint64(0) // in units
	auth.GasPrice = gasPrice
	auth.From = fromAddr

	//===================================//
	methodName := []byte("poke(bytes32)")
	methodID := crypto.Keccak256(methodName)[:4]

	price := new(big.Int)
    price.SetString("12345", 10)
    paddedPrice := common.LeftPadBytes(price.Bytes(), 32)

	var data []byte
    data = append(data, methodID...)
    data = append(data, paddedPrice...)

	// 交易创建
	tx := types.NewTransaction(nonce, toAddr, amount, gasLimit, gasPrice, data) //[]byte{}

	// 交易签名
	signedTx ,err := auth.Signer(auth.From, tx)
	//signedTx ,err := types.SignTx(tx,types.HomesteadSigner{},fromPrivkey)
	if err != nil {
        fmt.Println("auth signer error: %s", err)
    }
	//require.NoError(t, err)

	// 交易发送
	serr := client.SendTransaction(context.Background(), signedTx)
	if serr != nil {
		fmt.Println(serr)
	}

	// 等待挖矿完成
	bind.WaitMined(context.Background(), client, signedTx)
}

func main() {
	fmt.Println("Feed price started.")
	SendTx()
	fmt.Println("Feed price finished.")
}
