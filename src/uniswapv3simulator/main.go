package main

import (
	"log"
	"os"
	"strconv"
	"time"
	"uniswapv3simulator/onlinestrategy"
	"uniswapv3simulator/swaputils"
)

func parseArgs(args []string) (filePath string, feeTier float64, tickSpacing uint32, concentrationController uint32, err error) {
	filePath = os.Args[1]

	feeTier, err = strconv.ParseFloat(os.Args[2], 64)
	if err != nil {
		return
	}

	tickSpacing64, err := strconv.ParseUint(os.Args[3], 0, 32)
	if err != nil {
		return
	}
	tickSpacing = uint32(tickSpacing64)

	concentrationController64, err := strconv.ParseUint(os.Args[4], 0, 32)
	if err != nil {
		return
	}
	concentrationController = uint32(concentrationController64)

	return
}

func main() {
	filePath, feeTier, tickSpacing, concentrationController, err := parseArgs(os.Args)
	if err != nil {
		log.Fatal(err)
	}

	swaps, err := swaputils.ReadSwapsFile(filePath)
	if err != nil {
		log.Fatal(err)
	}

	if len(swaps) <= 0 {
		log.Fatal("No swaps provided")
	}

	strategy := onlinestrategy.OnlineStrategy{
		UninvestedToken0: 0.5,
		UninvestedToken1: 0.5 * swaps[0].Price,
		ActiveRange:      onlinestrategy.Range{TickSpacing: tickSpacing},
		LastSwap:         swaps[0],
		Strategy: &onlinestrategy.StaticStrategy{
			ConcentrationController: concentrationController,
			Duration:                time.Hour,
			LastDeposit:             time.Unix(0, 0),
		},
		FeeTier: feeTier,
	}

	err = strategy.PrintResults(swaps)
	if err != nil {
		log.Fatal(err)
	}
}
