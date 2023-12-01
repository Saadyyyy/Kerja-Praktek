package model

import (
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Username          string `json:"Username"`
	Password          string `json:"password"`
	PhoneNumber       string `json:"phone_number"`
	Email             string `json:"email"`
	IsVerified        bool   `gorm:"default:false" json:"is_verified"`
	VerificationToken string `json:"verification_token"`
	ConfirmPassword   string `json:"confirm_password"`
	IsAdmin           bool   `gorm:"default:false" json:"is_admin"` // Menambahkan nilai default

}

// Buat struct untuk permintaan perubahan kata sandi
type ChangePasswordRequest struct {
	CurrentPassword string `json:"currentPassword" binding:"required"`
	NewPassword     string `json:"newPassword" binding:"required"`
	ConfirmPassword string `json:"confirmPassword" binding:"required"` // Tambahkan ConfirmPassword
}
