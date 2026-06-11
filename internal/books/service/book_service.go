package service

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"mime/multipart"
	"time"

	"github.com/bellapacx/kids-utopia/internal/access/service"
	accessservice "github.com/bellapacx/kids-utopia/internal/access/service"
	"github.com/bellapacx/kids-utopia/internal/books/dto"
	"github.com/bellapacx/kids-utopia/internal/books/events"
	"github.com/bellapacx/kids-utopia/internal/books/model"
	"github.com/bellapacx/kids-utopia/internal/books/repository"
	"github.com/bellapacx/kids-utopia/pkg/sqs"
	"github.com/bellapacx/kids-utopia/pkg/storage"
	"github.com/google/uuid"
)

type BookService struct {
	repo repository.BookRepository
	storage  storage.Storage
queue   *sqs.Client
accessService *service.Service
}

func NewBookService(
	repo repository.BookRepository,
	storage storage.Storage,
	queue *sqs.Client,
	accessService *accessservice.Service,
) *BookService {

	return &BookService{
		repo:     repo,
		storage:  storage,
		queue: queue,
		accessService: accessService,
	}
}

func (s *BookService) CreateBook(ctx context.Context, req dto.CreateBookRequest) (*model.Book, error) {
	book := &model.Book{
		ID:          uuid.NewString(),
		Title:       req.Title,
		Description: req.Description,
		Author:      req.Author,
		Status:      "draft",
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	err := s.repo.Create(ctx, book)
	if err != nil {
		return nil, err
	}

	return book, nil
}
func (s *BookService) GetBookByID(ctx context.Context, id string) (*model.Book, error) {
	return s.repo.FindByID(ctx, id)
}
func (s *BookService) UploadBook(fileName string, fileURL string) (*model.Book, error) {
    log.Println("🔥 CreateUploadedBook called:", fileName)
	book := &model.Book{
		ID:          uuid.NewString(),
		Title:       fileName,
		Description: "",
		Author:      "",
		CoverURL:    fileURL,
		Status:      "processing",
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	err := s.repo.Create(context.Background(), book)
	if err != nil {
		return nil, err
	}

	// TODO: Kafka event trigger later
	

	return book, nil
}
func (s *BookService) UploadToStorage(ctx context.Context, file multipart.File, fileName string) (string, error) {
	return s.storage.UploadFile(ctx, file, fileName)
}
func (s *BookService) CreateUploadedBook(
	ctx context.Context,
	title string,
	author string,
	url string,
) (*model.Book, error) {
    log.Println("🔥 CreateUploadedBook called:", title)
	book := &model.Book{
	ID:        uuid.NewString(),
	Title:     title,
	Author:    author,
	CoverURL:  url,
	Status:    "processing",
	CreatedAt: time.Now(),
	UpdatedAt: time.Now(),
}
	

	// 1. Save to DB
	err := s.repo.Create(ctx, book)
	if err != nil {
		return nil, err
	}

	// 2. Build Kafka event
	event := events.BookUploadedEvent{
	BookID:   book.ID,
	ObjectKey: book.CoverURL,
	Status:    book.Status,
}

	// 3. Marshal safely
	data, err := json.Marshal(event)
	if err != nil {
		return nil, err
	}

	// 4. Publish to Kafka (IMPORTANT: don't ignore error)
	err = s.queue.Send(string(data))
if err != nil {
	return nil, err
}

	return book, nil
}
func (s *BookService) ListBooks(
	ctx context.Context,
) ([]model.Book, error) {

	return s.repo.ListBooks(ctx)
}
func (s *BookService) GetBookByIDs(
	ctx context.Context,
	bookID string,
	userID string,
	role string,
) (*model.Book, []model.BookPage, error) {

	// =========================
	// FETCH BOOK
	// =========================
	book, err := s.repo.GetBookByID(ctx, bookID)
	if err != nil {
		return nil, nil, err
	}

	// =========================
	// ACCESS CHECK
	// =========================
	allowed, preview, err := s.accessService.CanAccessBook(ctx, userID, book)
	if err != nil {
		return nil, nil, err
	}

	// =========================
	// FULL ACCESS
	// =========================
	if allowed {
		pages, err := s.repo.GetBookPages(ctx, bookID)
		if err != nil {
			return nil, nil, err
		}
		return book, pages, nil
	}

	// =========================
	// PREVIEW ACCESS
	// =========================
	if preview {
		pages, err := s.repo.GetBookPreview(ctx, bookID)
		if err != nil {
			return nil, nil, err
		}
		return book, pages, nil
	}

	// =========================
	// NO ACCESS (SAFETY FALLBACK)
	// =========================
	return nil, nil, fmt.Errorf("access denied")
}
func (s *BookService) GetBook(
	ctx context.Context,
	bookID string,
) (*model.Book, []model.BookPage, error) {

	// =========================
	// FETCH BOOK
	// =========================
	book, err := s.repo.GetBookByID(ctx, bookID)
	if err != nil {
		return nil, nil, err
	}

	// =========================
	// FETCH FULL PAGES
	// =========================
	pages, err := s.repo.GetBookPages(ctx, bookID)
	if err != nil {
		return nil, nil, err
	}

	return book, pages, nil
}

func (s *BookService) GetBookPreview(
	ctx context.Context,
	bookID string,
) (*model.Book, []model.BookPage, error) {

	// =========================
	// FETCH BOOK
	// =========================
	book, err := s.repo.GetBookByID(ctx, bookID)
	if err != nil {
		return nil, nil, err
	}

	// =========================
	// FETCH PREVIEW PAGES
	// =========================
	pages, err := s.repo.GetBookPreview(ctx, bookID)
	if err != nil {
		return nil, nil, err
	}

	return book, pages, nil
}

func (s *BookService) GetBookMeta(
	ctx context.Context,
	bookID string,
) (*model.Book, error) {

	return s.repo.GetBookByID(ctx, bookID)
}
func (s *BookService) CreateUploadedBooks(
	ctx context.Context,
	req dto.CreateUploadedBookRequest,
	fileURL string,
) (*model.Book, error) {

	log.Println("🔥 CreateUploadedBook called:", req.Title)

	// defaults
	if req.Language == "" {
		req.Language = "en"
	}

	if req.AccessType == "" {
		req.AccessType = "free"
	}

	book := &model.Book{
		ID:          uuid.NewString(),
		Title:       req.Title,
		Description: req.Description,
		Author:      req.Author,

		CoverURL: fileURL,

		AccessType: req.AccessType,

		AgeMin:   req.AgeMin,
		AgeMax:   req.AgeMax,
		Language: req.Language,
		Category: req.Category,

		Status:          "processing",
		PopularityScore: 0,

		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	// Save book
	if err := s.repo.Create(ctx, book); err != nil {
		return nil, err
	}

	// Publish processing event
	event := events.BookUploadedEvent{
		Type:      "book.uploaded",
		BookID:    book.ID,
		ObjectKey: book.CoverURL,
		Status:    book.Status,
	}

	data, err := json.Marshal(event)
	if err != nil {
		return nil, err
	}

	if err := s.queue.Send(string(data)); err != nil {
		return nil, err
	}

	return book, nil
}
func (s *BookService) CreateBookWithFirstVariant(
	ctx context.Context,
	req dto.CreateFirstVariantRequest,
	fileURL string,
) (*dto.CreateBookResponse, error) {

	// =========================
	// 1. CREATE BOOK (IF NEEDED)
	// =========================
	bookID := req.BookID

	var book *model.Book

	if bookID == "" {
		book = &model.Book{
			ID:          uuid.NewString(),
			Title:       req.Title,
			Description: req.Description,
			Author:      req.Author,

			AccessType: req.AccessType,
			AgeMin:     req.AgeMin,
			AgeMax:     req.AgeMax,
			Category:   req.Category,

			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}

		if err := s.repo.Create(ctx, book); err != nil {
			return nil, err
		}
	} else {
		b, err := s.repo.FindByID(ctx, bookID)
		if err != nil {
			return nil, err
		}
		book = b
	}

	// =========================
	// 2. CREATE FIRST VARIANT
	// =========================
	variant := &model.BookVariant{
		ID:        uuid.NewString(),
		BookID:    book.ID,
		Language:  req.Language,
		Title:     req.Title,
		FileURL:   fileURL,

		Status:   "processing",
		Progress: 0,

		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	if err := s.repo.CreateVariant(ctx, variant); err != nil {
		return nil, err
	}

	// =========================
	// 3. TRIGGER WORKER
	// =========================
	event := events.BookVariantUploaded{
		Type: "book.uploaded",
		BookID:    book.ID,
		VariantID: variant.ID,
		ObjectKey:   fileURL,
		Language:  req.Language,
	}

	data, _ := json.Marshal(event)
	s.queue.Send(string(data))

	// =========================
	// RETURN BOTH
	// =========================
	return &dto.CreateBookResponse{
		Book:    book,
		Variant: variant,
	}, nil
}
func (s *BookService) CreateBookVariant(
	ctx context.Context,
	req dto.CreateVariantRequest,
) (*model.BookVariant, error) {

	// =========================
	// 1. VALIDATE BOOK EXISTS
	// =========================
	book, err := s.repo.FindByID(ctx, req.BookID)
	if err != nil {
		return nil, fmt.Errorf("book not found: %w", err)
	}

	// =========================
	// 2. CREATE VARIANT
	// =========================
	variant := &model.BookVariant{
		ID:       uuid.NewString(),
		BookID:   book.ID,

		Language: req.Language,
		Title:    req.Title,

		FileURL: req.FileURL,

		Status:   "processing",
		Progress: 0,

		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	if err := s.repo.CreateVariant(ctx, variant); err != nil {
		return nil, err
	}

	// =========================
	// 3. TRIGGER WORKER
	// =========================
	event := events.BookVariantUploaded{
		Type:      "book.uploaded",
		BookID:    book.ID,
		VariantID: variant.ID,
		ObjectKey: variant.FileURL,
		Language:  variant.Language,
	}

	data, err := json.Marshal(event)
	if err != nil {
		return nil, err
	}

	if err := s.queue.Send(string(data)); err != nil {
		return nil, err
	}

	return variant, nil
}
func (s *BookService) ListBooksWithVariants(
	ctx context.Context,
) ([]dto.BookWithVariants, error) {

	books, err := s.repo.ListBooks(ctx)
	if err != nil {
		return nil, err
	}

	result := make([]dto.BookWithVariants, 0, len(books))

	for _, book := range books {

		variants, err := s.repo.ListVariantsByBookID(ctx, book.ID)
		if err != nil {
			return nil, err
		}

		bookDTO := dto.Book{
			ID:       book.ID,
			CoverURL: book.CoverURL,
			Title: buildVariantTitles(variants),
		}

		result = append(result, dto.BookWithVariants{
			Book:     bookDTO,
			Variants: variants,
		})
	}

	return result, nil
}

	
func (s *BookService) GetVariantsWithPages(
	ctx context.Context,
	bookID string,
) ([]dto.ReaderVariant, error) {

	variants, err := s.repo.ListVariantsByBookID(ctx, bookID)
	if err != nil {
		return nil, err
	}

	result := make([]dto.ReaderVariant, 0, len(variants))

	for _, v := range variants {

		pages, err := s.repo.GetPagesByVariantIDD(
			ctx,
			v.ID,
		)
		if err != nil {
			return nil, err
		}

		result = append(result, dto.ReaderVariant{
			ID:       v.ID,
			Language: v.Language,
			Pages:    pages, // []model.BookPage
		})
	}

	return result, nil
}
func (s *BookService) GetVariantsWithPreview(
	ctx context.Context,
	bookID string,
	limit int,
) ([]dto.ReaderVariant, error) {

	variants, err := s.repo.ListVariantsByBookID(ctx, bookID)
	if err != nil {
		return nil, err
	}

	result := make([]dto.ReaderVariant, 0, len(variants))

	for _, v := range variants {

		pages, err := s.repo.GetPagesByVariantIDD(ctx, v.ID)
		if err != nil {
			return nil, err
		}

		// 🔐 apply preview per variant
		if limit > 0 && len(pages) > limit {
			pages = pages[:limit]
		}

		result = append(result, dto.ReaderVariant{
			ID:       v.ID,
			Language: v.Language,
			Pages:    pages,
		})
	}

	return result, nil
}
var langMap = map[string]string{
	"am": "Amharic",
	"en": "English",
	"ti": "Tigrigna",
	"om": "Oromic",
}
func buildVariantTitles(variants []model.BookVariant) []string {
	titles := make([]string, 0, len(variants))

	for _, v := range variants {

		if v.Title == "" {
			continue
		}

		label := v.Title

		if langName, ok := langMap[v.Language]; ok {
			label = langName + ": " + v.Title
		} else if v.Language != "" {
			label = v.Language + ": " + v.Title
		}

		titles = append(titles, label)
	}

	return titles
}