package main

import (
	"log"

	"os"
	"user/pkg/handler"
	"user/pkg/middleware"
	"user/pkg/repository"
	"user/pkg/service"

	_ "user/docs"

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
	userService := &service.UserService{UserRepo: userRepo}
	userHandler := &handler.UserHandler{Service: userService}

	r := gin.Default()
	r.Use(cors.Default())
	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})

	// Protected user profile endpoints
	r.GET("/me", middleware.JWTAuthMiddleware(), userHandler.GetProfile)
	r.PUT("/me", middleware.JWTAuthMiddleware(), userHandler.UpdateProfile)
	r.GET("/doctors", userHandler.ListDoctors)

	// Swagger docs endpoint
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	if err := r.Run(":8081"); err != nil {
		log.Fatalf("failed to run server: %v", err)
	}
}
