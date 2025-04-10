# ChaINDEX-Go

[![GitHub release (latest SemVer)](https://img.shields.io/github/v/release/SeriousPassenger/ChaINDEX-Go?sort=semver)](https://github.com/SeriousPassenger/ChaINDEX-Go/releases/latest)
[![GitHub all releases](https://img.shields.io/github/downloads/SeriousPassenger/ChaINDEX-Go/total.svg)](https://github.com/SeriousPassenger/ChaINDEX-Go/releases)
[![GitHub Release Date](https://img.shields.io/github/release-date/SeriousPassenger/ChaINDEX-Go)](https://github.com/SeriousPassenger/ChaINDEX-Go/releases/latest)

A simple set of tools for scanning an Ethereum based blockchain using the Ethereum JSON-RPC API.

## TODO

* [ ] Scan and save block data given a start and end block.  (with/without full transactions)
* [ ] Scan and save receipts given a blocks data file.
* [ ] Extract contract creations from a given receipts data file.
* [ ] Extract all unique accounts from a given full transactions block data file.
* [ ] Scan the balances of a given set of accounts (including contracts), for given blocks (for historical analysis of balance) or the latest block.
* [ ] Scan the contract codes of given set of contract addresses.
* [ ] Targeting/filtering a set of addresses.
* [ ] Create a data structure that combines all of the data with ability to export.
* [ ] Resuming/updating the data with the latest new data.