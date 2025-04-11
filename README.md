# ChaINDEX-Go

[![GitHub release (latest SemVer)](https://img.shields.io/github/v/release/SeriousPassenger/ChaINDEX-Go?sort=semver)](https://github.com/SeriousPassenger/ChaINDEX-Go/releases/latest)
[![GitHub all releases](https://img.shields.io/github/downloads/SeriousPassenger/ChaINDEX-Go/total.svg)](https://github.com/SeriousPassenger/ChaINDEX-Go/releases)
[![GitHub Release Date](https://img.shields.io/github/release-date/SeriousPassenger/ChaINDEX-Go)](https://github.com/SeriousPassenger/ChaINDEX-Go/releases/latest)

A simple set of tools for scanning an Ethereum based blockchain using the Ethereum JSON-RPC API.

## TODO

* [X] ~~*Scan and save full block data given a start and end block.*~~ [2025-04-11] 

* [X] ~~*Scan receipt data for given blocks data file*~~ [2025-04-11]

* [X] ~~*Scan all accounts data using debug_accountRange, with the ability to filter by contracts addresses.*~~ [2025-04-11]

* [ ] Make the buffer logic for scanning and saving all account data better, dump to a single file using buffering instead.

* [ ] Batch scan the code of contracts at a specific block.

* [ ] Batch scan the balance of accounts at a specific block.

* [ ] Scan receipt data for given json file of transactions array.