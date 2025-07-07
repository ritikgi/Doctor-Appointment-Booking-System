package main

import (
	"log"
	"os"

	_ "auth/docs"
	"auth/pkg/handler"
	"auth/pkg/middleware"
	"auth/pkg/repository"
	"auth/pkg/service"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func setupDatabase() *gorm.DB {
	// Read DB config from environment variables
	dsn := "host=" + os.Getenv("DB_HOST") +
		" user=" + os.Getenv("DB_USER") +
		" password=" + os.Getenv("DB_PASSWORD") +
		" dbname=" + os.Getenv("DB_NAME") +
		" port=" + os.Getenv("DB_PORT") +
		" sslmode=disable"

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}

	// Auto-migrate the User model (creates table if not exists)
	// if err := db.AutoMigrate(&model.User{}); err != nil {
	// 	log.Fatalf("failed to migrate database: %v", err)
	// }

	return db
}

func main() {
	// Load environment variables from .env file (for local development)
	_ = godotenv.Load()

	db := setupDatabase() // Connect to DB and migrate

	// Initialize repository, service, and handler
	userRepo := &repository.UserRepository{DB: db}
	authService := &service.AuthService{UserRepo: userRepo}
	authHandler := &handler.AuthHandler{Service: authService}

	r := gin.Default()
	r.Use(cors.Default())

	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})

	// Register auth routes
	r.POST("/register", authHandler.Register) // User registration
	r.POST("/login", authHandler.Login)       // User login

	// Protected route example
	r.GET("/me", middleware.JWTAuthMiddleware(), func(c *gin.Context) {
		userID := c.GetInt("user_id")
		email := c.GetString("email")
		role := c.GetString("role")
		c.JSON(200, gin.H{
			"user_id": userID,
			"email":   email,
			"role":    role,
		})
	})

	// Swagger docs endpoint
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	if err := r.Run(":8080"); err != nil {
		log.Fatalf("failed to run server: %v", err)
	}
}
