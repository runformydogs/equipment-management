package main

import (
	"equipment-management/internal/controller"
	"equipment-management/internal/middleware"
	"equipment-management/internal/service"
	"equipment-management/pkg/auth"
	"fmt"
	"github.com/joho/godotenv"
	"gorm.io/gorm"
	"log"
	"net/http"
	"time"

	"equipment-management/internal/config"
	"equipment-management/internal/models"
	"equipment-management/internal/repository"
	"github.com/gin-contrib/cors"
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

	if cfg.SeedTestData {
		log.Println("SeedTestData flag enabled, creating test users...")

		var userCount int64
		if err := db.Model(&models.User{}).Count(&userCount).Error; err != nil {
			log.Fatal("Failed to count users: ", err)
		}

		if userCount == 0 {
			hashedPasswordAdmin, _ := auth.HashPassword("admin123")

			adminUser := models.User{
				Login:    "admin",
				Password: hashedPasswordAdmin,
				Role:     "admin",
			}

			if err := db.Create(&adminUser).Error; err != nil {
				log.Fatal("Failed to create admin user: ", err)
			}
			log.Println("Admin user created (login: admin, password: admin123)")

			hashedPasswordViewer, _ := auth.HashPassword("viewer123")

			viewerUser := models.User{
				Login:    "viewer",
				Password: hashedPasswordViewer,
				Role:     "viewer",
			}

			if err := db.Create(&viewerUser).Error; err != nil {
				log.Fatal("Failed to create viewer user: ", err)
			}
			log.Println("Viewer user created (login: viewer, password: viewer123)")
		}

		var nodeCount int64
		if err := db.Model(&models.NetworkNode{}).Count(&nodeCount).Error; err != nil {
			log.Fatal("Failed to count nodes: ", err)
		}

		if nodeCount == 0 {
			log.Println("Seeding test data...")
			if err := seedTestData(db); err != nil {
				log.Fatal("Failed to seed test data: ", err)
			}
			log.Println("Test data seeded successfully")
		}
	}

	deviceRepo := repository.NewDeviceRepository(repository.DB)
	networkNodeRepo := repository.NewNetworkNodeRepository(db)

	deviceService := service.NewDeviceService(deviceRepo)
	networkNodeService := service.NewNetworkNodeService(networkNodeRepo)

	deviceController := controller.NewDeviceController(deviceService)
	networkNodeController := controller.NewNetworkNodeController(networkNodeService)

	r := gin.Default()

	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:63342", "http://localhost:5500", "http://localhost:8080", "http://localhost"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	r.Use(func(c *gin.Context) {
		log.Printf("Запрос: %s %s\n", c.Request.Method, c.Request.URL)
		c.Next()
	})

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
				adminDeviceGroup.POST("", deviceController.CreateDevice)
				adminDeviceGroup.PUT("/:id", deviceController.UpdateDevice)
				adminDeviceGroup.DELETE("/:id", deviceController.DeleteDevice)
			}
		}

		nodeGroup := authGroup.Group("/network-nodes")
		{
			nodeGroup.GET("/tree", networkNodeController.GetFullTree)
			nodeGroup.GET("", networkNodeController.GetAllNodes)
			nodeGroup.GET("/:id", networkNodeController.GetNode)

			adminNodeGroup := nodeGroup.Group("")
			adminNodeGroup.Use(middleware.RoleMiddleware("admin"))
			{
				adminNodeGroup.POST("", networkNodeController.CreateNode)
				adminNodeGroup.PUT("/:id", networkNodeController.UpdateNode)
				adminNodeGroup.DELETE("/:id", networkNodeController.DeleteNode)
			}
		}
	}

	log.Printf("Server starting on :%s...\n", cfg.ServerPort)
	if err := r.Run(":" + cfg.ServerPort); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}

func seedTestData(db *gorm.DB) error {
	rootNodes := []models.NetworkNode{
		{Name: "Главный офис", Description: "Центральный узел сети"},
		{Name: "Филиал Восток", Description: "Восточное подразделение"},
		{Name: "Филиал Запад", Description: "Западное подразделение"},
		{Name: "Серверная", Description: "Основная серверная комната"},
		{Name: "Резервный ЦОД", Description: "Центр обработки данных"},
	}

	for i := range rootNodes {
		if err := db.Create(&rootNodes[i]).Error; err != nil {
			return err
		}
	}

	var createChildren func(parent *models.NetworkNode, depth int) error
	createChildren = func(parent *models.NetworkNode, depth int) error {
		if depth > 3 {
			return nil
		}

		childCount := 2
		if depth == 1 {
			childCount = 3
		}

		for i := 1; i <= childCount; i++ {
			child := models.NetworkNode{
				Name:        fmt.Sprintf("Дочерний узел %d-%d", parent.ID, i),
				Description: fmt.Sprintf("Дочерний элемент узла %s", parent.Name),
				ParentID:    &parent.ID,
			}

			if err := db.Create(&child).Error; err != nil {
				return err
			}

			devices := []models.Device{
				{
					Type:   "Компьютер",
					Vendor: "Dell",
					Model:  fmt.Sprintf("OptiPlex %d", i),
					Serial: fmt.Sprintf("DL-%d-%d", parent.ID, i),
				},
				{
					Type:   "Принтер",
					Vendor: "HP",
					Model:  fmt.Sprintf("LaserJet %d", i),
					Serial: fmt.Sprintf("HP-%d-%d", parent.ID, i),
				},
			}

			for j := range devices {
				devices[j].NetworkNodeID = &child.ID
				if err := db.Create(&devices[j]).Error; err != nil {
					return err
				}
			}

			if err := createChildren(&child, depth+1); err != nil {
				return err
			}
		}
		return nil
	}

	for i := range rootNodes {
		if err := createChildren(&rootNodes[i], 1); err != nil {
			return err
		}
	}

	for i := range rootNodes {
		device := models.Device{
			Type:          "Сервер",
			Vendor:        "IBM",
			Model:         fmt.Sprintf("System X-%d", i),
			Serial:        fmt.Sprintf("IBM-%d", i),
			NetworkNodeID: &rootNodes[i].ID,
		}
		if err := db.Create(&device).Error; err != nil {
			return err
		}
	}

	return nil
}
