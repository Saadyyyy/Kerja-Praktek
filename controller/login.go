package controller

import (
	"errors"
	"kerja-praktek/emails"
	help "kerja-praktek/helper"
	"kerja-praktek/middleware"
	"kerja-praktek/model"
	"net/http"

	"github.com/labstack/echo/v4"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

func SignIn(db *gorm.DB, secretKey []byte) echo.HandlerFunc {
	return func(c echo.Context) error {
		var user model.User
		if err := c.Bind(&user); err != nil {
			errorResponse := help.ErrorResponse{Code: http.StatusBadRequest, Message: err.Error()}
			return c.JSON(http.StatusBadRequest, errorResponse)
		}

		// Validasi apakah username dan password telah diisi
		if user.Username == "" {
			errorResponse := help.ErrorResponse{Code: http.StatusBadRequest, Message: "Username is required"}
			return c.JSON(http.StatusBadRequest, errorResponse)
		}
		if user.Password == "" {
			errorResponse := help.ErrorResponse{Code: http.StatusBadRequest, Message: "Password is required"}
			return c.JSON(http.StatusBadRequest, errorResponse)
		}

		// Mengecek apakah username ada dalam database
		var existingUser model.User
		result := db.Where("username = ?", user.Username).First(&existingUser)
		if result.Error != nil {
			if errors.Is(result.Error, gorm.ErrRecordNotFound) {
				errorResponse := help.ErrorResponse{Code: http.StatusUnauthorized, Message: "Invalid username"}
				return c.JSON(http.StatusUnauthorized, errorResponse)
			} else {
				errorResponse := help.ErrorResponse{Code: http.StatusInternalServerError, Message: "Failed to check username"}
				return c.JSON(http.StatusInternalServerError, errorResponse)
			}
		}

		// Membandingkan password yang dimasukkan dengan password yang di-hash
		err := bcrypt.CompareHashAndPassword([]byte(existingUser.Password), []byte(user.Password))
		if err != nil {
			errorResponse := help.ErrorResponse{Code: http.StatusUnauthorized, Message: "Invalid password"}
			return c.JSON(http.StatusUnauthorized, errorResponse)
		}

		if existingUser.IsAdmin {
			errorResponse := help.ErrorResponse{Code: http.StatusForbidden, Message: "Access denied. Admin cannot use this endpoint."}
			return c.JSON(http.StatusForbidden, errorResponse)
		}

		if !existingUser.IsVerified {
			errorResponse := help.ErrorResponse{Code: http.StatusUnauthorized, Message: "Akun tidak terverifikasi. Harap verifikasi email Anda sebelum login."}
			return c.JSON(http.StatusUnauthorized, errorResponse)
		}

		// Generate JWT token
		tokenString, err := middleware.GenerateToken(existingUser.Username, secretKey)
		if err != nil {
			errorResponse := help.ErrorResponse{Code: http.StatusInternalServerError, Message: "Failed to generate token"}
			return c.JSON(http.StatusInternalServerError, errorResponse)
		}

		// Mengirim email notifikasi
		if err := emails.SendLoginNotification(existingUser.Email, existingUser.Username); err != nil {
			errorResponse := help.ErrorResponse{Code: http.StatusInternalServerError, Message: "Failed to send notification email"}
			return c.JSON(http.StatusInternalServerError, errorResponse)
		}

		// Menyertakan ID pengguna dalam respons
		return c.JSON(http.StatusOK, map[string]interface{}{"code": http.StatusOK, "error": false, "message": "User login successful", "token": tokenString, "id": existingUser.ID})
	}
}
