package main

import (
	"crypto/sha512"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"syscall/js"
)

type VerifierData struct {
	UUID    string
	Balance string
}

type InputData struct {
	UUID    string     `json:"uuid"`
	Balance [][]string `json:"balance"`
}

func (vd VerifierData) CalculateHash() ([]byte, error) {
	hash := sha512.Sum512([]byte(vd.UUID + vd.Balance))
	return hash[:], nil
}

func jsGetLeafHash(this js.Value, args []js.Value) interface{} {
	jsonArgs := args[0].String()
	var verifierData InputData
	err := json.Unmarshal([]byte(jsonArgs), &verifierData)
	if err != nil {
		return ""
	}
	val, _ := GetLeafHash(verifierData.UUID, verifierData.Balance)

	return val
}

func GetLeafHash(uuid string, balance [][]string) (string, error) {

	balanceStr := fmt.Sprintf("%v", balance)
	dataToHash := VerifierData{UUID: uuid, Balance: balanceStr}

	leafHash, err := dataToHash.CalculateHash()
	if err != nil {
		return "", err
	}

	return hex.EncodeToString(leafHash), nil
}

func main() {
	c := make(chan struct{}, 0)
	js.Global().Set("getLeafHash", js.FuncOf(jsGetLeafHash))
	<-c
}
