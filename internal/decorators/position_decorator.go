package decorators

import (
	"log"
	"github.com/CatalinVoineag/bani/internal/database"
	"github.com/Rhymond/go-money"
	"github.com/google/uuid"
)

type DecoratedPosition struct {
  Id uuid.UUID
  Ticker string
  Quantity float64
  AveragePrice float64
  CurrentPrice string
  Ppl float64
  PreviousClosePrice string
  DailyGain string
  DailyGainPercentage float64
  DailyGainNumber int64
}

func DecoratePosition(position database.Position) DecoratedPosition {
  dailyGain := calculateDailyGain(position)

  decoratedPosition := DecoratedPosition {
    Id: position.ID,
    Ticker: position.Ticker,
    Quantity: position.Quantity,
    AveragePrice: position.AveragePrice,
    CurrentPrice: currentPrice(position),
    Ppl: position.Ppl,
    PreviousClosePrice: previousPrice(position),
    DailyGain: dailyGain.Value,
    DailyGainPercentage: dailyGain.Percentage,
    DailyGainNumber: dailyGain.Number,
  }

  return decoratedPosition
}

type DailyGain struct {
  Value string
  Percentage float64
  Number int64
}

func currentPrice (position database.Position) string {
  formatter := money.NewFormatter(3, ".", "", "£", "£1")

  result := money.New(position.CurrentPrice, money.GBP).Display()

  if position.Securitytype.String != "etf" {
    result = formatter.Format(position.CurrentPrice)
  }

  return result
}

func previousPrice (position database.Position) string {
  formatter := money.NewFormatter(3, ".", "", "£", "£1")

  result := money.New(position.PreviousClosePrice.Int64, money.GBP).Display()

  if position.Securitytype.String != "etf" {
    result = formatter.Format(position.PreviousClosePrice.Int64)
  }

  return result
}

func calculateDailyGain(position database.Position) DailyGain {
  formatter := money.NewFormatter(3, ".", "", "£", "£1")

  currentPrice := money.New(position.CurrentPrice, money.GBP)

  previousPrice := money.New(
    position.PreviousClosePrice.Int64,
    money.GBP,
  )

  valuePerShare, err := currentPrice.Subtract(previousPrice)

  if err != nil {
    log.Fatal("valuePerShare conversion failed ", err)
  }

  number := int64(float64(valuePerShare.Amount()) * position.Quantity)
  totalValuePence := money.New(
     number,
    money.GBP,
  )

  totalValue := totalValuePence.Display()
  if position.Securitytype.String != "etf" {
    totalValue = formatter.Format(totalValuePence.Amount())
  }

  percentage := (float64(position.CurrentPrice - position.PreviousClosePrice.Int64) / float64(position.CurrentPrice)) * 100

  return DailyGain {
    Value: totalValue,
    Percentage: percentage,
    Number: totalValuePence.Amount(),
  }
}
