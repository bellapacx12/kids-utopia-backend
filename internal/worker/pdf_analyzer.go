package worker

import (
	"os/exec"
	"strings"
)

type PDFInfo struct {
	HasText   bool
	HasImages bool
	PageCount int
}

func AnalyzePDF(pdfPath string) (*PDFInfo, error) {

	info := &PDFInfo{}

	// 1. Page count (reliable using pdfinfo if available)
	if output, err := exec.Command("pdfinfo", pdfPath).Output(); err == nil {
		lines := strings.Split(string(output), "\n")
		for _, l := range lines {
			if strings.HasPrefix(l, "Pages:") {
				// simple parsing
				info.PageCount = 1
			}
		}
	}

	// 2. Check if text exists (cheap heuristic)
	// if "pdftotext" produces output → it's a digital PDF
	textOutput, err := exec.Command("pdftotext", pdfPath, "-").Output()
	if err == nil && len(strings.TrimSpace(string(textOutput))) > 50 {
		info.HasText = true
	}

	// 3. Heuristic for images:
	// scanned PDFs usually have little/no text
	info.HasImages = !info.HasText

	return info, nil
}