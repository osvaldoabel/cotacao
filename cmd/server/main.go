package main

import (
	"context"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"time"

	"github.com/osvaldoabel/cotacao/pkg/utils"
)

const (
	RouteGetIndex           = "/"
	ExchageProviderUrl      = "https://economia.awesomeapi.com.br/json/last/USD-BRL"
	ExchangeProviderTimeout = 15 * time.Second
)

type ExchangeProvider interface {
	Execute(ctx context.Context) (utils.GetExchangeResponse, error)
}

type HttpHandler interface {
	Index(w http.ResponseWriter, r *http.Request)
}

type conversionHandler struct {
	exchangeProvider ExchangeProvider
}

// NewConversionHandler
func NewConversionHandler() HttpHandler {
	return &conversionHandler{
		exchangeProvider: NewExchangeProvider(),
	}
}

// Index
func (h *conversionHandler) Index(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), ExchangeProviderTimeout)
	defer cancel()

	data, err := h.exchangeProvider.Execute(ctx)
	if err != nil {
		log.Fatalf("error: %v", err)
	}

	utils.JsonResponse(w, data.USDBRL, http.StatusOK)
}

// exchangeProvider
type exchangeProvider struct{}

// NewExchangeProvider
func NewExchangeProvider() ExchangeProvider {
	return &exchangeProvider{}
}

// Execute
func (h *exchangeProvider) Execute(ctx context.Context) (utils.GetExchangeResponse, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, ExchageProviderUrl, nil)
	if err != nil {
		return utils.GetExchangeResponse{}, err
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return utils.GetExchangeResponse{}, err
	}
	defer resp.Body.Close()

	byteBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return utils.GetExchangeResponse{}, err
	}

	var result utils.GetExchangeResponse
	err = json.Unmarshal(byteBody, &result)
	if err != nil {
		return utils.GetExchangeResponse{}, err
	}

	// fmt.Println(result)

	return result, err
	// io.Copy(os.Stdout, resp.Body)
}

func main() {
	handler := NewConversionHandler()
	http.HandleFunc(RouteGetIndex, handler.Index)
	http.ListenAndServe(":8808", nil)

}
