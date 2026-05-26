package engine

import (
	"context"
	"errors"
	"fmt"

	accesssvc "github.com/bellapacx/kids-utopia/internal/access/service"
	bookmodel "github.com/bellapacx/kids-utopia/internal/books/model"
	booksvc "github.com/bellapacx/kids-utopia/internal/books/service"
	"github.com/bellapacx/kids-utopia/internal/progress/repository"
	progresssvc "github.com/bellapacx/kids-utopia/internal/progress/service"
	"github.com/bellapacx/kids-utopia/internal/reader/constants"
	"github.com/bellapacx/kids-utopia/internal/reader/events"
	readermodel "github.com/bellapacx/kids-utopia/internal/reader/model"
	streakservice "github.com/bellapacx/kids-utopia/internal/reader/streak/service"
	sessionsvc "github.com/bellapacx/kids-utopia/internal/reader_session/service"
	"github.com/gin-gonic/gin"
)

type Engine struct {
	accessService   *accesssvc.Service
	bookService     *booksvc.BookService
	sessionService  *sessionsvc.Service
	progressService *progresssvc.ProgressService
	streakService   *streakservice.StreakService
	eventBus        *events.Bus
}

func New(
	access *accesssvc.Service,
	book *booksvc.BookService,
	session *sessionsvc.Service,
	progress *progresssvc.ProgressService,
	streakService *streakservice.StreakService,
	eventBus *events.Bus,
) *Engine {

	return &Engine{
		accessService:   access,
		bookService:     book,
		sessionService:  session,
		progressService: progress,
		streakService:   streakService,
		eventBus:        eventBus,
	}
}
func (e *Engine) Open(
	ctx context.Context,
	userID string,
	childID string,
	bookID string,
) (*readermodel.ReadingState, error) {

	// =========================
	// LOAD BOOK META
	// =========================

	book, err := e.bookService.GetBookMeta(ctx, bookID)
	if err != nil {
		return nil, err
	}

	// =========================
	// ACCESS CHECK
	// =========================

	allowed, preview, err := e.accessService.CanAccessBook(ctx, userID, book)
	if err != nil {
		return nil, err
	}

	locked := false

	var pages []bookmodel.BookPage

	if allowed {

		_, pages, err = e.bookService.GetBook(ctx, bookID)
		if err != nil {
			return nil, err
		}

	} else if preview {

		locked = true

		_, pages, err = e.bookService.GetBookPreview(ctx, bookID)
		if err != nil {
			return nil, err
		}

	} else {
		return nil, fmt.Errorf("access denied")
	}

	maxPage := len(pages)

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
	return nil, err
}

	// =========================
	// PROGRESS (GET OR CREATE)
	// =========================

	progress, err := e.progressService.GetProgress(ctx, childID, bookID)

	if err != nil {

		if errors.Is(err, repository.ErrNotFound) {

			progress, err = e.progressService.CreateProgress(
				ctx,
				childID,
				bookID,
				0,
			)

			if err != nil {
				return nil, err
			}

		} else {
			return nil, err
		}
	}

	currentPage := progress.CurrentPage
	completed := progress.Completed
	percent := progress.ProgressPercent

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
    
	e.eventBus.Publish(events.Event{
	Type:    events.ProgressUpdated,
	ChildID: childID,
	BookID:  bookID,
	Page:    page,
})
	// =========================
	// TOTAL PAGES
	// =========================

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
		return err
	}

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