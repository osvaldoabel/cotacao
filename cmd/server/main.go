package main

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"log"
	"net/http"
	"time"

	"github.com/osvaldoabel/cotacao/pkg/utils"

	_ "github.com/mattn/go-sqlite3"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

const (
	RouteGetCotacao         = "/cotacao"
	ExchageProviderUrl      = "https://economia.awesomeapi.com.br/json/last/USD-BRL"
	ExchangeProviderTimeout = 200 * time.Millisecond
	DatabaseInsertTimeout   = 10 * time.Millisecond

	dbFile = "database.db"
)

type ExchangeProvider interface {
	Execute(ctx context.Context) (Cotation, error)
}

type HttpHandler interface {
	Index(w http.ResponseWriter, r *http.Request)
}

type conversionHandler struct {
	repository       DatabaseRepository
	exchangeProvider ExchangeProvider
}

type GetExchangeResponse struct {
	USDBRL Cotation `json:"USDBRL"`
}

// NewConversionHandler
func NewConversionHandler(repo DatabaseRepository) HttpHandler {
	return &conversionHandler{
		repository:       repo,
		exchangeProvider: NewExchangeProvider(),
	}
}

// Index
func (h *conversionHandler) Index(w http.ResponseWriter, r *http.Request) {
	cotationResult, err := h.exchangeProvider.Execute(r.Context())
	if err != nil {

		log.Default().Printf("error while trying to call exchange provider to execute conversion.", err)
		utils.JsonResponse(w, nil, http.StatusExpectationFailed)
		return
	}

	if _, err := h.repository.Insert(r.Context(), cotationResult); err != nil {
		log.Default().Printf("error while trying to insert data into the database.", err)
		utils.JsonResponse(w, nil, http.StatusInternalServerError)
		return
	}

	utils.JsonResponse(w, cotationResult, http.StatusOK)
}

// exchangeProvider
type exchangeProvider struct{}

// NewExchangeProvider
func NewExchangeProvider() ExchangeProvider {
	return &exchangeProvider{}
}

// formatCotationResonse
func (h *exchangeProvider) formatCotationResonse(body io.Reader) (Cotation, error) {
	log.Default().Println("======================")
	byteBody, err := io.ReadAll(body)
	if err != nil {
		return Cotation{}, err
	}

	var result GetExchangeResponse
	err = json.Unmarshal(byteBody, &result)
	if err != nil {
		log.Default().Printf("error while trying to Unmarshal response body", err)
		return Cotation{}, err
	}

	return result.USDBRL, nil
}

// Execute
func (h *exchangeProvider) Execute(ctx context.Context) (Cotation, error) {
	ctx, cancel := context.WithTimeout(ctx, ExchangeProviderTimeout)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, ExchageProviderUrl, nil)
	if err != nil {
		return Cotation{}, err
	}

	apiResponse := make(chan *http.Response)
	go func() {
		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			log.Default().Printf("error while executing http request.", err)
			return
		}

		apiResponse <- resp
	}()

	select {
	case result := <-apiResponse:
		defer result.Body.Close()
		return h.formatCotationResonse(result.Body)
	case <-ctx.Done():
		log.Default().Printf("ctx done. it was canceled or timed out.", ctx.Err())
		return Cotation{}, errors.New("error: ctx done. it was canceled or timed out.")
	}
}

func main() {
	conn, err := GetDBConnection()
	if err != nil {
		log.Fatalf("error while trying to create a database connection. %v", err)
	}

	repo, err := NewSqliteRepository(conn)
	if err != nil {
		log.Fatalf("error while trying to create sqlLite connection", err)
	}

	handler := NewConversionHandler(repo)
	http.HandleFunc(RouteGetCotacao, handler.Index)
	http.ListenAndServe(":8080", nil)

}

// ////////////////////////////////| repository |\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\

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
