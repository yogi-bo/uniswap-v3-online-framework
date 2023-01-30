package uniswapmath

import "math"

const IntervalSize = 1.0001

func getSqrtPrices(currentTick, lowerTick, upperTick int32) (sqrtCurrentPrice, sqrtLowerPrice, sqrtUpperPrice float64) {
	sqrtCurrentPrice = math.Pow(IntervalSize, float64(currentTick)/2)
	sqrtLowerPrice = math.Pow(IntervalSize, float64(lowerTick)/2)
	sqrtUpperPrice = math.Pow(IntervalSize, float64(upperTick)/2)

	return
}

func liquidityToAmount0(sqrtLowerPrice, sqrtUpperPrice float64, liquidity float64) float64 {
	return liquidity * (sqrtUpperPrice - sqrtLowerPrice) / (sqrtLowerPrice * sqrtUpperPrice)
}

func liquidityToAmount1(sqrtLowerPrice, sqrtUpperPrice float64, liquidity float64) float64 {
	return liquidity * (sqrtUpperPrice - sqrtLowerPrice)
}

func LiquidityToAmounts(currentTick, lowerTick, upperTick int32, liquidity float64) (amount0, amount1 float64) {
	sqrtCurrentPrice, sqrtLowerPrice, sqrtUpperPrice := getSqrtPrices(currentTick, lowerTick, upperTick)

	if sqrtCurrentPrice < sqrtLowerPrice {
		amount0 = liquidityToAmount0(sqrtLowerPrice, sqrtUpperPrice, liquidity)
	} else if sqrtCurrentPrice > sqrtUpperPrice {
		amount1 = liquidityToAmount1(sqrtLowerPrice, sqrtUpperPrice, liquidity)
	} else {
		amount0 = liquidityToAmount0(sqrtCurrentPrice, sqrtUpperPrice, liquidity)
		amount1 = liquidityToAmount1(sqrtLowerPrice, sqrtCurrentPrice, liquidity)
	}
	return
}

func amount0ToLiquidity(sqrtLowerPrice, sqrtUpperPrice float64, amount0 float64) float64 {
	return amount0 * (sqrtLowerPrice * sqrtUpperPrice) / (sqrtUpperPrice - sqrtLowerPrice)
}

func amount1ToLiquidity(sqrtLowerPrice, sqrtUpperPrice float64, amount1 float64) float64 {
	return amount1 / (sqrtUpperPrice - sqrtLowerPrice)
}

func AmountsToLiquidity(currentTick, lowerTick, upperTick int32, amount0, amount1 float64) float64 {
	sqrtCurrentPrice, sqrtLowerPrice, sqrtUpperPrice := getSqrtPrices(currentTick, lowerTick, upperTick)

	if sqrtCurrentPrice < sqrtLowerPrice {
		return amount0ToLiquidity(sqrtLowerPrice, sqrtUpperPrice, amount0)
	} else if sqrtCurrentPrice > sqrtUpperPrice {
		return amount1ToLiquidity(sqrtLowerPrice, sqrtUpperPrice, amount1)
	} else {
		return math.Min(
			amount0ToLiquidity(sqrtCurrentPrice, sqrtUpperPrice, amount0),
			amount1ToLiquidity(sqrtLowerPrice, sqrtCurrentPrice, amount1),
		)
	}
}

func PriceToTick(price float64) int32 {
	return int32(math.Log(price) / math.Log(IntervalSize))
}
