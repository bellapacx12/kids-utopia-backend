package worker

import (
	"github.com/ledongthuc/pdf"
)

func extractPDFPages(filePath string) ([]string, error) {

	f, reader, err := pdf.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	var pages []string

	totalPage := reader.NumPage()

	for i := 1; i <= totalPage; i++ {

		page := reader.Page(i)

		if page.V.IsNull() {
			continue
		}

		text, err := page.GetPlainText(nil)
		if err != nil {
			continue
		}

		pages = append(pages, text)
	}

	return pages, nil
}