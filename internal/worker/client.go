package worker

import (
	"bytes"
	"encoding/json"
	"net/http"
)

type PythonRequest struct {
	PDFPath   string `json:"pdf_path"`
	OutputDir string `json:"output_dir"`
}

type PythonResponse struct {
	Success bool     `json:"success"`
	Pages   []string `json:"pages"`
	Count   int      `json:"count"`
	Error   string   `json:"error"`
}
func RunPythonWorker(pdfPath, outputDir string) (*PythonResponse, error) {

	reqBody := PythonRequest{
		PDFPath:   pdfPath,
		OutputDir: outputDir,
	}

	data, err := json.Marshal(reqBody)
	if err != nil {
		return nil, err
	}

	resp, err := http.Post(
		"http://pdf-worker:8001/process",
		"application/json",
		bytes.NewBuffer(data),
	)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var result PythonResponse

	err = json.NewDecoder(resp.Body).Decode(&result)
	if err != nil {
		return nil, err
	}

	return &result, nil
}