package jobs

import (
	"context"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"os"
	"sync"
	"time"

	"github.com/CatalinVoineag/bani/internal/database"
	"github.com/google/uuid"
  _"github.com/lib/pq"
)

type apiConfig struct {
  DB *database.Queries
}

type Position struct {
  Ticker string `json:"ticker"` 
  Quantity float64 `json:"quantity"`
  AveragePrice float64 `json:"averagePrice"`
  CurrentPrice float64 `json:"currentPrice"`
  Ppl float64 `json:"ppl"`
}

type Positions = []Position

func Start(db *database.Queries, timeBetweenRequest time.Duration) {
  log.Printf("Fetching positions every %s duration", timeBetweenRequest)

  ticker := time.NewTicker(timeBetweenRequest)

  for ; ; <-ticker.C {
    wg := &sync.WaitGroup{}
    wg.Add(1)

    go currentPositionsWorker(db, wg)

    wg.Wait()
    log.Printf("Finished")
  }
}

func currentPositionsWorker(db *database.Queries, wg *sync.WaitGroup) {
  log.Printf("Start")
  defer wg.Done()
  positions := getTradingTwoOneTwoPositions()
  if len(positions) > 0 {
    for _, position := range positions {
      
      record, err := db.CreatePosition(context.Background(), database.CreatePositionParams {
        ID: uuid.New(),
        CreatedAt: time.Now(),
        UpdatedAt: time.Now(),
        Quantity: position.Quantity,
        AveragePrice: position.AveragePrice,
        CurrentPrice: position.CurrentPrice,
        Ppl: position.Ppl,
        Ticker: position.Ticker,
      })

      if err != nil {
        log.Println("Position cannot be created ", err)
      }
      
      log.Println("Position Created ticker: ", record.Ticker)
      lastPositions, err := db.GetLastPositionsTodayByTickerExcludingCurrent(
        context.Background(),
        database.GetLastPositionsTodayByTickerExcludingCurrentParams{
          Ticker: record.Ticker,
          ID: record.ID,
        },
      )

      if err != nil {
        log.Println("No last position")
      } else {
        for _, position := range lastPositions {
          _, err = db.DeletePoistion(context.Background(), position.ID)
          if err != nil {
            log.Printf("Could not delete position %s error: %e", position.ID, err)
          }
        }
        log.Println("Last positions deleted")
      }
    }
  } else {
    log.Println("No positions")
  }
}

func getTradingTwoOneTwoPositions() Positions {
  reqUrl := "https://live.trading212.com/api/v0/equity/portfolio"
  req, err := http.NewRequest("GET", reqUrl, nil)
  if err != nil {
    panic(err)
  }
  req.Header.Add("Authorization", os.Getenv("API_KEY"))
  res, err := http.DefaultClient.Do(req)
  if err != nil {
    panic(err)
  }
  defer res.Body.Close()
  body, err := io.ReadAll(res.Body)
  if err != nil {
    panic(err)
  }

  var positions Positions
  json.Unmarshal([]byte(body), &positions)

  return positions
}
