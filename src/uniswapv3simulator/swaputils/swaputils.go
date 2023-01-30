package swaputils

import (
	"compress/gzip"
	"encoding/csv"
	"errors"
	"os"
	"strconv"
	"time"
)

type Swap struct {
	BlockNumber      uint64
	TransactionIndex uint64
	LogIndex         uint64
	BlockTimestamp   time.Time
	Amount0          float64
	Amount1          float64
	Price            float64
	Liquidity        float64
	Tick             int32
}

func parseSwapCsvLine(line []string) (swap Swap, err error) {
	if len(line) != 11 {
		err = errors.New("CSV line length is wrong")
		return
	}

	swap.BlockNumber, err = strconv.ParseUint(line[1], 0, 64)
	if err != nil {
		return
	}

	swap.TransactionIndex, err = strconv.ParseUint(line[2], 0, 64)
	if err != nil {
		return
	}

	swap.LogIndex, err = strconv.ParseUint(line[3], 0, 64)
	if err != nil {
		return
	}

	swap.BlockTimestamp, err = time.Parse("2006-01-02 15:04:05-07:00", line[4])
	if err != nil {
		return
	}

	swap.Amount0, err = strconv.ParseFloat(line[6], 64)
	if err != nil {
		return
	}

	swap.Amount1, err = strconv.ParseFloat(line[7], 64)
	if err != nil {
		return
	}

	swap.Price, err = strconv.ParseFloat(line[8], 64)
	if err != nil {
		return
	}

	swap.Liquidity, err = strconv.ParseFloat(line[9], 64)
	if err != nil {
		return
	}

	tick64, err := strconv.ParseFloat(line[10], 64)
	if err != nil {
		return
	}
	swap.Tick = int32(tick64)

	return
}

func ReadSwapsFile(filePath string) ([]Swap, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	gzipReader, err := gzip.NewReader(file)
	if err != nil {
		return nil, err
	}
	defer gzipReader.Close()

	swapData, err := csv.NewReader(gzipReader).ReadAll()
	if err != nil {
		return nil, err
	}
	swapData = swapData[1:] // remove headers

	swaps := make([]Swap, len(swapData))
	for i, line := range swapData {
		swaps[i], err = parseSwapCsvLine(line)
		if err != nil {
			return nil, err
		}
	}

	return swaps, nil
}
