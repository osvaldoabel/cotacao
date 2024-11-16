package main

import (
	"fmt"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func getTempFolder() string {
	return os.TempDir()
}
func TestWrite2File(t *testing.T) {
	tmpFolder := getTempFolder()
	tmpFile := fmt.Sprintf("%scotacao.txt", tmpFolder)
	defer os.RemoveAll(tmpFile)
	err := Write2File[string](tmpFile, "simple content")
	assert.Nil(t, err)
}
