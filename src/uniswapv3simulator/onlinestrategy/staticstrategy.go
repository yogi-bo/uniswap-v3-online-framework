package onlinestrategy

import (
	"math"
	"time"
	"uniswapv3simulator/uniswapmath"
)

type StaticStrategy struct {
	ConcentrationController uint32
	Duration                time.Duration
	LastDeposit             time.Time
}

func (l *StaticStrategy) NextRange(onlineStrategy *OnlineStrategy) (Range, error) {
	if l.LastDeposit.Add(l.Duration).After(onlineStrategy.LastSwap.BlockTimestamp) {
		return onlineStrategy.ActiveRange, nil
	}
	l.LastDeposit = l.LastDeposit.Add((onlineStrategy.LastSwap.BlockTimestamp.Sub(l.LastDeposit) / l.Duration) * l.Duration)

	controllerFactor := math.Pow(uniswapmath.IntervalSize, float64(onlineStrategy.ActiveRange.TickSpacing)*float64(l.ConcentrationController))
	return Range{
		TickSpacing: onlineStrategy.ActiveRange.TickSpacing,
		High:        onlineStrategy.LastSwap.Price * controllerFactor,
		Low:         onlineStrategy.LastSwap.Price / controllerFactor,
	}, nil
}
