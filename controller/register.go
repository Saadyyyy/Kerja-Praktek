package controller

import (
	"errors"
	"kerja-praktek/emails"
	"kerja-praktek/helper"
	"kerja-praktek/middleware"
	"kerja-praktek/model"

	"net/http"
	"regexp"

	"github.com/labstack/echo/v4"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

func SignUp(db *gorm.DB, secretKey []byte) echo.HandlerFunc {
	return func(c echo.Context) error {
		// Mendeklarasikan variabel untuk menyimpan data pengguna
		var user model.User
		if err := c.Bind(&user); err != nil {
			errorResponse := helper.ErrorResponse{Code: http.StatusBadRequest, Message: err.Error()}
			return c.JSON(http.StatusBadRequest, errorResponse)
		}

		// Validasi data username
		if len(user.Username) < 5 {
			errorResponse := helper.ErrorResponse{Code: http.StatusBadRequest, Message: "Username must be at least 5 characters"}
			return c.JSON(http.StatusBadRequest, errorResponse)
		}

		// Validasi data password
		if len(user.Password) < 8 || !helper.IsValidPassword(user.Password) {
			errorResponse := helper.ErrorResponse{Code: http.StatusBadRequest, Message: "Password must be at least 8 characters and contain a combination of letters and numbers"}
			return c.JSON(http.StatusBadRequest, errorResponse)
		}

		// Validasi format email
		emailPattern := `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`
		match, _ := regexp.MatchString(emailPattern, user.Email)
		if !match {
			errorResponse := helper.ErrorResponse{Code: http.StatusBadRequest, Message: "Invalid email format"}
			return c.JSON(http.StatusBadRequest, errorResponse)
		}

		// Validasi data phone number
		if !helper.IsValidPhoneNumber(user.PhoneNumber) {
			errorResponse := helper.ErrorResponse{Code: http.StatusBadRequest, Message: "Invalid phone number format"}
			return c.JSON(http.StatusBadRequest, errorResponse)
		}

		if user.Password != user.ConfirmPassword {
			errorResponse := helper.ErrorResponse{Code: http.StatusBadRequest, Message: "Password confirmation does not match"}
			return c.JSON(http.StatusBadRequest, errorResponse)
		}

		// Memeriksa apakah username sudah ada di database
		var existingUser model.User
		result := db.Where("username = ?", user.Username).First(&existingUser)
		if result.Error == nil {
			errorResponse := helper.ErrorResponse{Code: http.StatusConflict, Message: "Username already exists"}
			return c.JSON(http.StatusConflict, errorResponse)
		} else if !errors.Is(result.Error, gorm.ErrRecordNotFound) {
			errorResponse := helper.ErrorResponse{Code: http.StatusInternalServerError, Message: "Failed to check username"}
			return c.JSON(http.StatusInternalServerError, errorResponse)
		}

		// Memeriksa apakah email sudah ada di database
		result = db.Where("email = ?", user.Email).First(&existingUser)
		if result.Error == nil {
			errorResponse := helper.ErrorResponse{Code: http.StatusConflict, Message: "Email already exists"}
			return c.JSON(http.StatusConflict, errorResponse)
		} else if !errors.Is(result.Error, gorm.ErrRecordNotFound) {
			errorResponse := helper.ErrorResponse{Code: http.StatusInternalServerError, Message: "Failed to check email"}
			return c.JSON(http.StatusInternalServerError, errorResponse)
		}

		// Memeriksa apakah nomor telepon sudah ada di database
		result = db.Where("phone_number = ?", user.PhoneNumber).First(&existingUser)
		if result.Error == nil {
			errorResponse := helper.ErrorResponse{Code: http.StatusConflict, Message: "Phone number already exists"}
			return c.JSON(http.StatusConflict, errorResponse)
		} else if !errors.Is(result.Error, gorm.ErrRecordNotFound) {
			errorResponse := helper.ErrorResponse{Code: http.StatusInternalServerError, Message: "Failed to check phone number"}
			return c.JSON(http.StatusInternalServerError, errorResponse)
		}

		// Mengenkripsi password
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
		if err != nil {
			errorResponse := helper.ErrorResponse{Code: http.StatusInternalServerError, Message: "Failed to hash password"}
			return c.JSON(http.StatusInternalServerError, errorResponse)
		}

		// Membuat token verifikasi unik
		uniqueToken := helper.GenerateUniqueToken()
		user.VerificationToken = uniqueToken

		// Menyimpan password yang dienkripsi ke database
		user.Password = string(hashedPassword)
		db.Create(&user)
		user.Password = ""

		// Generate token otentikasi
		tokenString, err := middleware.GenerateToken(user.Username, secretKey)
		if err != nil {
			errorResponse := helper.ErrorResponse{Code: http.StatusInternalServerError, Message: "Failed to generate token"}
			return c.JSON(http.StatusInternalServerError, errorResponse)
		}

		// Mengirim email selamat datang
		if err := emails.SendWelcomeEmail(user.Email, user.Username, uniqueToken); err != nil {
			errorResponse := helper.ErrorResponse{Code: http.StatusInternalServerError, Message: "Failed to send welcome email"}
			return c.JSON(http.StatusInternalServerError, errorResponse)
		}

		// Respon sukses
		response := map[string]interface{}{
			"code":    http.StatusOK,
			"message": "User created successfully, Please check your email to verify your account",
			"token":   tokenString,
			"id":      user.ID,
		}

		return c.JSON(http.StatusOK, response)
	}
}
