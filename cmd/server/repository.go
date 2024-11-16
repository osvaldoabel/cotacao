package main

import (
	"context"
	"log"

	_ "github.com/mattn/go-sqlite3"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

const (
	dbFile = "database.db"
)

type Cotation struct {
	ID          uint   ` gorm:"primarykey"`
	Code        string `json:"code"`        // "code": "USD",
	CodeIn      string `json:"codein"`      //"codein": "BRL",
	Name        string `json:"name"`        //"name": "Dólar Americano/Real Brasileiro",
	High        string `json:"high"`        //"high": "5.8296",
	Low         string `json:"low"`         //"low": "5.7215",
	VarBid      string `json:"varBid"`      //"varBid": "-0.0249",
	PctChange   string `json:"pctChange"`   //"pctChange": "-0.43",
	Bid         string `json:"bid"`         //"bid": "5.7809",
	Ask         string `json:"ask"`         //"ask": "5.7814",
	CurrentTime string `json:"timestamp"`   //"timestamp": "1731604922",
	CreateDate  string `json:"create_date"` //"create_date": "2024-11-14 14:22:02"
}

type DatabaseRepository interface {
	Insert(ctx context.Context, c Cotation) (Cotation, error)
}

type sqliteRepository struct {
	Conn *gorm.DB
}

func GetDBConnection() (*gorm.DB, error) {
	conn, err := gorm.Open(sqlite.Open(dbFile), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	err = conn.AutoMigrate(&Cotation{})
	if err != nil {
		return nil, err
	}

	return conn, nil

}

// NewSqliteRepository
func NewSqliteRepository(conn *gorm.DB) (DatabaseRepository, error) {
	return &sqliteRepository{Conn: conn}, nil
}

func (r *sqliteRepository) Insert(ctx context.Context, c Cotation) (Cotation, error) {
	ctx, cancel := context.WithTimeout(ctx, DatabaseInsertTimeout)
	defer cancel()
	result := r.Conn.Create(&c)
	if result.Error != nil {
		log.Default().Println("error while trying to insert data into the database.")
		return c, result.Error
	}
	log.Default().Println("data inserted successfully")
	return c, nil
}
