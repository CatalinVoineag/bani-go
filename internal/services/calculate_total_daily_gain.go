package total_daily_gain

import (
	//"log"

	decorators "github.com/CatalinVoineag/bani/internal/decorators"
)

type TotalGain struct {
  Value float64
  Percentage float64
  PortfolioValue float64
}

func Call(positions []decorators.DecoratedPosition) TotalGain {
  //var totalGain TotalGain
  //var positionsValueForClosedPrice float64
  //var portfolioValue float64

  //for _, position := range positions {
  //  totalGain.Value += position.DailyGain
  //  positionsValueForClosedPrice += position.Quantity * position.PreviousClosePrice.Float64 
  //  portfolioValue += position.Quantity * position.CurrentPrice
  //}

  //totalGain.Percentage = (positionsValueForClosedPrice / totalGain.Value) / 1000
  //totalGain.PortfolioValue = portfolioValue
  //log.Println("PORTFOLIO VALUE ", portfolioValue)

  return TotalGain{
    Value: 1,
    Percentage: 1,
    PortfolioValue: 1,
  }

  //return totalGain
}
