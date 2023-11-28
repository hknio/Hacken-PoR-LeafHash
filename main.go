package main

import (
	"crypto/sha512"
	"encoding/hex"
	"fmt"
)

type VerifierData struct {
	UUID    string
	Balance string
}

func (vd VerifierData) CalculateHash() ([]byte, error) {
	hash := sha512.Sum512([]byte(vd.UUID + vd.Balance))
	return hash[:], nil
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
