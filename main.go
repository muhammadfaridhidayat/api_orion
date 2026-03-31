package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"api_orion/api"
	"api_orion/db"
	"api_orion/middleware"
	"api_orion/model"
	"api_orion/repo"
	"api_orion/service"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"

	// _ "github.com/lib/pq"
	"gorm.io/gorm"
	// "gorm.io/driver/postgres"
)

type APIHandler struct {
	UserAPIHandler   api.UserAPI
	MemberAPIHandler api.MemberAPI
	BatchAPIHandler  api.BatchAPI
}

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	gin.SetMode(gin.ReleaseMode)

	router := gin.New()

	router.Use(gin.LoggerWithFormatter(func(param gin.LogFormatterParams) string {
		return fmt.Sprintf("[%s] \"%s %s %s\"\n",
			param.TimeStamp.Format(time.RFC822),
			param.Method,
			param.Path,
			param.ErrorMessage,
		)
	}))
	router.Use(gin.Recovery())

	// --- CORS SETUP HERE ---
	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{os.Getenv("CORS_ALLOWED_ORIGINS")},
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))
	// --- END CORS ---

	// Serve uploaded files
	router.Static("/uploads", "./uploads")


	// Get DATABASE_URL
	databaseURL := os.Getenv("DATABASE_URL")
	if databaseURL == "" {
		panic("DATABASE_URL not found in .env file")
	}

	database := db.Postgres{}
	conn, err := database.ConnectUrl(databaseURL)
	if err != nil {
		panic(err)
	}

	// Auto-migrate the User model
	err = conn.AutoMigrate(&model.User{}, &model.NewMember{}, &model.Batch{})
	if err != nil {
		log.Fatal("Failed to auto-migrate User model:", err)
	}

	fmt.Println("Successfully connected to database and auto-migrated User model")

	// Route
	router = RunServer(router, conn)

	// Get Port
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	fmt.Printf("✅ Server is running on port %s\n", port)
	if err := router.Run(":" + port); err != nil {
		panic(err)
	}
}

func RunServer(router *gin.Engine, conn *gorm.DB) *gin.Engine {
	dbConn := conn

	userAPI := api.NewUserAPI(service.NewUserService(repo.NewUserRepo(dbConn)))
	batchServ := service.NewBatchService(repo.NewBatchRepo(dbConn))
	memberAPI := api.NewMemberAPI(service.NewMemberService(repo.NewMemberRepo(dbConn)), batchServ)
	batchAPI := api.NewBatchAPI(batchServ)

	apiHandler := APIHandler{
		UserAPIHandler:   userAPI,
		MemberAPIHandler: memberAPI,
		BatchAPIHandler:  batchAPI,
	}

	apiV1 := router.Group("/api/v1")

	user := apiV1.Group("/user")
	{
		user.POST("/register", apiHandler.UserAPIHandler.Register)
		user.POST("/login", apiHandler.UserAPIHandler.Login)
		user.POST("/logout", apiHandler.UserAPIHandler.Logout)

		user.Use(middleware.Auth())
		user.GET("/profile", apiHandler.UserAPIHandler.GetUserProfile)
	}

	member := apiV1.Group("/member")
	{
		member.POST("/register", apiHandler.MemberAPIHandler.CreateMember)

		member.Use(middleware.Auth())
		member.GET("/all", apiHandler.MemberAPIHandler.GetAllMember)
		member.GET("/:id", apiHandler.MemberAPIHandler.GetMemberByID)
		member.GET("/nim/:nim", apiHandler.MemberAPIHandler.GetMemberByNim)
		member.PUT("/:id", apiHandler.MemberAPIHandler.Update)
		member.PUT("/:id/status", apiHandler.MemberAPIHandler.UpdateStatus)
		member.GET("/trend", apiHandler.MemberAPIHandler.GetRegistrationTrend)
		member.DELETE("/:id", apiHandler.MemberAPIHandler.Delete)
	}

	batch := apiV1.Group("/batch")
	{
		batch.GET("/active", apiHandler.BatchAPIHandler.GetActiveBatch)

		batch.Use(middleware.Auth())
		batch.POST("/create", apiHandler.BatchAPIHandler.CreateBatch)
		batch.GET("/all", apiHandler.BatchAPIHandler.GetAllBatch)
		batch.GET("/:id", apiHandler.BatchAPIHandler.GetBatchByID)
		batch.PUT("/:id", apiHandler.BatchAPIHandler.Update)
		batch.PUT("/:id/active", apiHandler.BatchAPIHandler.UpdateActiveStatus)
		batch.DELETE("/:id", apiHandler.BatchAPIHandler.Delete)
	}

	return router
}
