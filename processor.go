package main

import (
	"fmt"
	"math/big"
)

/*
type RpcBlockMinimal struct {
	Difficulty       string   `json:"difficulty"`
	ExtraData        string   `json:"extraData"`
	GasLimit         string   `json:"gasLimit"`
	GasUsed          string   `json:"gasUsed"`
	Hash             string   `json:"hash"`
	LogsBloom        string   `json:"logsBloom"`
	Miner            string   `json:"miner"`
	MixHash          string   `json:"mixHash"`
	Nonce            string   `json:"nonce"`
	Number           string   `json:"number"`
	ParentHash       string   `json:"parentHash"`
	ReceiptsRoot     string   `json:"receiptsRoot"`
	Sha3Uncles       string   `json:"sha3Uncles"`
	Size             string   `json:"size"`
	StateRoot        string   `json:"stateRoot"`
	Timestamp        string   `json:"timestamp"`
	TotalDifficulty  string   `json:"totalDifficulty"`
	Transactions     []string `json:"transactions"`
	TransactionsRoot string   `json:"transactionsRoot"`
	Uncles           []string `json:"uncles"`
}

*/

func RpcBlockMinimalToBlockMinimal(rpcBlock *RpcBlockMinimal) (*BlockMinimal, error) {
	hexStrings := []string{
		rpcBlock.Difficulty,
		rpcBlock.GasLimit,
		rpcBlock.GasUsed,
		rpcBlock.Nonce,
		rpcBlock.Number,
		rpcBlock.Size,
		rpcBlock.Timestamp,
		rpcBlock.TotalDifficulty,
	}

	// Convert hex strings to big.Int
	values, err := HexToBigIntMultiple(hexStrings)

	if err != nil {
		return nil, fmt.Errorf("failed to convert hex strings to big.Int: %w", err)
	}

	// Create a map to associate hex strings with their corresponding big.Int values
	hexToBigIntMap := map[string]*big.Int{
		rpcBlock.Difficulty:      values[0],
		rpcBlock.GasLimit:        values[1],
		rpcBlock.GasUsed:         values[2],
		rpcBlock.Nonce:           values[3],
		rpcBlock.Number:          values[4],
		rpcBlock.Size:            values[5],
		rpcBlock.Timestamp:       values[6],
		rpcBlock.TotalDifficulty: values[7],
	}

	block := &BlockMinimal{
		Timestamp:        hexToBigIntMap[rpcBlock.Timestamp],
		Difficulty:       hexToBigIntMap[rpcBlock.Difficulty],
		ExtraData:        rpcBlock.ExtraData,
		GasLimit:         hexToBigIntMap[rpcBlock.GasLimit],
		GasUsed:          hexToBigIntMap[rpcBlock.GasUsed],
		Hash:             rpcBlock.Hash,
		LogsBloom:        rpcBlock.LogsBloom,
		Miner:            rpcBlock.Miner,
		MixHash:          rpcBlock.MixHash,
		Nonce:            hexToBigIntMap[rpcBlock.Nonce],
		Number:           hexToBigIntMap[rpcBlock.Number],
		ParentHash:       rpcBlock.ParentHash,
		ReceiptsRoot:     rpcBlock.ReceiptsRoot,
		Sha3Uncles:       rpcBlock.Sha3Uncles,
		Size:             hexToBigIntMap[rpcBlock.Size],
		StateRoot:        rpcBlock.StateRoot,
		TotalDifficulty:  hexToBigIntMap[rpcBlock.TotalDifficulty],
		TransactionsRoot: rpcBlock.TransactionsRoot,
		Uncles:           rpcBlock.Uncles,
		Transactions:     rpcBlock.Transactions,
	}

	return block, nil
}

func RpcBlockFullToBlockFull(rpcBlock *RpcBlockFull) (*BlockFull, error) {
	hexStrings := []string{
		rpcBlock.Difficulty,
		rpcBlock.GasLimit,
		rpcBlock.GasUsed,
		rpcBlock.Nonce,
		rpcBlock.Number,
		rpcBlock.Size,
		rpcBlock.Timestamp,
		rpcBlock.TotalDifficulty,
	}

	// Convert hex strings to big.Int
	values, err := HexToBigIntMultiple(hexStrings)

	if err != nil {
		return nil, fmt.Errorf("failed to convert hex strings to big.Int: %w", err)
	}

	// Create a map to associate hex strings with their corresponding big.Int values
	hexToBigIntMap := map[string]*big.Int{
		rpcBlock.Difficulty:      values[0],
		rpcBlock.GasLimit:        values[1],
		rpcBlock.GasUsed:         values[2],
		rpcBlock.Nonce:           values[3],
		rpcBlock.Number:          values[4],
		rpcBlock.Size:            values[5],
		rpcBlock.Timestamp:       values[6],
		rpcBlock.TotalDifficulty: values[7],
	}

	transactions := make([]TransactionFull, len(rpcBlock.Transactions))

	for i, tx := range rpcBlock.Transactions {
		hexStrings := []string{
			tx.BlockNumber,
			tx.Gas,
			tx.GasPrice,
			tx.Nonce,
			tx.ChainId,
			tx.Type,
			tx.TransactionIndex,
			tx.Value,
		}

		// Convert hex strings to big.Int
		values, err := HexToBigIntMultiple(hexStrings)

		if err != nil {
			return nil, fmt.Errorf("failed to convert hex strings to big.Int: %w", err)
		}

		// Create a map to associate hex strings with their corresponding big.Int values
		hexToBigIntMap := map[string]*big.Int{
			tx.BlockNumber:      values[0],
			tx.Gas:              values[1],
			tx.GasPrice:         values[2],
			tx.Nonce:            values[3],
			tx.ChainId:          values[4],
			tx.Type:             values[5],
			tx.TransactionIndex: values[6],
			tx.Value:            values[7],
		}

		transactions[i] = TransactionFull{
			BlockHash:        tx.BlockHash,
			BlockNumber:      hexToBigIntMap[tx.BlockNumber],
			From:             tx.From,
			Gas:              hexToBigIntMap[tx.Gas],
			GasPrice:         hexToBigIntMap[tx.GasPrice],
			Hash:             tx.Hash,
			Input:            tx.Input,
			Nonce:            hexToBigIntMap[tx.Nonce],
			To:               tx.To,
			TransactionIndex: hexToBigIntMap[tx.TransactionIndex],
			Value:            hexToBigIntMap[tx.Value],
			Type:             hexToBigIntMap[tx.Type],
			ChainId:          hexToBigIntMap[tx.ChainId],
			V:                tx.V,
			R:                tx.R,
			S:                tx.S,
		}

	}

	block := &BlockFull{
		Timestamp:        hexToBigIntMap[rpcBlock.Timestamp],
		Difficulty:       hexToBigIntMap[rpcBlock.Difficulty],
		ExtraData:        rpcBlock.ExtraData,
		GasLimit:         hexToBigIntMap[rpcBlock.GasLimit],
		GasUsed:          hexToBigIntMap[rpcBlock.GasUsed],
		Hash:             rpcBlock.Hash,
		LogsBloom:        rpcBlock.LogsBloom,
		Miner:            rpcBlock.Miner,
		MixHash:          rpcBlock.MixHash,
		Nonce:            hexToBigIntMap[rpcBlock.Nonce],
		Number:           hexToBigIntMap[rpcBlock.Number],
		ParentHash:       rpcBlock.ParentHash,
		ReceiptsRoot:     rpcBlock.ReceiptsRoot,
		Sha3Uncles:       rpcBlock.Sha3Uncles,
		Size:             hexToBigIntMap[rpcBlock.Size],
		StateRoot:        rpcBlock.StateRoot,
		TotalDifficulty:  hexToBigIntMap[rpcBlock.TotalDifficulty],
		TransactionsRoot: rpcBlock.TransactionsRoot,
		Uncles:           rpcBlock.Uncles,
		Transactions:     transactions,
	}

	return block, nil
}

func RpcReceiptToReceipt(rpcReceipt *RpcReceipt) (*Receipt, error) {
	hexStrings := []string{
		rpcReceipt.BlockNumber,
		rpcReceipt.CumulativeGasUsed,
		rpcReceipt.EffectiveGasPrice,
		rpcReceipt.GasUsed,
		rpcReceipt.TransactionIndex,
		rpcReceipt.ContractAddress,
	}

	// Convert hex strings to big.Int
	values, err := HexToBigIntMultiple(hexStrings)

	if err != nil {
		return nil, fmt.Errorf("failed to convert hex strings to big.Int: %w", err)
	}

	// Create a map to associate hex strings with their corresponding big.Int values
	hexToBigIntMap := map[string]*big.Int{
		rpcReceipt.BlockNumber:       values[0],
		rpcReceipt.CumulativeGasUsed: values[1],
		rpcReceipt.EffectiveGasPrice: values[2],
		rpcReceipt.GasUsed:           values[3],
		rpcReceipt.TransactionIndex:  values[4],
		rpcReceipt.Type:              values[5],
	}

	logs := make([]Log, len(rpcReceipt.Logs))

	for i, log := range rpcReceipt.Logs {
		hexStrings := []string{
			log.LogIndex,
		}

		// Convert hex strings to big.Int
		values, err := HexToBigIntMultiple(hexStrings)

		if err != nil {
			return nil, fmt.Errorf("failed to convert hex strings to big.Int: %w", err)
		}

		hexToBigIntMap[log.LogIndex] = values[0]

		logs[i] = Log{
			Address:          log.Address,
			Topics:           log.Topics,
			Data:             log.Data,
			BlockNumber:      hexToBigIntMap[log.BlockNumber],
			TransactionHash:  log.TransactionHash,
			TransactionIndex: hexToBigIntMap[log.TransactionIndex],
			BlockHash:        log.BlockHash,
			LogIndex:         hexToBigIntMap[log.LogIndex],
			Removed:          log.Removed,
		}
	}

	return &Receipt{
		BlockHash:         rpcReceipt.BlockHash,
		BlockNumber:       hexToBigIntMap[rpcReceipt.BlockNumber],
		ContractAddress:   rpcReceipt.ContractAddress,
		CumulativeGasUsed: hexToBigIntMap[rpcReceipt.CumulativeGasUsed],
		GasUsed:           hexToBigIntMap[rpcReceipt.GasUsed],
		Status:            rpcReceipt.Status,
		To:                rpcReceipt.To,
		TransactionHash:   rpcReceipt.TransactionHash,
		TransactionIndex:  hexToBigIntMap[rpcReceipt.TransactionIndex],
		Logs:              logs,
		LogsBloom:         rpcReceipt.LogsBloom,
		From:              rpcReceipt.From,
		EffectiveGasPrice: hexToBigIntMap[rpcReceipt.EffectiveGasPrice],
		Type:              hexToBigIntMap[rpcReceipt.Type],
	}, nil
}
