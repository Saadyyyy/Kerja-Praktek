package routes

import (
	"io/ioutil"
	"kerja-praktek/controller"
	"net/http"
	"os"

	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

func ServeHTML(c echo.Context) error {
	htmlData, err := ioutil.ReadFile("index.html")
	if err != nil {
		return err
	}
	return c.HTML(http.StatusOK, string(htmlData))
}

func SetupRoutes(e *echo.Echo, db *gorm.DB) {
	// e.Use(Logger())
	godotenv.Load(".env")

	// if err != nil {
	// 	log.Fatal("Error loading .env file")
	// }

	secretKey := []byte(os.Getenv("SECRET_JWT"))

	// Menggunakan routes yang telah dipisahkan
	e.GET("/", ServeHTML)
	e.POST("/register", controller.SignUp(db, secretKey))
	e.GET("/verify", controller.VerifyEmail(db))
	e.POST("/login", controller.SignIn(db, secretKey))
}
