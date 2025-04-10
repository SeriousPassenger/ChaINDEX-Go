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
		return nil, nil
	}

	sub := hexString[2:]
	value := new(big.Int)

	if _, ok := value.SetString(sub, 16); !ok {
		return nil, nil
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

func SaveStructToJSONFile(data interface{}, path string) error {
	// Create folders in the path if they do not exist
	// Split the path into directories and file name
	parts := strings.Split(path, "/")

	// Create directories if they do not exist
	for i := 0; i < len(parts)-1; i++ {
		dir := strings.Join(parts[:i+1], "/")
		if _, err := os.Stat(dir); os.IsNotExist(err) {
			if err := os.MkdirAll(dir, os.ModePerm); err != nil {
				return fmt.Errorf("failed to create directory %s: %w", dir, err)
			}
		}
	}

	// Save rawData to a file indent 4
	file, err := os.Create(path)

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

func JSONToStruct(path string, data interface{}) error {
	file, err := os.Open(path)

	if err != nil {
		return fmt.Errorf("failed to open file: %w", err)
	}

	defer file.Close()

	decoder := json.NewDecoder(file)

	if err := decoder.Decode(data); err != nil {
		return fmt.Errorf("failed to decode JSON from file: %w", err)
	}

	return nil
}
