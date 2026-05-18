package worker

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
)

func extractPDFImages(
	pdfURL string,
	outputDir string,
) ([]string, error) {

	// ensure directory exists
	err := os.MkdirAll(outputDir, os.ModePerm)
	if err != nil {
		return nil, err
	}

	// send URL instead of local file path
	reqBody := map[string]string{
		"pdf_url":   pdfURL,
		"output_dir": outputDir,
	}

	data, _ := json.Marshal(reqBody)

	resp, err := http.Post(
		"http://pdf-worker:8001/process",
		"application/json",
		bytes.NewBuffer(data),
	)

	if err != nil {
		return nil, fmt.Errorf(
			"python worker request failed: %v",
			err,
		)
	}

	defer resp.Body.Close()

	var result struct {
		Success bool     `json:"success"`
		Pages   []string `json:"pages"`
		Count   int      `json:"count"`
		Error   string   `json:"error"`
	}

	err = json.NewDecoder(resp.Body).Decode(&result)
	if err != nil {
		return nil, fmt.Errorf(
			"invalid python response: %v",
			err,
		)
	}

	if !result.Success {
		return nil, fmt.Errorf(
			"python worker failed: %s",
			result.Error,
		)
	}

	return result.Pages, nil
}