package service

import (
	"context"
	"encoding/json"
	"log"
	"mime/multipart"
	"time"

	"github.com/google/uuid"

	"github.com/bellapacx/kids-utopia/internal/books/dto"
	"github.com/bellapacx/kids-utopia/internal/books/events"
	"github.com/bellapacx/kids-utopia/internal/books/model"
	"github.com/bellapacx/kids-utopia/internal/books/repository"
	"github.com/bellapacx/kids-utopia/pkg/sqs"
	"github.com/bellapacx/kids-utopia/pkg/storage"
)

type BookService struct {
	repo repository.BookRepository
	storage  storage.Storage
queue   *sqs.Client
}

func NewBookService(
	repo repository.BookRepository,
	storage storage.Storage,
	queue *sqs.Client,
) *BookService {

	return &BookService{
		repo:     repo,
		storage:  storage,
		queue: queue,
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