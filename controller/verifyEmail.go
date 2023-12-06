package controller

import (
	"html/template" // Import pustaka template HTML
	"kerja-praktek/model"
	"net/http"

	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

func VerifyEmail(db *gorm.DB) echo.HandlerFunc {
	return func(c echo.Context) error {
		token := c.QueryParam("token")

		var user model.User
		result := db.Where("verification_token = ?", token).First(&user)
		if result.Error != nil {
			return c.String(http.StatusUnauthorized, "Invalid verification token")
		}

		user.IsVerified = true
		user.VerificationToken = "" // Setelah verifikasi, hapus token verifikasi
		db.Save(&user)

		// Baca template HTML dari file
		tmpl, err := template.ParseFiles("emails/verification.html")
		if err != nil {
			return c.String(http.StatusInternalServerError, "Internal Server Error")
		}

		// Eksekusi template dan kirimkan sebagai respons
		err = tmpl.Execute(c.Response().Writer, nil)
		if err != nil {
			return c.String(http.StatusInternalServerError, "Internal Server Error")
		}

		return nil
	}
}
