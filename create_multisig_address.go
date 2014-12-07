package main

import (
	"github.com/soroushjp/go-bitcoin-multisig/base58check"
	"github.com/soroushjp/go-bitcoin-multisig/btcutils"

	"encoding/csv"
	"encoding/hex"
	"flag"
	"fmt"
	"log"
	"strings"
)

var flagPublicKeys string
var flagM int
var flagN int

const REQUIRED_FLAG_COUNT = 3

func main() {
	//Parse flags
	flag.IntVar(&flagM, "m", 0, "M, the minimum number of keys needed to spend Bitcoin in M-of-N multisig transaction.")
	flag.IntVar(&flagN, "n", 0, "N, the total number of possible keys that can be used to spend Bitcoin in M-of-N multisig transaction.")
	flag.StringVar(&flagPublicKeys, "public-keys", "", "Comma separated list of private keys to sign with. Whitespace is stripped and quotes may be placed around keys. Eg. key1,key2,\"key3\" .")
	flag.Parse()
	if flag.NFlag() != REQUIRED_FLAG_COUNT {
		//We only need to check flag count because Go will automatically throw an error for undefined flags
		log.Fatal("Please provide all required flags.")
	}
	//Convert public keys argument into slice of public key bytes with necessary tidying
	flagPublicKeys = strings.Replace(flagPublicKeys, "'", "\"", -1) //Replace single quotes with double since csv package only recognizes double quotes
	publicKeyStrings, err := csv.NewReader(strings.NewReader(flagPublicKeys)).Read()
	if err != nil {
		log.Fatal(err)
	}
	publicKeys := make([][]byte, len(publicKeyStrings))
	for i, publicKeyString := range publicKeyStrings {
		publicKeyString = strings.TrimSpace(publicKeyString)   //Trim whitespace
		publicKeys[i], err = hex.DecodeString(publicKeyString) //Get private keys as slice of raw bytes
		if err != nil {
			log.Fatal(err)
		}
	}
	//Create redeemScript from public keys
	//redeemScript := btcutils.NewTwoOfTwoRedeemScript(publicKeys[0], publicKeys[1])
	redeemScript, err := btcutils.NewMOfNRedeemScript(flagM, flagN, publicKeys)
	if err != nil {
		log.Fatal(err)
	}
	redeemScriptHash, err := btcutils.Hash160(redeemScript)
	if err != nil {
		log.Fatal(err)
	}
	//Get P2SH address by base58 encodin with P2SH prefix 0x05
	P2SHAddress := base58check.Encode("05", redeemScriptHash)
	//Output P2SH and redeemScript
	fmt.Println("---------------------")
	fmt.Println("Your *P2SH ADDRESS* is:")
	fmt.Println(P2SHAddress)
	fmt.Println("Give this to sender funding multisig address with Bitcoin.")
	fmt.Println("---------------------")
	fmt.Println("---------------------")
	fmt.Println("Your *REDEEM SCRIPT* is:")
	fmt.Println(hex.EncodeToString(redeemScript))
	fmt.Println("Keep private and provide this to redeem multisig balance later.")
	fmt.Println("---------------------")
}
