package main

import (
	"fmt"
	"log"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"

	"github.com/bellapacx/kids-utopia/internal/auth"

	// BOOKS
	bookhandler "github.com/bellapacx/kids-utopia/internal/books/handler"
	bookrepo "github.com/bellapacx/kids-utopia/internal/books/repository"
	bookroutes "github.com/bellapacx/kids-utopia/internal/books/routes"
	bookservice "github.com/bellapacx/kids-utopia/internal/books/service"

	// ACCESS
	accessmiddleware "github.com/bellapacx/kids-utopia/internal/access/middleware"
	accessservice "github.com/bellapacx/kids-utopia/internal/access/service"

	// SUBSCRIPTION
	subhandler "github.com/bellapacx/kids-utopia/internal/subscriptions/handler"
	subrepo "github.com/bellapacx/kids-utopia/internal/subscriptions/repository"
	subroutes "github.com/bellapacx/kids-utopia/internal/subscriptions/routes"
	subservice "github.com/bellapacx/kids-utopia/internal/subscriptions/service"

	// EDITOR (BOOKS)
	editorhandler "github.com/bellapacx/kids-utopia/internal/books/handler"
	editorroutes "github.com/bellapacx/kids-utopia/internal/books/routes"
	editorservice "github.com/bellapacx/kids-utopia/internal/books/service"

	// USERS
	usersRoutes "github.com/bellapacx/kids-utopia/internal/users"
	usersHandler "github.com/bellapacx/kids-utopia/internal/users/handler"
	usersRepo "github.com/bellapacx/kids-utopia/internal/users/repository"
	usersService "github.com/bellapacx/kids-utopia/internal/users/service"

	// CHILDREN
	childrenRoutes "github.com/bellapacx/kids-utopia/internal/children"
	childrenHandler "github.com/bellapacx/kids-utopia/internal/children/handler"
	childrenRepo "github.com/bellapacx/kids-utopia/internal/children/repository"
	childrenService "github.com/bellapacx/kids-utopia/internal/children/service"

	// PROGRESS
	progressHandler "github.com/bellapacx/kids-utopia/internal/progress/handler"
	progressRepo "github.com/bellapacx/kids-utopia/internal/progress/repository"
	progressRoutes "github.com/bellapacx/kids-utopia/internal/progress/routes"
	progressService "github.com/bellapacx/kids-utopia/internal/progress/service"

	// INFRA
	"github.com/bellapacx/kids-utopia/pkg/config"
	"github.com/bellapacx/kids-utopia/pkg/database"
	"github.com/bellapacx/kids-utopia/pkg/logger"
	"github.com/bellapacx/kids-utopia/pkg/middleware"
	redisClient "github.com/bellapacx/kids-utopia/pkg/redis"
	sqsClient "github.com/bellapacx/kids-utopia/pkg/sqs"
	"github.com/bellapacx/kids-utopia/pkg/storage"

	// NOTIFICATIONS
	"github.com/bellapacx/kids-utopia/internal/notifications/email"
	"github.com/bellapacx/kids-utopia/internal/notifications/otp"
	"github.com/bellapacx/kids-utopia/internal/notifications/sms"
)

func main() {

	cfg := config.Load()
	logger.Init()

	// =========================
	// DB
	// =========================
	dbURL := fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s?sslmode=%s",
		cfg.DBUser,
		cfg.DBPassword,
		cfg.DBHost,
		cfg.DBPort,
		cfg.DBName,
		cfg.DBSSLMode,
	)

	database.Connect(dbURL)
	log.Println("✅ PostgreSQL connected")

	// =========================
	// STORAGE
	// =========================
	storageClient, err := storage.NewS3Storage(
		cfg.S3Endpoint,
		cfg.S3AccessKey,
		cfg.S3SecretKey,
		cfg.S3Bucket,
		cfg.S3PublicURL,
	)
	if err != nil {
		log.Fatal(err)
	}

	// =========================
	// QUEUE
	// =========================
	queue, err := sqsClient.New(cfg.SQSQueueURL, cfg.AWSRegion)
	if err != nil {
		log.Fatal(err)
	}

	// =========================
	// REDIS
	// =========================
	redisClient.Connect(cfg.RedisURL)

	// =========================
	// GIN
	// =========================
	r := gin.Default()

	r.RedirectTrailingSlash = false
// =========================
// HEALTH CHECK
// =========================
r.GET("/health", func(c *gin.Context) {
	c.JSON(200, gin.H{
		"status": "ok",
	})
})
	

r.Use(cors.New(cors.Config{
	AllowOriginFunc: func(origin string) bool {
		return origin == "http://localhost:3000" || origin == ""
	},
	AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
	AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
	AllowCredentials: true,
}))

r.OPTIONS("/*path", func(c *gin.Context) {
	c.Status(204)
})

	// =========================
	// AUTH
	// =========================
	emailSender := email.NewSendGrid(cfg.SendGridAPIKey, cfg.FromEmail)
	smsSender := sms.NewSender()

	otpRouter := otp.NewRouter(emailSender, smsSender)
	otpService := otp.NewService(otpRouter)

	authRepo := &auth.Repository{}
	authService := auth.NewService(authRepo, otpService, cfg.JWTSecret)
	authHandler := auth.NewHandler(authService)

	auth.NewRoutes(authHandler).Register(r.Group("/api/v1"))

	// =========================
	// SUBSCRIPTION + ACCESS
	// =========================
	subRepo := subrepo.New(database.DB)
	subService := subservice.New(subRepo)
	subHandler := subhandler.New(subService)

	accessSvc := accessservice.New(subService)

	// =========================
	// BOOKS
	// =========================
	bookRepo := bookrepo.NewBookRepository()

	bookService := bookservice.NewBookService(
		bookRepo,
		storageClient,
		queue,
		accessSvc,
	)

	bookHandler := bookhandler.NewBookHandler(bookService)

	// =========================
	// ACCESS MIDDLEWARE
	// =========================
	accessMw := accessmiddleware.New(accessSvc, bookRepo)

	// =========================
	// BOOK ROUTES (PROTECTED)
	// =========================
	readerGroup := r.Group("/api/v1/books")

readerGroup.Use(
	middleware.AuthMiddleware(cfg.JWTSecret),
)

bookroutes.RegisterReaderRoutes(
	readerGroup,
	bookHandler,
	accessMw,
)
editorBooks := r.Group("/api/v1/books")

editorBooks.Use(
	middleware.AuthMiddleware(cfg.JWTSecret),
)

editorBooks.Use(
	middleware.RequireRoles("editor", "admin"),
)

bookroutes.RegisterEditorBookRoutes(
	editorBooks,
	bookHandler,
)
	

	// =========================
	// SUBSCRIPTIONS ROUTES
	// =========================
	subroutes.RegisterSubscriptionRoutes(r.Group("/api/v1/subscriptions"), subHandler)

	// =========================
	// EDITOR (ADMIN ONLY)
	// =========================
	editorRepo := bookrepo.NewBookRepository()
	editorPagesRepo := bookrepo.NewBookPagesRepository(database.DB)

	editorService := editorservice.NewEditorService(
		editorRepo,
		editorPagesRepo,
		storageClient,
	)

	editorHandler := editorhandler.NewEditorHandler(editorService)

	editorGroup := r.Group("/api/v1")
	editorGroup.Use(middleware.AuthMiddleware(cfg.JWTSecret))
	editorGroup.Use(middleware.RequireRoles("editor", "admin"))

	editorroutes.RegisterEditorRoutes(editorGroup, editorHandler)

	// =========================
	// USERS
	// =========================
	userRepo := usersRepo.NewUserRepository(database.DB)
	userService := usersService.NewUserService(userRepo)
	userHandler := usersHandler.NewUserHandler(userService)

	usersRoutes.RegisterRoutes(
		r.Group("/api/v1"),
		userHandler,
		middleware.AuthMiddleware(cfg.JWTSecret),
	)

	// =========================
	// CHILDREN
	// =========================
	childRepo := childrenRepo.NewChildRepository(database.DB)
	childService := childrenService.NewChildService(childRepo)
	childHandler := childrenHandler.NewChildHandler(childService)

	childrenRoutes.RegisterRoutes(
		r.Group("/api/v1"),
		childHandler,
		middleware.AuthMiddleware(cfg.JWTSecret),
	)

	// =========================
	// PROGRESS
	// =========================
	progRepo := progressRepo.NewProgressRepository(database.DB)
	progService := progressService.NewProgressService(progRepo)
	progHandler := progressHandler.NewProgressHandler(progService)

	progressRoutes.RegisterRoutes(
		r.Group("/api/v1"),
		progHandler,
		middleware.AuthMiddleware(cfg.JWTSecret),
	)

	// =========================
	// ADMIN
	// =========================
	admin := r.Group("/api/v1/admin")
	admin.Use(
		middleware.AuthMiddleware(cfg.JWTSecret),
		middleware.RoleGuard("admin", "super_admin"),
	)

	admin.GET("/dashboard", func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "admin access granted"})
	})

	// =========================
	// START SERVER
	// =========================
	log.Printf("🚀 Server running on %s", cfg.AppPort)

	if err := r.Run(":" + cfg.AppPort); err != nil {
		log.Fatal(err)
	}
}