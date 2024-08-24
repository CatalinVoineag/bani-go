package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"html/template"
	"io"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/CatalinVoineag/bani/internal/database"
	"github.com/CatalinVoineag/bani/internal/jobs"
	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	_ "github.com/lib/pq"
)

type Templates struct {
  templates *template.Template
}

func (t *Templates) Render(
	w io.Writer,
	name string,
	data interface{},
	c echo.Context,
) error {
  return t.templates.ExecuteTemplate(w, name, data)
}

func newTemplate() *Templates {
  return &Templates {
    templates: template.Must(template.ParseGlob("views/*.html")),
  }
}

type TradingPosition struct {
  Ticker string `json:"ticker"` 
  Quantity float64 `json:"quantity"`
  AveragePrice float64 `json:"averagePrice"`
  CurrentPrice float64 `json:"currentPrice"`
  Ppl float32 `json:"ppl"`
}

func newPosition(ticker string, quantity float64, averagePrice float64, currentPrice float64, ppl float32, previousClosePrice float64) Position {
  return Position {
    Ticker: ticker,
    Quantity: quantity,
    AveragePrice: averagePrice,
    CurrentPrice: currentPrice,
    Ppl: ppl,
    PreviousClosePrice: previousClosePrice,
  }
}

type Positions = []TradingPosition

type Data struct {
  Positions []database.Position
}

type apiConfig struct {
  DB *database.Queries
}

func getPositions() Positions {
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

type Page struct {
  Data Data
}

type Position struct {
  Id int64
  Ticker string
  Quantity float64
  AveragePrice float64
  CurrentPrice float64
  Ppl float32
  PreviousClosePrice float64
}

func main() {
  godotenv.Load(".env")
  e := echo.New()
  e.Use(middleware.Logger())

  e.Renderer = newTemplate()
  
  dbURL := os.Getenv("DB_URL")
  if dbURL == "" {
    log.Fatal("DB_URL not found")
  }

  conn, err := sql.Open("postgres", dbURL)
  if err != nil {
    log.Fatal("Can't connect to DB")
  }

  db := database.New(conn)

  go jobs.Start(db, (15 * time.Minute))
  go jobs.ScrapePreviousClosePrice(db)

  e.GET("/", func(c echo.Context) error {
    positions, err := db.GetTodayPositions(context.Background())

    if err != nil {
      e.Logger.Fatal("No positions")
    }

    page := Page {
      Data: Data { 
        Positions: positions,
      },
    }

    return c.Render(200, "index", page)
  })

  e.Logger.Fatal(e.Start(":3000"))
}
