package main

import (
	"encoding/json"
	"fmt"
	"log"
	"math"
	"math/big"
	"os"
	"strings"
	"time"
	"unsafe"

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

	batchSize := config.Scan.BlockScanConfig.BatchSize
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

	batchSize := config.Scan.ReceiptScanConfig.BatchSize
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

/*
{"jsonrpc":"2.0","id":1,"result":{"root":"0xf2c2b49d3bcbebb3d26d98fa9d96ac2b62f9ba2b6a3a741eeebe543e27ca89cc","accounts":{"0x0000000000000000000000000000000000efface":{"balance":"1200000000000000","nonce":0,"root":"0x56e81f171bcc55a6ff8345e692c0f86e5b48e01b996cadc001622fb5e363b421","codeHash":"0xc5d2460186f7233c927e7db2dcc703c0e500b653ca82273b7bfad8045d85a470"},"0x0000000000000000000000000000000001664799":{"balance":"400000000000000","nonce":0,"root":"0x56e81f171bcc55a6ff8345e692c0f86e5b48e01b996cadc001622fb5e363b421","codeHash":"0xc5d2460186f7233c927e7db2dcc703c0e500b653ca82273b7bfad8045d85a470"},"0x0000000000000000000000000000000005f5e100":{"balance":"100000000000","nonce":0,"root":"0x56e81f171bcc55a6ff8345e692c0f86e5b48e01b996cadc001622fb5e363b421","codeHash":"0xc5d2460186f7233c927e7db2dcc703c0e500b653ca82273b7bfad8045d85a470"},"0x00000000000000000000000000000000069f6bc7":{"balance":"8000000000000000","nonce":0,"root":"0x56e81f171bcc55a6ff8345e692c0f86e5b48e01b996cadc001622fb5e363b421","codeHash":"0xc5d2460186f7233c927e7db2dcc703c0e500b653ca82273b7bfad8045d85a470"},"0x000000000000000000000000000000000deada02":{"balance":"5000000000000000","nonce":0,"root":"0x56e81f171bcc55a6ff8345e692c0f86e5b48e01b996cadc001622fb5e363b421","codeHash":"0xc5d2460186f7233c927e7db2dcc703c0e500b653ca82273b7bfad8045d85a470"},"0x0000000000000000000000000000000027bc86aa":{"balance":"9000000000000000","nonce":0,"root":"0x56e81f171bcc55a6ff8345e692c0f86e5b48e01b996cadc001622fb5e363b421","codeHash":"0xc5d2460186f7233c927e7db2dcc703c0e500b653ca82273b7bfad8045d85a470"},"0x000000000000000000000000000000004377dead":{"balance":"200000000000000","nonce":0,"root":"0x56e81f171bcc55a6ff8345e692c0f86e5b48e01b996cadc001622fb5e363b421","codeHash":"0xc5d2460186f7233c927e7db2dcc703c0e500b653ca82273b7bfad8045d85a470"},"0x000000000000000000000000000000005349dead":{"balance":"10000000000000","nonce":0,"root":"0x56e81f171bcc55a6ff8345e692c0f86e5b48e01b996cadc001622fb5e363b421","codeHash":"0xc5d2460186f7233c927e7db2dcc703c0e500b653ca82273b7bfad8045d85a470"},"0x0000000000000000000000000000000062d48537":{"balance":"4970580747186068","nonce":0,"root":"0x56e81f171bcc55a6ff8345e692c0f86e5b48e01b996cadc001622fb5e363b421","codeHash":"0xc5d2460186f7233c927e7db2dcc703c0e500b653ca82273b7bfad8045d85a470"},"0x0000000000000000000000000000000062d4b7e8":{"balance":"9138424108689073","nonce":0,"root":"0x56e81f171bcc55a6ff8345e692c0f86e5b48e01b996cadc001622fb5e363b421","codeHash":"0xc5d2460186f7233c927e7db2dcc703c0e500b653ca82273b7bfad8045d85a470"}},"next":"AAAAAAAAAAAAAAAAAAAAAGNyabg="}}
*/

type RpcAccount struct {
	Balance    string `json:"balance"`
	Nonce      uint64 `json:"nonce"`
	Root       string `json:"root"`
	CodeHash   string `json:"codeHash"`
	IsContract bool   `json:"isContract"`
}

type RpcAccountPageResult struct {
	Root     string                `json:"root"`
	Accounts map[string]RpcAccount `json:"accounts"`
	Next     string                `json:"next"`
}

func ScanAllAccounts(config *Config) error {
	if err := validateConfig(config); err != nil {
		return fmt.Errorf("invalid configuration: %w", err)
	}

	log.Printf("Starting the account scanner...\n")

	client, err := GetRpcClient(config.Rpc.Url)
	if err != nil {
		return fmt.Errorf("failed to create RPC client: %w", err)
	}
	log.Printf("Connected to RPC server: %s\n", config.Rpc.Url)

	batchSize := config.Scan.AccountScanConfig.BatchSize
	scanKey := config.Scan.AccountScanConfig.StartKey
	blockNumber := config.Scan.AccountScanConfig.BlockNumber
	maxAccounts := config.Scan.AccountScanConfig.MaxAccounts
	baseOutputFileName := config.Scan.AccountScanConfig.OutputFileName

	if maxAccounts == 0 {
		maxAccounts = math.MaxUint64
	}

	log.Printf("Batch size: %d\n", batchSize)
	log.Printf("Max accounts: %d\n", maxAccounts)
	log.Printf("Start key: %s\n", scanKey)
	log.Printf("Block number: %d\n", blockNumber)
	log.Printf("Output file name: %s\n", baseOutputFileName)
	log.Printf("Delay between requests: %d ms\n", config.Rpc.Delay)
	log.Printf("Starting the account scan...\n")

	// Our in-memory store for accounts.
	accounts := make(map[string]RpcAccount)
	nonContractCodeHash := "0xc5d2460186f7233c927e7db2dcc703c0e500b653ca82273b7bfad8045d85a470"
	numContracts := 0
	numNonContracts := 0
	numContractsWithBalance := 0

	// Set up flush threshold and a slice to track chunk files.
	const flushThresholdBytes = 10 * 1024 * 1024 // 10 MB
	chunkFiles := []string{}

	// flushAccounts flushes the current accounts map to disk and resets it.
	flushAccounts := func() error {
		chunkFileName := fmt.Sprintf("%s_%d.json", baseOutputFileName, len(chunkFiles)+1)
		filePath := fmt.Sprintf("%s/%s", config.Scan.OutputDir, chunkFileName)
		if err := SaveStructToJSONFile(accounts, filePath); err != nil {
			return fmt.Errorf("failed to save accounts chunk to file: %w", err)
		}
		memoryUsageMB := (len(accounts) * int(unsafe.Sizeof(RpcAccount{}))) / (1024 * 1024)
		log.Printf("Flushed %d accounts (approx %d MB) to file: %s\n", len(accounts), memoryUsageMB, filePath)
		chunkFiles = append(chunkFiles, filePath)
		// Reset the in-memory store.
		accounts = make(map[string]RpcAccount)
		return nil
	}

	totalAccounts := 0

	latestAccount := &RpcAccount{
		Balance:    "0",
		Nonce:      0,
		Root:       "",
		CodeHash:   "",
		IsContract: false,
	}

	latestAddress := "0x"

	iters := 0

	lastTime := time.Now()

	lastAccountCount := 1

	seenAccounts := make(map[string]bool)

	// Main scanning loop.
	for {
		iters++
		if uint64(len(accounts)) >= maxAccounts {
			log.Printf("Reached the maximum number of accounts to scan: %d\n", maxAccounts)
			break
		}

		var raw json.RawMessage
		err = client.Call(&raw, "debug_accountRange", blockNumber, scanKey, batchSize, true, true)
		if err != nil {
			return fmt.Errorf("failed to fetch accounts: %w", err)
		}

		// Delay between requests.
		time.Sleep(time.Duration(config.Rpc.Delay) * time.Millisecond)

		var result RpcAccountPageResult
		if err := json.Unmarshal(raw, &result); err != nil {
			return fmt.Errorf("failed to unmarshal JSON: %w", err)
		}

		// Update the scan key.
		newKey := result.Next
		scanKey = newKey

		completedScan := false

		// Process each account in the batch.
		for address, account := range result.Accounts {
			if seenAccounts[address] {
				log.Printf("Completely scanned all accounts !\n")
				completedScan = true
			}

			seenAccounts[address] = true

			if uint64(len(seenAccounts)) >= maxAccounts {
				log.Printf("Reached the maximum number of accounts to scan: %d !\n", maxAccounts)
				completedScan = true
			}

			if account.CodeHash == nonContractCodeHash {
				numNonContracts++
				account.IsContract = false
			} else {
				numContracts++
				account.IsContract = true
				if account.Balance != "0" {
					numContractsWithBalance++

					accounts[address] = account

					latestAccount = &account
					latestAddress = address
				}

			}

		}

		if completedScan {
			log.Printf("Completed scan!\n")
			break
		}

		if iters%20 == 0 {
			log.Printf("===========================================\n")

			log.Printf("Total: %d accounts (%d contracts, %d non-contracts, %d contracts with balance)\n",
				len(accounts), numContracts, numNonContracts, numContractsWithBalance)
			log.Printf("===========================================\n")
			scanKeyBigInt, _ := Base64ToBigInt(scanKey)
			log.Printf("Current scan key: %s (%d)\n", scanKey, scanKeyBigInt)
			log.Printf("===========================================\n")
			log.Printf("Latest account item: \n")
			log.Printf("Address: %s\n", latestAddress)
			log.Printf("Balance: %s\n", latestAccount.Balance)
			log.Printf("Nonce: %d\n", latestAccount.Nonce)
			log.Printf("===========================================\n")
			log.Printf("Time elapsed: %s\n", time.Since(lastTime))

			accountsIncreased := len(accounts) - lastAccountCount
			lastAccountCount = len(accounts)
			accountPerMillisec := float64(accountsIncreased) / float64(time.Since(lastTime).Milliseconds())
			log.Printf("Accounts per second: %.2f\n", accountPerMillisec*1000)

			log.Printf("===========================================\n")
			lastTime = time.Now()
		}

		// Flush if the estimated in-memory size exceeds the threshold.
		currMemoryBytes := len(accounts) * int(unsafe.Sizeof(RpcAccount{}))
		if currMemoryBytes >= flushThresholdBytes {
			totalAccounts += len(accounts)
			if err := flushAccounts(); err != nil {
				return err
			}
		}
	}

	// Flush any remaining accounts.
	if len(accounts) > 0 {
		totalAccounts += len(accounts)
		chunkFileName := fmt.Sprintf("%s_%d.json", baseOutputFileName, len(chunkFiles)+1)
		filePath := fmt.Sprintf("%s/%s", config.Scan.OutputDir, chunkFileName)
		if err := SaveStructToJSONFile(accounts, filePath); err != nil {
			return fmt.Errorf("failed to save remaining accounts to file: %w", err)
		}
		log.Printf("Flushed final %d accounts to file: %s\n", len(accounts), filePath)
		chunkFiles = append(chunkFiles, filePath)
	}

	// At this point, multiple chunk files have been created.
	// Now merge them into one final file.
	if len(chunkFiles) > 0 {
		mergedAccounts := make(map[string]RpcAccount)
		for _, chunkFile := range chunkFiles {
			data, err := os.ReadFile(chunkFile)
			if err != nil {
				return fmt.Errorf("failed to read chunk file %s: %w", chunkFile, err)
			}
			var chunkData map[string]RpcAccount
			if err := json.Unmarshal(data, &chunkData); err != nil {
				return fmt.Errorf("failed to unmarshal chunk file %s: %w", chunkFile, err)
			}
			for k, v := range chunkData {
				mergedAccounts[k] = v
			}
		}
		mergedFile := fmt.Sprintf("%s/%s_merged.json", config.Scan.OutputDir, baseOutputFileName)
		if err := SaveStructToJSONFile(mergedAccounts, mergedFile); err != nil {
			return fmt.Errorf("failed to save merged accounts to file: %w", err)
		}
		log.Printf("Merged %d chunk files into final file: %s\n", len(chunkFiles), mergedFile)
	} else {
		// No chunk files were written so far. This means no flush was needed,
		// and you already have the complete data in memory.
		mergedFile := fmt.Sprintf("%s/%s.json", config.Scan.OutputDir, baseOutputFileName)
		if err := SaveStructToJSONFile(accounts, mergedFile); err != nil {
			return fmt.Errorf("failed to save accounts to file: %w", err)
		}
		log.Printf("Accounts saved to final file: %s\n", mergedFile)
	}

	log.Printf("\nTask completed successfully!\n")
	log.Printf("Total contracts: %d\n", numContracts)
	log.Printf("Total non-contracts: %d\n", numNonContracts)
	log.Printf("Last scan key: %s\n", scanKey)
	return nil
}
