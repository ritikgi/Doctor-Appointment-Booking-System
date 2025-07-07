package main

import (
	"log"

	"appointment/pkg/handler"
	"appointment/pkg/middleware"
	"appointment/pkg/repository"
	"appointment/pkg/service"
	"os"

	_ "appointment/docs"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func setupDatabase() *gorm.DB {
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

	// Auto-migrate the Appointment model (creates table if not exists)
	// if err := db.AutoMigrate(&model.Appointment{}); err != nil {
	// 	log.Fatalf("failed to migrate database: %v", err)
	// }

	return db
}

func main() {
	db := setupDatabase() // Connect to DB and migrate

	appRepo := &repository.AppointmentRepository{DB: db}
	appService := &service.AppointmentService{Repo: appRepo}
	appHandler := &handler.AppointmentHandler{Service: appService}

	r := gin.Default()
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

	r.POST("/appointments", middleware.JWTAuthMiddleware(), appHandler.CreateAppointment)
	r.GET("/appointments", middleware.JWTAuthMiddleware(), appHandler.GetAppointments)
	r.PUT("/appointments/:id", middleware.JWTAuthMiddleware(), appHandler.UpdateAppointment)
	r.DELETE("/appointments/:id", middleware.JWTAuthMiddleware(), appHandler.DeleteAppointment)

	// Swagger docs endpoint
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	if err := r.Run(":8083"); err != nil {
		log.Fatalf("failed to run server: %v", err)
	}
}
