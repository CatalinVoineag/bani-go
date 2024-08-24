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
      lastPosition, err := db.GetLastPositionTodayByTicker(
        context.Background(),
        position.Ticker,
      )

      if err != nil {
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
          log.Printf("Position cannot be created for %s error: %e", position.Ticker, err)
        } else {
          log.Println("Position created for: ", record.Ticker)
        }
      } else {
        record, err := db.UpdatePosition(
          context.Background(),
          database.UpdatePositionParams{
            Quantity: position.Quantity ,  
            AveragePrice: position.AveragePrice,
            CurrentPrice: position.CurrentPrice,
            Ppl: position.Ppl,
            Ticker: position.Ticker,
            ID: lastPosition.ID,
          },
        )

        if err != nil {
          log.Printf("Position cannot be updated for %s error: %e ", record.Ticker, err)
        } else {
          log.Println("Position be updated for %r", record.Ticker)
        } 
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
