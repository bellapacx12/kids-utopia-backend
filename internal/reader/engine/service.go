package engine

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"time"

	accesssvc "github.com/bellapacx/kids-utopia/internal/access/service"
	bookdto "github.com/bellapacx/kids-utopia/internal/books/dto"
	bookmodel "github.com/bellapacx/kids-utopia/internal/books/model"
	booksvc "github.com/bellapacx/kids-utopia/internal/books/service"
	"github.com/bellapacx/kids-utopia/internal/events"
	"github.com/bellapacx/kids-utopia/internal/progress/repository"
	progresssvc "github.com/bellapacx/kids-utopia/internal/progress/service"
	"github.com/bellapacx/kids-utopia/internal/reader/constants"
	readermodel "github.com/bellapacx/kids-utopia/internal/reader/model"
	sessionsvc "github.com/bellapacx/kids-utopia/internal/reader_session/service"
	streakservice "github.com/bellapacx/kids-utopia/internal/streak/service"
	"github.com/bellapacx/kids-utopia/pkg/kafka"
	"github.com/bellapacx/kids-utopia/pkg/sqs"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type Engine struct {
	accessService   *accesssvc.Service
	bookService     *booksvc.BookService
	sessionService  *sessionsvc.Service
	progressService *progresssvc.ProgressService
	streakService   *streakservice.StreakService
	sqsClient       *sqs.Client
	producer        *kafka.Producer
}

func New(
	access *accesssvc.Service,
	book *booksvc.BookService,
	session *sessionsvc.Service,
	progress *progresssvc.ProgressService,
	streakService *streakservice.StreakService,
	sqsClient *sqs.Client,
	producer *kafka.Producer,
) *Engine {

	return &Engine{
		accessService:   access,
		bookService:     book,
		sessionService:  session,
		progressService: progress,
		streakService:   streakService,
		sqsClient:       sqsClient,
		producer: producer,
	}
}
func (e *Engine) Open(
	ctx context.Context,
	userID string,
	childID string,
	bookID string,
) (*readermodel.ReadingState, error) {

	log.Printf("📖 OPEN: user=%s child=%s book=%s", userID, childID, bookID)

	// =========================
	// LOAD BOOK META
	// =========================

	book, err := e.bookService.GetBookMeta(ctx, bookID)
	if err != nil {
		log.Printf("❌ GetBookMeta failed book=%s err=%v", bookID, err)
		return nil, err
	}

	log.Printf("📘 Book loaded: %s", bookID)

	// =========================
	// ACCESS CHECK
	// =========================

	allowed, preview, err := e.accessService.CanAccessBook(ctx, userID, book)
	if err != nil {
		log.Printf("❌ Access check failed user=%s book=%s err=%v", userID, bookID, err)
		return nil, err
	}

	log.Printf("🔐 Access result: allowed=%v preview=%v", allowed, preview)

	locked := false

var variants []bookdto.ReaderVariant

	if allowed {

		variants, err = e.bookService.GetVariantsWithPages(ctx, bookID)
		if err != nil {
			log.Printf("❌ GetBookVariant failed book=%s err=%v", bookID, err)
			return nil, err
		}

		log.Printf("📚 Full book loaded pages=%d", len(variants))

	} else if preview {

		locked = true

		variants, err = e.bookService.GetVariantsWithPreview(ctx, bookID, 3)
		if err != nil {
			log.Printf("❌ GetBookPreview failed book=%s err=%v", bookID, err)
			return nil, err
		}

		log.Printf("👀 Preview book loaded pages=%d", len(variants))

	} else {
		log.Printf("⛔ Access denied user=%s book=%s", userID, bookID)
		return nil, fmt.Errorf("access denied")
	}

	maxPage := 0
for _, v := range variants {
	maxPage += len(v.Pages)
}

	// =========================
	// SESSION (GET OR CREATE)
	// =========================

	session, err := e.sessionService.GetOrCreateActiveSession(
		ctx,
		userID,
		childID,
		bookID,
		0,
	)

	if err != nil {
		log.Printf("❌ Session error user=%s child=%s book=%s err=%v",
			userID, childID, bookID, err)
		return nil, err
	}

	log.Printf("🟢 Session active: sessionID=%s", session.ID)

	// =========================
	// PUBLISH EVENT
	// =========================

	event := events.Event{
		EventID:   uuid.NewString(),
		Type:      events.SessionStarted,
		SessionID: session.ID,
		UserID:    userID,
		ChildID:   childID,
		BookID:    bookID,
		Page:      0,
		Timestamp: time.Now(),
	}

	e.publish(event)

	log.Printf("📤 Event published: type=%s session=%s", event.Type, session.ID)

	// =========================
	// PROGRESS
	// =========================

	progress, err := e.progressService.GetProgress(ctx, childID, bookID)

	if err != nil {

		if errors.Is(err, repository.ErrNotFound) {

			log.Printf("📊 No progress found → creating child=%s book=%s", childID, bookID)

			progress, err = e.progressService.CreateProgress(
				ctx,
				childID,
				bookID,
				0,
			)

			if err != nil {
				log.Printf("❌ CreateProgress failed err=%v", err)
				return nil, err
			}
		} else {
			log.Printf("❌ GetProgress failed err=%v", err)
			return nil, err
		}
	}

	log.Printf("📊 Progress loaded page=%d percent=%d",
		progress.CurrentPage,
		progress.ProgressPercent,
	)

	// =========================
	// RESPONSE
	// =========================

	log.Printf("✅ OPEN complete session=%s book=%s", session.ID, bookID)

	return &readermodel.ReadingState{
		SessionID: session.ID,

		Book: readermodel.BookResponse{
			Info:  book,
			Variants: variants,
		},

		Reader: readermodel.ReaderProgress{
			CurrentPage:     progress.CurrentPage,
			Completed:       progress.Completed,
			ProgressPercent: progress.ProgressPercent,
		},

		Access: readermodel.ReaderAccess{
			Allowed: allowed,
			Preview: preview,
			Locked:  locked,
			MaxPage: maxPage,
		},

		Features: readermodel.ReaderFeatures{
			Audio:     true,
			Bookmarks: true,
		},
	}, nil
}
func (e *Engine) Update(
	ctx context.Context,
	userID string,
	childID string,
	bookID string,
	page int,
) error {

	// =========================
	// LOAD BOOK META
	// =========================

	book, err := e.bookService.GetBookMeta(
		ctx,
		bookID,
	)

	if err != nil {
		return err
	}

	// =========================
	// ACCESS CHECK
	// =========================

	allowed, preview, err := e.accessService.CanAccessBook(
		ctx,
		userID,
		book,
	)

	if err != nil {
		return err
	}

	// =========================
	// PAGE LIMIT VALIDATION
	// =========================

	if preview && page >= constants.DefaultPreviewPages {
		page = constants.DefaultPreviewPages - 1
	}

	// =========================
	// ACTIVE SESSION
	// =========================

	session, err := e.sessionService.GetActiveSession(
		ctx,
		userID,
		childID,
		bookID,
	)

	if err != nil {
		return err
	}

	if session == nil {
		return fmt.Errorf("no active session")
	}

	// =========================
	// UPDATE SESSION
	// =========================

	err = e.sessionService.UpdateSession(
		ctx,
		session.ID,
		page,
	)

	if err != nil {
		return err
	}
    var prevPage int

if session.EndPage != nil {
    prevPage = *session.EndPage
} else {
    prevPage = 0
}
var totalPages int

	if allowed {

		_, pages, err := e.bookService.GetBook(
			ctx,
			bookID,
		)

		if err != nil {
			return err
		}

		totalPages = len(pages)

	} else {

		totalPages = constants.DefaultPreviewPages
	}

	if totalPages <= 0 {
		totalPages = 1
	}
e.publish(events.Event{
	EventID:   uuid.NewString(),
	Type:      events.ProgressUpdated,
	SessionID: session.ID,
	UserID:    userID,
	ChildID:   childID,
	BookID:    bookID,
	Page:      page,
	PreviousPage: prevPage,
	Timestamp: time.Now(),
	TotalPages: totalPages,
})
	// =========================
	// TOTAL PAGES
	// =========================

	

	// =========================
	// UPDATE PROGRESS
	// =========================

	return e.progressService.UpdateProgress(
		ctx,
		childID,
		bookID,
		page,
		totalPages,
	)
}
func (e *Engine) Close(
	ctx context.Context,
	userID string,
	childID string,
	bookID string,
	page int,
) error {

	// =========================
	// LOAD BOOK META
	// =========================

	book, err := e.bookService.GetBookMeta(
		ctx,
		bookID,
	)

	if err != nil {
		return err
	}

	// =========================
	// ACCESS CHECK
	// =========================

	allowed, preview, err := e.accessService.CanAccessBook(
		ctx,
		userID,
		book,
	)

	if err != nil {
		return err
	}

	// =========================
	// LOAD TOTAL PAGES
	// =========================

	totalPages := constants.DefaultPreviewPages

	if allowed {

		_, pages, err := e.bookService.GetBook(
			ctx,
			bookID,
		)

		if err != nil {
			return err
		}

		totalPages = len(pages)

	} else if preview {

		_, pages, err := e.bookService.GetBookPreview(
			ctx,
			bookID,
		)

		if err != nil {
			return err
		}

		totalPages = len(pages)

	} else {

		return fmt.Errorf("access denied")
	}

	if totalPages <= 0 {
		totalPages = 1
	}

	// =========================
	// PREVIEW LIMIT PROTECTION
	// =========================

	if preview && page >= totalPages {
		page = totalPages - 1
	}

	// =========================
	// ACTIVE SESSION
	// =========================

	session, err := e.sessionService.GetActiveSession(
		ctx,
		userID,
		childID,
		bookID,
	)

	if err != nil {
		return err
	}

	if session == nil {
		return fmt.Errorf("no active session")
	}

	// =========================
	// FINAL PROGRESS UPDATE
	// =========================

	err = e.progressService.UpdateProgress(
		ctx,
		childID,
		bookID,
		page,
		totalPages,
	)

	if err != nil {
		return fmt.Errorf("UpdateProgress: %w", err)
	}
    log.Println("close")
	e.publish(events.Event{
	EventID:   uuid.NewString(),
	Type:      events.SessionEnded,
	SessionID: session.ID,
	UserID:    userID,
	ChildID:   childID,
	BookID:    bookID,
	Page:      page,
	Timestamp: time.Now(),
})
	// =========================
	// END SESSION
	// =========================
    
	return e.sessionService.EndSession(
		ctx,
		session.ID,
		page,
	)
}
func (e *Engine) State(
	ctx context.Context,
	userID string,
	childID string,
	bookID string,
) (*readermodel.ReadingState, error) {

	// =========================
	// LOAD BOOK
	// =========================

	book, err := e.bookService.GetBookMeta(
		ctx,
		bookID,
	)

	if err != nil {
		return nil, err
	}

	// =========================
	// ACCESS CHECK
	// =========================

	allowed, preview, err := e.accessService.CanAccessBook(
		ctx,
		userID,
		book,
	)

	if err != nil {
		return nil, err
	}

	locked := false

	// =========================
	// LOAD CONTENT
	// =========================

	var pages []bookmodel.BookPage

	if allowed {

		_, pages, err = e.bookService.GetBook(
			ctx,
			bookID,
		)

		if err != nil {
			return nil, err
		}

	} else if preview {

		locked = true

		_, pages, err = e.bookService.GetBookPreview(
			ctx,
			bookID,
		)

		if err != nil {
			return nil, err
		}

	} else {

		return nil, fmt.Errorf("access denied")
	}

	// =========================
	// ACTIVE SESSION
	// =========================

	session, err := e.sessionService.GetActiveSession(
		ctx,
		userID,
		childID,
		bookID,
	)

	if err != nil {
		return nil, err
	}

	// no active session
	if session == nil {
		return nil, fmt.Errorf("no active session")
	}

	// =========================
	// LOAD PROGRESS
	// =========================

	progress, _ := e.progressService.GetProgress(
		ctx,
		childID,
		bookID,
	)

	currentPage := 0
	completed := false
	percent := 0

	if progress != nil {
		currentPage = progress.CurrentPage
		completed = progress.Completed
		percent = progress.ProgressPercent
	}

	// =========================
	// RESPONSE
	// =========================

	return &readermodel.ReadingState{
		SessionID: session.ID,

		Book: gin.H{
			"info":  book,
			"pages": pages,
		},

		Reader: readermodel.ReaderProgress{
			CurrentPage:     currentPage,
			Completed:       completed,
			ProgressPercent: percent,
		},

		Access: readermodel.ReaderAccess{
			Allowed: allowed,
			Preview: preview,
			Locked:  locked,
			MaxPage: len(pages),
		},

		Features: readermodel.ReaderFeatures{
			Audio:     true,
			Bookmarks: true,
		},
	}, nil
}
func (e *Engine) publish(event events.Event) {

	b, err := json.Marshal(event)
	if err != nil {
		log.Printf("❌ marshal failed: %v", err)
		return
	}

	if err := e.producer.Publish(
		context.Background(),
		"kids-utopia.events",
		event.SessionID,
		b, // 👈 TEMP FORCE RAW JSON
	); err != nil {
		log.Printf("❌ kafka publish failed: %v", err)
	}
}