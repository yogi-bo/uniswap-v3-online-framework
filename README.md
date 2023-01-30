# Uniswap v3 Online Strategy Framework and Simulator

## Introduction

This repository contains a framework to develop and back-test online strategies for Uniswap v3 liquidity provision, including:

1. [uniswapv3data](src/uniswapv3data): Python lib to download Uniswap swap history from Google BigQuery.
2. [uniswapv3simulator](src/uniswapv3simulator): Go package to simulate strategies on historical Uniswap data efficiently.
3. [Exponential Weighted Market Maker.ipynb](notebooks/Exponential%20Weighted%20Market%20Maker.ipynb): Jupyter Notebook example using the framework to simulate the Exponential Weights strategy described in [Uniswap Liquidity Provision: An Online Learning Approach](https://arxiv.org/abs/2302.00610).

## Implement your own strategy

Adding another strategy to [uniswapv3simulator](src/uniswapv3simulator) is easy. Implement the following interface in [onlinestrategy](src/uniswapv3simulator/onlinestrategy):
```go
type Strategy interface {
	NextRange(onlineStrategy *OnlineStrategy) (Range, error)
}
```
`NextRange` should return in which price range to provide liquidity on the next step.

Use your implementation instead of `StaticStrategy` in the [main function](src/uniswapv3simulator/main.go).

## Download Uniswap Historical Data

You can use [uniswapv3simulator](src/uniswapv3simulator) to download historical swap data from the Ethereum dataset on Google BitQuery.

First, you should generate a [service account key](https://cloud.google.com/iam/docs/creating-managing-service-account-keys) and download the authentication details to a local json file. Then you can use the lib as follows:
```python
from uniswapv3data.bigquery_downloader import BigQueryUniswapV3Downloader

BIGQUERY_CREDENTIALS_PATH = 'data/bigquery_auth.json'
OUTPUT_FILE_PATH = 'data/uniswap_data.csv.gz'
POOL_ADDRESS = '0x3416cF6C708Da44DB2624D63ea0AAef7113527C6'

BigQueryUniswapV3Downloader(BIGQUERY_CREDENTIALS_PATH).download_pool_data(POOL_ADDRESS, DATE_BEGIN, DATE_END).to_csv(OUTPUT_FILE_PATH)
```

## Notice

This is a very simplified model that does not take into account all effects, such as gas fees. Use at your own risk.
