package main

import (
	"encoding/json"
	"fmt"
	"math/big"
	"os"
	"strings"
)

func HexToBigInt(hexString string) (*big.Int, error) {
	hexString = strings.ToLower(hexString)

	if !strings.HasPrefix(hexString, "0x") {
		return nil, fmt.Errorf("invalid hex prefix: %s", hexString)
	}

	sub := hexString[2:]
	value := new(big.Int)

	if _, ok := value.SetString(sub, 16); !ok {
		return nil, fmt.Errorf("invalid hex digits: %s", hexString)
	}

	return value, nil
}

func HexToBigIntMultiple(hexStrings []string) ([]*big.Int, error) {
	values := make([]*big.Int, len(hexStrings))

	for i, hexString := range hexStrings {
		value, err := HexToBigInt(hexString)
		if err != nil {
			return nil, fmt.Errorf("failed to convert hex string %q: %w", hexString, err)
		}
		values[i] = value
	}

	return values, nil
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
