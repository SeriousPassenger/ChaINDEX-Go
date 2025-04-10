package main

import (
	"fmt"
	"log"
	"math/big"
	"strings"
	"time"

	"github.com/schollz/progressbar/v3"
)

// -----------------------------------------
// Complete scanning logic
// -----------------------------------------

func validateConfig(config *Config) error {
	if config.Scan.FromBlock > config.Scan.ToBlock {
		return fmt.Errorf("from_block must be less or equal to to_block")
	}

	return nil
}

func ScanBlocksWithConfig(config *Config) error {
	if err := validateConfig(config); err != nil {
		return fmt.Errorf("invalid configuration: %w", err)
	}

	log.Printf("Starting the block scanner...\n")

	startBlock := config.Scan.FromBlock
	endBlock := config.Scan.ToBlock

	batchSize := config.Rpc.BatchSize
	totalBlocks := endBlock - startBlock + 1

	batchCount := totalBlocks / batchSize

	if totalBlocks%batchSize != 0 {
		batchCount++
	}

	log.Printf("Total blocks to scan: %d\n", totalBlocks)
	log.Printf("Batch size: %d\n", batchSize)
	log.Printf("Delay between requests: %d ms\n", config.Rpc.Delay)

	batches := make([][]*big.Int, batchCount)

	for i := uint64(0); i < batchCount; i++ {
		batchStart := startBlock + uint64(i)*batchSize
		batchEnd := batchStart + batchSize - 1

		if batchEnd > endBlock {
			batchEnd = endBlock
		}

		currentBatchBlockCount := batchEnd - batchStart + 1
		currentBatch := make([]*big.Int, currentBatchBlockCount)

		for j := uint64(0); j < currentBatchBlockCount; j++ {
			currentBatch[j] = big.NewInt(0).SetUint64(batchStart + j)
		}

		batches[i] = currentBatch
	}

	log.Printf("Total batches created: %d\n", len(batches))

	client, err := GetRpcClient(config.Rpc.Url)

	if err != nil {
		return fmt.Errorf("failed to create RPC client: %w", err)
	}

	log.Printf("Connected to RPC server: %s\n", config.Rpc.Url)

	bar := progressbar.NewOptions64(int64(batchCount),
		progressbar.OptionSetDescription("Fetching blocks..."),
		progressbar.OptionSetWriter(log.Writer()),
		progressbar.OptionSetWidth(20),
	)

	blocks := make([]RpcBlockFull, 0)

	log.Printf("Fetching blocks from %d to %d in batches of %d\n", startBlock, endBlock, batchSize)

	txCount := 0
	failedBlockCount := 0
	successfulBlockCount := 0

	for _, batch := range batches {
		if len(batch) == 0 {
			continue
		}

		bar.Describe(fmt.Sprintf("Txs: %d, Success (Block): %d Fail (Block): %d", txCount, successfulBlockCount, failedBlockCount))

		batchBlocks, err := GetBlocksBatch(client, batch)

		// Add a delay between requests
		time.Sleep(time.Duration(config.Rpc.Delay) * time.Millisecond)

		if err != nil {
			bar.Add(1)
			failedBlockCount += len(batch)
			return fmt.Errorf("failed to fetch blocks: %w", err)
		}

		for _, block := range batchBlocks {
			transactions := block.Transactions

			if len(config.Filter.Addresses) > 0 {
				filteredTransactions := make([]RpcTransactionFull, 0)

				for _, tx := range transactions {
					fromLower := strings.ToLower(tx.From)
					toLower := strings.ToLower(tx.To)

					for _, address := range config.Filter.Addresses {
						addressLower := strings.ToLower(address)

						if address == "ContractCreation" && tx.To == "" {
							filteredTransactions = append(filteredTransactions, tx)
							break
						}

						if fromLower == addressLower || toLower == addressLower {
							filteredTransactions = append(filteredTransactions, tx)
							break
						}

					}

				}
				block.Transactions = filteredTransactions

			}

			txCount += len(block.Transactions)

			successfulBlockCount++
			blocks = append(blocks, block)
		}

		bar.Add(1)

	}

	if len(blocks) == 0 {
		return fmt.Errorf("no blocks fetched")
	}

	outBlocks := make([]*BlockFull, len(blocks))

	for i, block := range blocks {
		blockData, err := RpcBlockFullToBlockFull(&block)

		if err != nil {
			return fmt.Errorf("failed to convert block data: %w", err)
		}

		outBlocks[i] = blockData
	}

	filePath := fmt.Sprintf("%s/blocks_%d_to_%d.json", config.Scan.OutputDir, startBlock, endBlock)

	if err := SaveStructToJSONFile(outBlocks, filePath); err != nil {
		return fmt.Errorf("failed to save blocks to file: %w", err)
	}

	log.Printf("\nTask completed successfully!\n")
	log.Printf("Blocks fetched and saved to %s.\n", filePath)

	bar.Finish()
	return nil
}

func ScanReceiptsWithConfig(config *Config, blockFile string) error {
	var blocks []*BlockFull

	err := JSONToStruct(blockFile, &blocks)

	if err != nil {
		return fmt.Errorf("failed to load blocks from file: %w", err)
	}

	transactions := make([]string, 0)

	for _, block := range blocks {
		if block == nil {
			continue
		}

		for _, tx := range block.Transactions {
			transactions = append(transactions, tx.Hash)
		}
	}

	if err := validateConfig(config); err != nil {
		return fmt.Errorf("invalid configuration: %w", err)
	}

	log.Printf("Starting the receipt scanner...\n")

	client, err := GetRpcClient(config.Rpc.Url)

	if err != nil {
		return fmt.Errorf("failed to create RPC client: %w", err)
	}

	log.Printf("Connected to RPC server: %s\n", config.Rpc.Url)

	batchSize := config.Rpc.BatchSize
	totalTransactions := uint64(len(transactions))
	batchCount := totalTransactions / batchSize

	if totalTransactions%batchSize != 0 {
		batchCount++
	}

	log.Printf("Total transactions to scan: %d\n", totalTransactions)

	log.Printf("Batch size: %d\n", batchSize)

	bar := progressbar.NewOptions64(int64(batchCount),
		progressbar.OptionSetDescription("Fetching receipts..."),
		progressbar.OptionSetWriter(log.Writer()),
		progressbar.OptionSetWidth(20),
	)

	receipts := make([]Receipt, 0)

	for i := uint64(0); i < batchCount; i++ {
		batchStart := i * batchSize
		batchEnd := batchStart + batchSize

		if batchEnd > totalTransactions {
			batchEnd = totalTransactions
		}

		currentBatch := transactions[batchStart:batchEnd]

		bar.Describe(fmt.Sprintf("Fetching receipts for transactions %d to %d", batchStart, batchEnd))

		batchReceipts, err := GetReceiptsBatch(client, currentBatch)
		// Add a delay between requests
		time.Sleep(time.Duration(config.Rpc.Delay) * time.Millisecond)

		if err != nil {
			bar.Add(1)
			return fmt.Errorf("failed to fetch receipts: %w", err)
		}

		for _, receipt := range batchReceipts {
			receiptData, err := RpcReceiptToReceipt(&receipt)

			if err != nil {
				bar.Add(1)
				return fmt.Errorf("failed to convert receipt data: %w", err)
			}

			receipts = append(receipts, *receiptData)
		}

		bar.Add(1)
	}

	if len(receipts) == 0 {
		return fmt.Errorf("no receipts fetched")
	}

	filePath := fmt.Sprintf("%s/receipts_%d_to_%d.json", config.Scan.OutputDir, config.Scan.FromBlock, config.Scan.ToBlock)

	if err := SaveStructToJSONFile(receipts, filePath); err != nil {
		return fmt.Errorf("failed to save receipts to file: %w", err)
	}

	log.Printf("\nTask completed successfully!\n")
	log.Printf("Receipts fetched and saved to %s.\n", filePath)
	bar.Finish()

	return nil
}
