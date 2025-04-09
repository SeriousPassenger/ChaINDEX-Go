package main

import (
	"encoding/json"
	"fmt"
	"math/big"
	"os"
)

func HexToBigInt(hexString string) (*big.Int, error) {
	if len(hexString) < 3 || hexString[:2] != "0x" {
		return nil, fmt.Errorf("invalid hex string: %s", hexString)
	}

	// Convert the hex string to a big.Int
	value := new(big.Int)
	value.SetString(hexString[2:], 16)

	return value, nil
}

func BigIntToHex(value *big.Int) string {
	hexString := fmt.Sprintf("0x%x", value)

	return hexString
}

func SaveStructToJSONFile(data interface{}, filename string) error {
	// Save rawData to a file indent 4
	file, err := os.Create(filename)

	if err != nil {
		return fmt.Errorf("failed to create file: %w", err)
	}

	defer file.Close()

	encoder := json.NewEncoder(file)

	encoder.SetIndent("", "    ")

	if err := encoder.Encode(data); err != nil {
		return fmt.Errorf("failed to encode JSON to file: %w", err)
	}

	return nil
}
