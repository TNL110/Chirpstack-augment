package main

import (
	"log"

	"go-auth-api/internal/auth"
	"go-auth-api/internal/config"
	"go-auth-api/internal/database"
	"go-auth-api/internal/handlers"
	"go-auth-api/internal/middleware"
	"go-auth-api/internal/repository"
	"go-auth-api/internal/service"

	"github.com/gin-gonic/gin"
)

func main() {
	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		log.Fatal("Failed to load config:", err)
	}

	// Connect to database
	db, err := database.Connect(cfg)
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}
	defer db.Close()

	// Connect to database with sqlx for device management
	dbx, err := database.ConnectX(cfg)
	if err != nil {
		log.Fatal("Failed to connect to database with sqlx:", err)
	}
	defer dbx.Close()

	// Initialize services
	jwtService := auth.NewJWTService(cfg.JWTSecret)
	userRepo := repository.NewUserRepository(db)
	chirpStackService := service.NewChirpStackService(cfg, userRepo)
	userService := service.NewUserService(userRepo, jwtService, chirpStackService)
	userHandler := handlers.NewUserHandler(userService)

	// Initialize device management
	deviceRepo := repository.NewDeviceRepository(dbx)
	deviceService := service.NewDeviceService(deviceRepo, userRepo, chirpStackService)
	deviceHandler := handlers.NewDeviceHandler(deviceService)

	// Setup Gin router
	r := gin.Default()

	// Health check endpoint
	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})

	// API routes
	api := r.Group("/api/v1")
	{
		// Public routes
		auth := api.Group("/auth")
		{
			auth.POST("/register", userHandler.Register)
			auth.POST("/login", userHandler.Login)
		}

		// Protected routes
		protected := api.Group("/user")
		protected.Use(middleware.AuthMiddleware(jwtService))
		{
			protected.GET("/profile", userHandler.Profile)
		}

		// User management routes (protected)
		users := api.Group("/users")
		users.Use(middleware.AuthMiddleware(jwtService))
		{
			users.GET("", userHandler.GetAllUsers)        // GET /api/v1/users
			users.GET("/search", userHandler.SearchUsers) // GET /api/v1/users/search
			users.GET("/:id", userHandler.GetUserByID)    // GET /api/v1/users/:id
			users.PUT("/:id", userHandler.UpdateUser)     // PUT /api/v1/users/:id
			users.DELETE("/:id", userHandler.DeleteUser)  // DELETE /api/v1/users/:id
		}

		// Device management routes (protected)
		devices := api.Group("/devices")
		devices.Use(middleware.AuthMiddleware(jwtService))
		{
			// Device version management
			devices.POST("/versions", deviceHandler.CreateDeviceVersion)
			devices.GET("/versions", deviceHandler.GetDeviceVersions)
			devices.GET("/versions/:id", deviceHandler.GetDeviceVersionByID)
			devices.PUT("/versions/:id", deviceHandler.UpdateDeviceVersion)
			devices.DELETE("/versions/:id", deviceHandler.DeleteDeviceVersion)

			// Allowed device management
			devices.POST("/allowed", deviceHandler.CreateAllowedDevice)
			devices.GET("/allowed", deviceHandler.GetAllowedDevices)
			devices.GET("/allowed/:devEUI", deviceHandler.GetAllowedDeviceByDevEUI)
			devices.PUT("/allowed/:devEUI", deviceHandler.UpdateAllowedDevice)
			devices.DELETE("/allowed/:devEUI", deviceHandler.DeleteAllowedDevice)

			// User device management
			devices.POST("", deviceHandler.CreateDevice)       // Create device for authenticated user
			devices.GET("/my", deviceHandler.GetMyDevices)     // Get devices for authenticated user
			devices.GET("/all", deviceHandler.GetAllDevices)   // Get all devices (admin)
			devices.GET("/:id", deviceHandler.GetDeviceByID)   // Get device by ID
			devices.PUT("/:id", deviceHandler.UpdateDevice)    // Update device
			devices.DELETE("/:id", deviceHandler.DeleteDevice) // Delete device
		}
	}

	// Start server
	log.Printf("Server starting on port %s", cfg.Port)
	if err := r.Run(":" + cfg.Port); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}
