package main

import (
	"crypto/sha512"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"math/big"
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
type HasherInput struct {
	H1 string `json:"h1"`
	H2 string `json:"h2"`
}

func SortAppend(sort bool, a, b []byte) []byte {
	if !sort {
		return append(a, b...)
	}
	var aBig, bBig big.Int
	aBig.SetBytes(a)
	bBig.SetBytes(b)
	if aBig.Cmp(&bBig) == -1 {
		return append(a, b...)
	}
	return append(b, a...)
}

// jsSumHashes is a wrapper for the SumHashes function
func jsSumHashes(this js.Value, args []js.Value) interface{} {
	jsonArgs := args[0].String()
	var hasherInput HasherInput
	err := json.Unmarshal([]byte(jsonArgs), &hasherInput)
	if err != nil {
		return ""
	}
	val := hasherInput.SumHashes()
	return val
}

// SumHashes returns the sum of two hashes
func (data HasherInput) SumHashes() string {
	h1bytes, _ := hex.DecodeString(data.H1)
	h2bytes, _ := hex.DecodeString(data.H2)
	bytes := SortAppend(true, h1bytes, h2bytes)
	hash := sha512.Sum512(bytes)
	return hex.EncodeToString(hash[:])
}

// jsGetLeafHash is a wrapper for the GetLeafHash function
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

// CalculateHash returns the hash of the VerifierData
func (vd VerifierData) CalculateHash() ([]byte, error) {
	hash := sha512.Sum512([]byte(vd.UUID + vd.Balance))
	return hash[:], nil
}

// GetLeafHash returns the hash of the VerifierData
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
	js.Global().Set("GetLeafHash", js.FuncOf(jsGetLeafHash))
	js.Global().Set("SumHashes", js.FuncOf(jsSumHashes))
	<-c
}
