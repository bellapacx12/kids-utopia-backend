package main

import (
	"fmt"
	"log"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"

	"github.com/bellapacx/kids-utopia/internal/auth"
	"github.com/bellapacx/kids-utopia/internal/books/handler"
	"github.com/bellapacx/kids-utopia/internal/books/repository"
	"github.com/bellapacx/kids-utopia/internal/books/service"

	"github.com/bellapacx/kids-utopia/internal/books/routes"

	accessmiddleware "github.com/bellapacx/kids-utopia/internal/access/middleware"
	accessservice "github.com/bellapacx/kids-utopia/internal/access/service"
	childrenRoutes "github.com/bellapacx/kids-utopia/internal/children"
	childrenHandler "github.com/bellapacx/kids-utopia/internal/children/handler"
	childrenRepo "github.com/bellapacx/kids-utopia/internal/children/repository"
	childrenService "github.com/bellapacx/kids-utopia/internal/children/service"
	"github.com/bellapacx/kids-utopia/internal/notifications/email"
	"github.com/bellapacx/kids-utopia/internal/notifications/otp"
	"github.com/bellapacx/kids-utopia/internal/notifications/sms"
	progressHandler "github.com/bellapacx/kids-utopia/internal/progress/handler"
	progressRepo "github.com/bellapacx/kids-utopia/internal/progress/repository"
	progressRoutes "github.com/bellapacx/kids-utopia/internal/progress/routes"
	progressService "github.com/bellapacx/kids-utopia/internal/progress/service"
	usersRoutes "github.com/bellapacx/kids-utopia/internal/users"
	usersHandler "github.com/bellapacx/kids-utopia/internal/users/handler"
	usersRepo "github.com/bellapacx/kids-utopia/internal/users/repository"
	usersService "github.com/bellapacx/kids-utopia/internal/users/service"

	editorroutes "github.com/bellapacx/kids-utopia/internal/books/routes"
	subscriptionhandler "github.com/bellapacx/kids-utopia/internal/subscriptions/handler"
	subscriptionrepository "github.com/bellapacx/kids-utopia/internal/subscriptions/repository"
	subscriptionroutes "github.com/bellapacx/kids-utopia/internal/subscriptions/routes"
	subscriptionservice "github.com/bellapacx/kids-utopia/internal/subscriptions/service"
	"github.com/bellapacx/kids-utopia/pkg/config"
	"github.com/bellapacx/kids-utopia/pkg/contextkeys"
	"github.com/bellapacx/kids-utopia/pkg/database"
	"github.com/bellapacx/kids-utopia/pkg/logger"
	"github.com/bellapacx/kids-utopia/pkg/middleware"
	"github.com/bellapacx/kids-utopia/pkg/storage"

	redisClient "github.com/bellapacx/kids-utopia/pkg/redis"
	sqsClient "github.com/bellapacx/kids-utopia/pkg/sqs"
)

func main() {

	// ================================
	// LOAD CONFIG
	// ================================

	cfg := config.Load()
    
	logger.Init()
	// ================================
	// CONNECT POSTGRES (NEON)
	// ================================

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

	// storage
     storageClient, err := storage.NewS3Storage(
	cfg.S3Endpoint,
	cfg.S3AccessKey,
	cfg.S3SecretKey,
	cfg.S3Bucket,
	cfg.S3PublicURL,
)

if err != nil {
	log.Fatal("storage init error:", err)
}
log.Println("RAW S3_ENDPOINT =", cfg.S3Endpoint)

// kafka
//kafkaProducer := kafka.NewProducer([]string{"kafka:9092"})
// ================================
	// AWS SQS MODULE
	// ================================
	queue, err := sqsClient.New(
	cfg.SQSQueueURL,
	cfg.AWSRegion,
)
if err != nil {
	log.Fatal(err)
}
log.Println("✅ SQS initialized:", queue != nil)
	// ================================
	// CONNECT REDIS (UPSTASH)
	// ================================

	redisClient.Connect(cfg.RedisURL)

	log.Println("✅ Redis connected")

	// ================================
	// INIT GIN
	// ================================

	r := gin.Default()
	

r.Use(cors.New(cors.Config{
    AllowOrigins: []string{
        "http://localhost:3000",
    },
    AllowMethods: []string{
        "GET",
        "POST",
        "PUT",
        "DELETE",
        "OPTIONS",
    },
    AllowHeaders: []string{
        "Origin",
        "Content-Type",
        "Authorization",
    },
    AllowCredentials: true,
    MaxAge: 12 * time.Hour,
}))
     
	api := r.Group("/api/v1")
	// ================================
	// HEALTH CHECK
	// ================================

	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status": "ok",
		})
	})
 //sesSender, err := email.NewSESSender(
//	cfg.AWSRegion,
//	cfg.SESFromEmail,
//)
//if err != nil {
//	log.Fatal(err)
//}
// 4. Init SendGrid
	emailSender := email.NewSendGrid(
		cfg.SendGridAPIKey,
		cfg.FromEmail,
	)
	
smsSender := sms.NewSender()
otpRouter := otp.NewRouter(emailSender, smsSender)
otpService := otp.NewService(otpRouter)
	// ================================
	// AUTH MODULE
	// ================================

	authRepo := &auth.Repository{}

	authService := auth.NewService(
		authRepo,
		otpService,
		cfg.JWTSecret,
	)

	authHandler := auth.NewHandler(authService)

	// ================================
	// AUTH ROUTES
	// ================================


// 👇 register routes via module
authRoutes := auth.NewRoutes(authHandler)
authRoutes.Register(r.Group("/api/v1"))
// ================================
	// Books ROUTES
	// ================================


    

bookRepo := repository.NewBookRepository()

bookService := service.NewBookService(
	bookRepo,
	storageClient,
	queue,
)

bookHandler := handler.NewBookHandler(bookService)

// =========================
// SUBSCRIPTIONS
// =========================

subRepo := subscriptionrepository.New(database.DB)

subService := subscriptionservice.New(subRepo)

subHandler := subscriptionhandler.New(subService)
subscriptionroutes.RegisterSubscriptionRoutes(api, subHandler)
// =========================
// ACCESS CONTROL
// =========================

accessSvc := accessservice.New(subService)

accessMw := accessmiddleware.New(
	accessSvc,
	bookRepo,
)

// =========================
// BOOK ROUTES
// =========================

readerGroup := r.Group("/api/v1/books")

readerGroup.Use(
	middleware.AuthMiddleware(cfg.JWTSecret),
)

readerGroup.Use(
	accessMw.CheckBookAccess(),
)

routes.RegisterBookRoutes(
	readerGroup,
	bookHandler,
)
bookPagesRepo := repository.NewBookPagesRepository(database.DB)

editorService := service.NewEditorService(
	bookRepo,
	bookPagesRepo,
	storageClient,
)

editorHandler := handler.NewEditorHandler(editorService)

editorGroup := r.Group("/api/v1/editor")

editorGroup.Use(
	middleware.AuthMiddleware(cfg.JWTSecret),
)

editorGroup.Use(
	middleware.RequireRoles("editor", "admin"),
)

editorroutes.RegisterEditorRoutes(
	editorGroup,
	editorHandler,
)
editorroutes.RegisterEditorRoutes(api, editorHandler)
// ================================
// USERS MODULE
// ================================

userRepo := usersRepo.NewUserRepository(database.DB)

userService := usersService.NewUserService(userRepo)

userHandler := usersHandler.NewUserHandler(userService)


usersRoutes.RegisterRoutes(
	api,
	userHandler,
	middleware.AuthMiddleware(cfg.JWTSecret),
)
// ================================
// CHILDREN MODULE
// ================================

childRepo := childrenRepo.NewChildRepository(database.DB)

childService := childrenService.NewChildService(childRepo)

childHandler := childrenHandler.NewChildHandler(childService)



childrenRoutes.RegisterRoutes(
	api,
	childHandler,
	middleware.AuthMiddleware(cfg.JWTSecret),
)
// ================================
// PROGRESS MODULE
// ================================

progRepo := progressRepo.NewProgressRepository(database.DB)

progService := progressService.NewProgressService(progRepo)

progHandler := progressHandler.NewProgressHandler(progService)

progressRoutes.RegisterRoutes(
	api,
	progHandler,
	middleware.AuthMiddleware(cfg.JWTSecret),
)
	// ================================
	// PROTECTED ROUTES
	// ================================

	protected := r.Group("/api/v1")

	protected.Use(
		middleware.AuthMiddleware(cfg.JWTSecret),
	)

	{
		protected.GET("/me", func(c *gin.Context) {

	c.JSON(200, gin.H{
		"user_id": c.GetString(contextkeys.UserID),
		"role":    c.GetString(contextkeys.Role),
	})
})
	}

	// ================================
	// ADMIN ROUTES
	// ================================

	admin := r.Group("/api/v1/admin")

	admin.Use(
		middleware.AuthMiddleware(cfg.JWTSecret),
		middleware.RoleGuard(
			"admin",
			"super_admin",
		),
	)

	{
		admin.GET("/dashboard", func(c *gin.Context) {

			c.JSON(200, gin.H{
				"message": "admin access granted",
			})
		})
	}

	// ================================
	// START SERVER
	// ================================

	log.Printf("🚀 Server running on port %s", cfg.AppPort)

	if err := r.Run(":" + cfg.AppPort); err != nil {
		log.Fatal(err)
	}
}