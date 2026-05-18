package worker

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
)

func renderPDFPagesToPNG(pdfPath string, outputDir string) ([]string, error) {

	// ensure output directory exists
	if err := os.MkdirAll(outputDir, os.ModePerm); err != nil {
		return nil, err
	}

	outputPrefix := filepath.Join(outputDir, "page")

	// run poppler tool
	cmd := exec.Command(
		"pdftoppm",
		"-png",
		pdfPath,
		outputPrefix,
	)

	output, err := cmd.CombinedOutput()
	if err != nil {
		return nil, fmt.Errorf("pdftoppm failed: %v - %s", err, string(output))
	}

	// collect generated images
	files, err := filepath.Glob(outputPrefix + "-*.png")
	if err != nil {
		return nil, err
	}

	return files, nil
}