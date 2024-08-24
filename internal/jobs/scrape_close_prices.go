package jobs

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"

	"strings"

	"github.com/CatalinVoineag/bani/internal/database"
	"github.com/google/uuid"
)

type Meta struct {
  Symbol string `json:"symbol"`
  PreviousClose float64 `json:"previousClose"`
}

type Result struct {
  Meta Meta `json:"meta"`
}

type Chart struct {
  Result []Result `json:"result"`
}

type Response struct {
  Chart Chart `json:"chart"`
}

type PreviousClose struct {
  Symbol string
  Close float64 
}

func ScrapePreviousClosePrice(db *database.Queries) {
  dbResults, err := db.GetTodayPositionsTickers(context.Background())
  if err != nil {
    log.Println("No positions to scrape open prices")
  }

  for _, tickerIdPair := range dbResults {
    pair, ok := tickerIdPair.([]uint8)

    if !ok {
      log.Println("Failed to get the scrape TickerId pair ", err)
    } else {
      pair := string(pair)
      mapPair := strings.Split(strings.Trim(string(pair), "()"), ",")
      positionId := mapPair[1]
      ticker := mapPair[0]

      tickerWithoutEQ, _, found := strings.Cut(ticker, "_")

      if found == false {
        log.Printf("Underscore separator not found for %s", ticker)
      } else {
        exchange := tickerWithoutEQ[len(tickerWithoutEQ) -1:]
        ticker := tickerWithoutEQ[:len(tickerWithoutEQ) -1]
        yahooQuoteParam := ticker + "." + exchange 

        closePrice := getPreviousClosePrice(yahooQuoteParam).Close

        log.Printf("Yahoo close price for %s %f", yahooQuoteParam, closePrice) 

        uid, _ := uuid.Parse(positionId)

        db.UpdatePreviousClosedPrice(context.Background(), database.UpdatePreviousClosedPriceParams {
          PreviousClosePrice: float64ToNullFloat64(closePrice),
          ID:                 uid,
        })
      }
    }
  }
}

func getPreviousClosePrice(ticker string) PreviousClose {
  reqUrl := "https://query1.finance.yahoo.com/v8/finance/chart/" + ticker
  req, err := http.NewRequest("GET", reqUrl, nil)

  if err != nil {
    panic(err)
  }

  res, err := http.DefaultClient.Do(req)

  if err != nil {
    panic(err)
  }

  defer res.Body.Close()
  body, err := io.ReadAll(res.Body)

  if err != nil {
    panic(err)
  }

  var response Response
  json.Unmarshal([]byte(body), &response)

  previousClose := PreviousClose {
    Symbol: response.Chart.Result[0].Meta.Symbol,
    Close: response.Chart.Result[0].Meta.PreviousClose,
  }

  return previousClose
}

func ScrapeOpenPrices(db *database.Queries) {
	fmt.Println("hello")

}

func float64ToNullFloat64(value float64) sql.NullFloat64 {
	return sql.NullFloat64{
		Float64: value,
		Valid:   true,
	}
}
