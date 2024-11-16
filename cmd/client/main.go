package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"time"
)

const (
	ApiEndoint = "http://localhost:8080/cotacao"
	MaxTimeout = 300 * time.Millisecond
	filePath   = "./cotacao.txt"
)

// CurrentCotation
type CurrentCotation struct {
	Bid string `json:"bid"` //"bid": "5.7809",
}

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), MaxTimeout)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, ApiEndoint, nil)
	if err != nil {
		log.Fatalf("error: failed to create an http request with context. %v", err)
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Default().Printf("error while executing http request.. %v", err)
		return
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatalf("error while trying to read response body. %v", err)
	}
	var data CurrentCotation
	if err = json.Unmarshal(body, &data); err != nil {
		log.Fatalf("error while trying to unmarshal json data. %v", err)
	}

	content := fmt.Sprintf(`Dolar: %s`, data.Bid)
	if err = Write2File[string](filePath, content); err != nil {
		log.Fatalf("error while trying to write the result into a file. %v", err)
	}
}

func Write2File[T string](file string, content T) error {
	if err := os.WriteFile(file, []byte(content), os.ModeAppend); err != nil {
		return err
	}

	return nil
}
