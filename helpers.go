package main

import (
	"encoding/base64"
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

// BigIntToBase64 converts a *big.Int* to a Base64 encoded string.
func BigIntToBase64(num *big.Int) string {
	// Get the big-endian byte slice representation of the number.
	bytes := num.Bytes()
	// Encode the byte slice to a Base64 string.
	return base64.StdEncoding.EncodeToString(bytes)
}

// Base64ToBigInt converts a Base64 encoded string back to a *big.Int*.
func Base64ToBigInt(encodedStr string) (*big.Int, error) {
	// Decode the Base64 string into a byte slice.
	decodedBytes, err := base64.StdEncoding.DecodeString(encodedStr)
	if err != nil {
		return nil, fmt.Errorf("error decoding base64 string: %w", err)
	}

	// Create a new big.Int and set its value from the byte slice.
	num := new(big.Int)
	num.SetBytes(decodedBytes)
	return num, nil
}

func TestBase64Conversion(num *big.Int) {
	encoded := BigIntToBase64(num)
	fmt.Println("Encoded:", encoded)
	fmt.Println("Hex:", "0x"+num.Text(16))

	decoded, err := Base64ToBigInt(encoded)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	fmt.Println("Decoded:", decoded)

	// Check if the original number and the decoded number are equal
	if num.Cmp(decoded) == 0 {
		fmt.Println("Success: The original and decoded numbers are equal.")
	} else {
		fmt.Println("Error: The original and decoded numbers are not equal.")
	}
	// Check if the original number and the decoded number are equal
}
