package main

import (
	"database/sql"

	_ "github.com/mattn/go-sqlite3"
)
type Money float64

type Cotation struct {
	Code        string       `json:"code"`      // "code": "USD",
	CodeIn      string       `json:"codein"`    //"codein": "BRL",
	Name        string       `json:"name"`      //"name": "Dólar Americano/Real Brasileiro",
	High        Money        `json:"high"`      //"high": "5.8296",
	Low         Money        `json:"low"`       //"low": "5.7215",
	VarBid      string       `json:"varBid"`    //"varBid": "-0.0249",
	PctChange   string       `json:"pctChange"` //"pctChange": "-0.43",
	Bid         string       `json:"bid"`       //"bid": "5.7809",
	Ask     	string    	 `json:"ask"` 	//"ask": "5.7814",
	CurrentTime time.Time    `json:"timestamp"`   //"timestamp": "1731604922",
	CreateDate  time.Time    `json:"create_date"` //"create_date": "2024-11-14 14:22:02"
}

type DatabaseRepository interface {
	Insert(ctx context.Context, c Cotation) error
}

type sqliteRepository struct {
	Conn * sql.DB
}

func NewSqliteRepository() DatabaseRepository {
	return sqliteRepository{
		Conn 
	}
}