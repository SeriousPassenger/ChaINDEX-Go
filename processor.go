package main

import (
	"fmt"
)

// ------------------------------------------------------
// Convert RPC structures to internal structures
// ------------------------------------------------------

func RpcTransactionToTransaction(tx *RpcTransaction) (*Transaction, error) {

	inputStrings := []string{
		tx.BlockNumber,
		tx.Value,
		tx.Gas,
		tx.GasPrice,
		tx.Nonce,
		tx.TransactionIndex,
		tx.ChainId,
		tx.Type,
	}

	bigIntValues, err := HexToBigIntMultiple(inputStrings)

	if err != nil {
		return nil, fmt.Errorf("failed to convert hex strings to big.Int: %w", err)
	}

	blockNumberInt, txValueInt, txGasInt, txGasPriceInt, txNonceInt, txIndexInt, chainIdInt, txTypeInt := bigIntValues[0], bigIntValues[1], bigIntValues[2], bigIntValues[3], bigIntValues[4], bigIntValues[5], bigIntValues[6], bigIntValues[7]

	txOut := &Transaction{
		BlockHash:        tx.BlockHash,
		BlockNumber:      blockNumberInt,
		From:             tx.From,
		Gas:              txGasInt,
		GasPrice:         txGasPriceInt,
		Hash:             tx.Hash,
		Input:            tx.Input,
		Nonce:            txNonceInt,
		To:               tx.To,
		TransactionIndex: txIndexInt,
		Value:            txValueInt,
		Type:             txTypeInt,
		ChainId:          chainIdInt,
		V:                tx.V,
		R:                tx.R,
		S:                tx.S,
	}

	return txOut, nil
}

func RpcLogToLog(log *RpcLog) (*Log, error) {

	inputHexStrings := []string{
		log.BlockNumber,
		log.TransactionIndex,
		log.LogIndex,
	}

	bigIntValues, err := HexToBigIntMultiple(inputHexStrings)

	if err != nil {
		return nil, fmt.Errorf("failed to convert hex strings to big.Int: %w", err)
	}

	blockNumberInt, transactionIndexInt, logIndexInt := bigIntValues[0], bigIntValues[1], bigIntValues[2]

	logOut := &Log{
		Address:          log.Address,
		Topics:           log.Topics,
		Data:             log.Data,
		BlockNumber:      blockNumberInt,
		TransactionHash:  log.TransactionHash,
		TransactionIndex: transactionIndexInt,
		BlockHash:        log.BlockHash,
		LogIndex:         logIndexInt,
		Removed:          log.Removed,
	}

	return logOut, nil
}

func RpcReceiptToReceipt(receipt *RpcReceipt) (*Receipt, error) {
	inputHexStrings := []string{
		receipt.BlockNumber,
		receipt.CumulativeGasUsed,
		receipt.GasUsed,
		receipt.TransactionIndex,
		receipt.EffectiveGasPrice,
	}

	bigIntValues, err := HexToBigIntMultiple(inputHexStrings)

	if err != nil {
		return nil, fmt.Errorf("failed to convert hex strings to big.Int: %w", err)
	}

	blockNumberInt, cumulativeGasUsedInt, gasUsedInt, transactionIndexInt, effectiveGasPriceInt := bigIntValues[0], bigIntValues[1], bigIntValues[2], bigIntValues[3], bigIntValues[4]

	// Convert the logs to the new format
	logs := make([]Log, len(receipt.Logs))

	for i, log := range receipt.Logs {
		logResult, err := RpcLogToLog(&log)

		if err != nil {
			return nil, fmt.Errorf("failed to convert log: %w", err)
		}

		logs[i] = *logResult
	}

	receiptOut := &Receipt{
		BlockHash:         receipt.BlockHash,
		BlockNumber:       blockNumberInt,
		ContractAddress:   receipt.ContractAddress,
		CumulativeGasUsed: cumulativeGasUsedInt,
		GasUsed:           gasUsedInt,
		Status:            receipt.Status,
		To:                receipt.To,
		TransactionHash:   receipt.TransactionHash,
		TransactionIndex:  transactionIndexInt,
		Logs:              logs,
		LogsBloom:         receipt.LogsBloom,
		From:              receipt.From,
		EffectiveGasPrice: effectiveGasPriceInt,
		Type:              receipt.Type,
	}

	return receiptOut, nil
}

func RpcBlockToBlockAndTransactions(block *RpcBlock) (*BaseBlock, []*Transaction, error) {
	hexStrings := []string{
		block.Difficulty,
		block.GasLimit,
		block.GasUsed,
		block.Nonce,
		block.Number,
		block.Size,
		block.Timestamp,
		block.TotalDifficulty,
	}

	bigIntValues, err := HexToBigIntMultiple(hexStrings)

	if err != nil {
		return nil, nil, fmt.Errorf("failed to convert hex strings to big.Int: %w", err)
	}

	diffInt, gasLimitInt, gasUsedInt, nonceInt, numberInt, sizeInt, timestampInt, totalDiffInt := bigIntValues[0], bigIntValues[1], bigIntValues[2], bigIntValues[3], bigIntValues[4], bigIntValues[5], bigIntValues[6], bigIntValues[7]

	baseBlock := &BaseBlock{
		Difficulty:       diffInt,
		ExtraData:        block.ExtraData,
		GasLimit:         gasLimitInt,
		GasUsed:          gasUsedInt,
		Hash:             block.Hash,
		LogsBloom:        block.LogsBloom,
		Miner:            block.Miner,
		MixHash:          block.MixHash,
		Nonce:            nonceInt,
		Number:           numberInt,
		ParentHash:       block.ParentHash,
		ReceiptsRoot:     block.ReceiptsRoot,
		Sha3Uncles:       block.Sha3Uncles,
		Size:             sizeInt,
		StateRoot:        block.StateRoot,
		Timestamp:        timestampInt,
		TotalDifficulty:  totalDiffInt,
		TransactionsRoot: block.TransactionsRoot,
		Uncles:           block.Uncles,
	}

	transactions := make([]*Transaction, len(block.Transactions))

	for i, tx := range block.Transactions {
		outTX, err := RpcTransactionToTransaction(&tx)

		if err != nil {
			return nil, nil, fmt.Errorf("failed to convert transaction: %w", err)
		}

		transactions[i] = outTX
	}

	return baseBlock, transactions, nil
}

// ------------------------------------------------------
// Combine internal structures into full structures
// ------------------------------------------------------

func TransactionToFullTransaction(tx *Transaction, receipt *Receipt) *FullTransaction {
	return &FullTransaction{
		BaseTransaction: *tx,
		Receipt:         *receipt,
	}
}

func GetFullBlock(block *BaseBlock, transactions []*Transaction, receipts []*Receipt) *FullBlock {
	fullTransactions := make([]FullTransaction, len(transactions))

	for i, tx := range transactions {
		fullTransactions[i] = *TransactionToFullTransaction(tx, receipts[i])
	}

	return &FullBlock{
		BaseBlock:        *block,
		FullTransactions: fullTransactions,
	}
}
