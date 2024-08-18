package db

import (
  "fmt"
	"github.com/go-pg/pg"
	"github.com/go-pg/pg/orm"
)

type Position struct {
  Id int64
  Quantity float64
  AveragePrice float64
  CurrentPrice float64
  Ppl float32
}

func createSchema(db *pg.DB) error {
  models := []interface{}{
    (*Position)(nil),
  }

  for _, model := range models {
    err := db.Model(model).CreateTable(&orm.CreateTableOptions{
      Temp: true,
    })
    if err != nil {
      return err
    }
  }
  return nil
}

func Setup() *pg.DB {
  fmt.Println("Setup")

  db := pg.Connect(&pg.Options{
    Database: "bani_development",
    User: "postgres",
    Password: "password",
  })

  err := createSchema(db)
  if err != nil {
    panic(err)
  }

  return db
}
