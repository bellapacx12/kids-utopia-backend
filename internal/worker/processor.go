package worker

import (
	"context"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/bellapacx/kids-utopia/internal/books/dto"
	"github.com/bellapacx/kids-utopia/internal/books/events"
	"github.com/bellapacx/kids-utopia/internal/books/repository"
	"github.com/bellapacx/kids-utopia/pkg/storage"
)


func ProcessBook(
	event events.BookUploadedEvent,
	st storage.Storage,
	repo repository.BookPagesRepository,
) error{

	ctx := context.Background()

	filePath := filepath.Join("/tmp", event.BookID+".pdf")

	// =========================
	// STEP 1: DOWNLOAD PDF (FIXED)
	// =========================
	downloadURL, err := st.GetPresignedURL(ctx, event.ObjectKey)
	if err != nil {
		log.Printf("❌ presigned url error (book=%s): %v", event.BookID, err)
		return err
	}

	if err := downloadPDF(downloadURL, filePath); err != nil {
		log.Printf("❌ download error (book=%s): %v", event.BookID, err)
		return err
	}

	log.Println("📥 Downloaded PDF for book:", event.BookID)

	// =========================
	// STEP 2: ANALYZE PDF
	// =========================
	info, err := AnalyzePDF(filePath)
	if err != nil {
		log.Printf("❌ PDF analysis failed (book=%s): %v", event.BookID, err)
		return err
	}

	log.Printf(
		"📊 PDF TYPE → pages:%d text:%v images:%v",
		info.PageCount,
		info.HasText,
		info.HasImages,
	)

	// =========================
	// STEP 3: EXTRACT TEXT
	// =========================
	textPages, err := extractPDFPages(filePath)
	if err != nil {
		log.Printf("❌ text extraction error (book=%s): %v", event.BookID, err)
		return err
	}

	log.Printf("📄 Extracted %d text pages", len(textPages))

	// =========================
	// STEP 4: RENDER IMAGES (PNG)
	// =========================
	outputDir := filepath.Join("/tmp", event.BookID)

	imagePaths, err := renderPDFPagesToPNG(filePath, outputDir)
	if err != nil {
		log.Printf("❌ render error (book=%s): %v", event.BookID, err)
		return err
	}

	if len(imagePaths) == 0 {
		log.Printf("❌ no pages rendered (book=%s)", event.BookID)
		return err
	}

	log.Printf("🖼️ Rendered %d pages", len(imagePaths))

	// =========================
	// STEP 5: UPLOAD IMAGES (STORE OBJECT KEYS ONLY)
	// =========================
	var imageKeys []string


for i, imagePath := range imagePaths {

	file, err := os.Open(imagePath)
	if err != nil {
		log.Printf("❌ open image failed: %v", err)
		continue
	}

	objectName := fmt.Sprintf("%s/page-%d.png", event.BookID, i+1)

	_, err = st.UploadFile(ctx, file, objectName)
	file.Close()

	if err != nil {
		log.Printf("❌ upload failed: %v", err)
		continue
	}

	imageKeys = append(imageKeys, objectName)
}

	// =========================
	// STEP 6: SAVE DRAFT
	// =========================
	pages := make([]dto.EditorPageDTO, len(textPages))

for i := range textPages {
	pages[i] = dto.EditorPageDTO{
		PageNumber: i + 1,
		Content:    textPages[i],
		ImageKey:   imageKeys[i],
	}
}

err = repo.SavePages(ctx, event.BookID, pages)
if err != nil {
	log.Printf("❌ save error (book=%s): %v", event.BookID, err)
	return err
}

	log.Printf(
		"✅ Book processed successfully: %s (pages=%d)",
		event.BookID,
		len(textPages),
	)
	return nil
}