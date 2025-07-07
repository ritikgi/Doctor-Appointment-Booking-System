package main

import (
	"log"

	"os"
	"schedule/pkg/handler"
	"schedule/pkg/middleware"
	"schedule/pkg/repository"
	"schedule/pkg/service"

	_ "schedule/docs"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
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

	// Auto-migrate the Slot model (creates table if not exists)
	// if err := db.AutoMigrate(&model.Slot{}); err != nil {
	// 	log.Fatalf("failed to migrate database: %v", err)
	// }

	return db
}

func main() {
	db := setupDatabase() // Connect to DB and migrate

	// Initialize repository, service, and handler
	slotRepo := &repository.SlotRepository{DB: db}
	slotService := &service.SlotService{Repo: slotRepo}
	slotHandler := &handler.SlotHandler{Service: slotService}

	r := gin.Default()
	// r.Use(cors.Default())
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:3000"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
	}))

	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})

	// Public: Get slots for a doctor
	r.GET("/slots", slotHandler.GetSlots)

	// Doctor-only: Create, update, delete slots (JWT-protected)
	doctorOnly := func(c *gin.Context) {
		role := c.GetString("role")
		if role != "doctor" {
			c.JSON(403, gin.H{"error": "forbidden: only doctors allowed"})
			c.Abort()
			return
		}
		c.Next()
	}

	r.POST("/slots", middleware.JWTAuthMiddleware(), doctorOnly, slotHandler.CreateSlot)
	r.PUT("/slots/:id", middleware.JWTAuthMiddleware(), doctorOnly, slotHandler.UpdateSlot)
	r.DELETE("/slots/:id", middleware.JWTAuthMiddleware(), doctorOnly, slotHandler.DeleteSlot)
	r.PUT("/slots/:id/book", slotHandler.BookSlot)
	r.PUT("/slots/:id/unbook", slotHandler.UnbookSlot)

	// Swagger docs endpoint
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	if err := r.Run(":8082"); err != nil {
		log.Fatalf("failed to run server: %v", err)
	}
}
