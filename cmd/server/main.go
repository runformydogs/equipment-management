package main

import (
	"equipment-management/internal/controller"
	"equipment-management/internal/middleware"
	"equipment-management/internal/service"
	"equipment-management/pkg/auth"
	"github.com/joho/godotenv"
	"log"
	"net/http"

	"equipment-management/internal/config"
	"equipment-management/internal/models"
	"equipment-management/internal/repository"
	"github.com/gin-gonic/gin"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Print("No .env file found")
	}

	cfg := config.LoadConfig()

	if err := repository.InitDB(cfg); err != nil {
		log.Fatal("Database init error: ", err)
	}

	db := repository.DB
	if err := db.AutoMigrate(
		&models.User{},
		&models.Device{},
		&models.NetworkNode{},
	); err != nil {
		log.Fatal("Migration failed: ", err)
	}
	log.Println("Migration completed")

	var userCount int64
	if err := db.Model(&models.User{}).Count(&userCount).Error; err != nil {
		log.Fatal("Failed to count users: ", err)
	}

	if userCount == 0 {
		hashedPassword, err := auth.HashPassword("admin123")
		if err != nil {
			log.Fatal("Failed to hash password: ", err)
		}

		adminUser := models.User{
			Login:    "admin",
			Password: hashedPassword,
			Role:     "admin",
		}

		if err := db.Create(&adminUser).Error; err != nil {
			log.Fatal("Failed to create admin user: ", err)
		}
		log.Println("Admin user created (login: admin, password: admin123)")
	}

	deviceRepo := repository.NewDeviceRepository(repository.DB)

	deviceService := service.NewDeviceService(deviceRepo)

	deviceController := controller.NewDeviceController(deviceService)

	r := gin.Default()

	r.POST("/login", controller.Login)

	authGroup := r.Group("/")
	authGroup.Use(middleware.AuthMiddleware())
	{
		authGroup.GET("/protected", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{"message": "You are authenticated!"})
		})

		deviceGroup := authGroup.Group("/devices")
		{
			deviceGroup.GET("", deviceController.GetAllDevices)
			deviceGroup.GET("/:id", deviceController.GetDevice)

			adminDeviceGroup := deviceGroup.Group("")
			adminDeviceGroup.Use(middleware.RoleMiddleware("admin"))
			{
				deviceGroup.POST("", deviceController.CreateDevice)
				deviceGroup.PUT("/:id", deviceController.UpdateDevice)
				deviceGroup.DELETE("/:id", deviceController.DeleteDevice)
			}
		}
	}
}
