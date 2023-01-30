package onlinestrategy

import (
	"fmt"
	"math"
	"uniswapv3simulator/swaputils"
	"uniswapv3simulator/uniswapmath"
)

type Strategy interface {
	NextRange(onlineStrategy *OnlineStrategy) (Range, error)
}

type Range struct {
	TickSpacing uint32
	Low         float64
	High        float64
}

func (r Range) toTicks() (lowerTick int32, upperTick int32) {
	lowerTick = int32(math.Floor(float64(uniswapmath.PriceToTick(r.Low))/float64(r.TickSpacing)) * float64(r.TickSpacing))
	upperTick = int32(math.Ceil(float64(uniswapmath.PriceToTick(r.High))/float64(r.TickSpacing)) * float64(r.TickSpacing))
	return
}

func (r Range) isActive(tick int32) bool {
	lowerTick, upperTick := r.toTicks()
	return tick <= upperTick && tick >= lowerTick
}

type OnlineStrategy struct {
	InvestedLiquidity float64
	UninvestedToken0  float64
	UninvestedToken1  float64
	ActiveRange       Range
	LastSwap          swaputils.Swap
	Strategy          Strategy
	FeeTier           float64
}

func (o *OnlineStrategy) totalValueInToken1() float64 {
	lowerTick, upperTick := o.ActiveRange.toTicks()
	investedAmount0, investedAmount1 := uniswapmath.LiquidityToAmounts(uniswapmath.PriceToTick(o.LastSwap.Price), lowerTick, upperTick, o.InvestedLiquidity)
	return (o.UninvestedToken0+investedAmount0)*o.LastSwap.Price + o.UninvestedToken1 + investedAmount1
}

func (o *OnlineStrategy) totalValueInToken0() float64 {
	return o.totalValueInToken1() / o.LastSwap.Price
}

func (o *OnlineStrategy) withdraw() {
	o.UninvestedToken1 = o.totalValueInToken1() / 2
	o.UninvestedToken0 = o.UninvestedToken1 / o.LastSwap.Price
	o.InvestedLiquidity = 0
}

func (o *OnlineStrategy) reinvest(newRange Range) {
	if newRange == o.ActiveRange {
		return
	}
	o.withdraw()

	o.ActiveRange = newRange
	lowerTick, upperTick := o.ActiveRange.toTicks()

	o.InvestedLiquidity = uniswapmath.AmountsToLiquidity(uniswapmath.PriceToTick(o.LastSwap.Price), lowerTick, upperTick, o.UninvestedToken0, o.UninvestedToken1)
	investedAmount0, investedAmount1 := uniswapmath.LiquidityToAmounts(uniswapmath.PriceToTick(o.LastSwap.Price), lowerTick, upperTick, o.InvestedLiquidity)
	o.UninvestedToken0 = math.Max(o.UninvestedToken0-investedAmount0, 0)
	o.UninvestedToken1 = math.Max(o.UninvestedToken1-investedAmount1, 0)
}

func (o *OnlineStrategy) earnFees(swap swaputils.Swap) {
	if !o.ActiveRange.isActive(swap.Tick) || swap.Liquidity <= 0 {
		return
	}

	if swap.Amount0 <= 0 {
		o.UninvestedToken0 += -swap.Amount0 * o.FeeTier * o.InvestedLiquidity / (o.InvestedLiquidity + swap.Liquidity)
	} else {
		o.UninvestedToken1 += -swap.Amount1 * o.FeeTier * o.InvestedLiquidity / (o.InvestedLiquidity + swap.Liquidity)
	}
}

func (o *OnlineStrategy) PrintResults(swaps []swaputils.Swap) error {
	for _, swap := range swaps {
		o.earnFees(swap)
		fmt.Println(o.totalValueInToken0())

		o.LastSwap = swap
		newRange, err := o.Strategy.NextRange(o)
		if err != nil {
			return err
		}

		o.reinvest(newRange)
	}

	return nil
}
