package main

import (
	"log"
	"max/auth/controllers"
	"max/auth/database"
	"max/auth/middleware"
	"max/auth/repositories"
	"max/auth/services"
	"max/auth/utils"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	// 1. Cargar variables de entorno
	if err := godotenv.Load(); err != nil {
		log.Println("No se encontró archivo .env, usando variables de sistema")
	}

	// 2. Conectar BD
	db := database.ConnectDB()
	defer db.Close()

	// 3. Inicializar capas (Inyección de Dependencias)
	userRepo := repositories.NewUserRepository(db)
	authService := services.NewAuthService(userRepo)
	authController := controllers.NewAuthController(authService)

	// 4. Configurar el Router (Gin)
	r := gin.Default()

	// Rutas Públicas
	public := r.Group("/api/auth")
	{
		public.POST("/register", authController.Register)
		public.POST("/login", authController.Login)
	}

	// Rutas Protegidas (Usan el Middleware)
	protected := r.Group("/api/user")
	protected.Use(middleware.RequireAuth())
	{
		// Ruta de prueba para verificar que el JWT funciona
		protected.GET("/profile", func(c *gin.Context) {
			userID, _ := c.Get("userID") // Obtenido desde el middleware
			utils.RespondJSON(c, http.StatusOK, "Perfil protegido cargado", gin.H{"userID": userID})
		})
	}

	// 5. Iniciar Servidor
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Servidor corriendo en puerto %s", port)
	r.Run(":" + port)
}
