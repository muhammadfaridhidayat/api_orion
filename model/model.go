package model

import (
	"time"
)

// ======================
// USER
// ======================

type User struct {
	ID        int       `gorm:"primaryKey" json:"id"`
	Fullname  string    `json:"fullname" gorm:"type:varchar(255);"`
	Email     string    `json:"email" gorm:"type:varchar(255);not null"`
	Password  string    `json:"password" gorm:"type:varchar(255);not null"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type UserRegister struct {
	Fullname string `json:"fullname" binding:"required"`
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// ======================
// NEW MEMBER
// ======================

type Devision string

const (
	Programming Devision = "PROGRAMMING"
	Electronics Devision = "ELECTRONICS"
	Mechanical  Devision = "MECHANICAL"
)

type Status string

const (
	Pending  Status = "PENDING"
	Verified Status = "VERIFIED"
	Rejected Status = "REJECTED"
)

type NewMember struct {
	ID        int       `gorm:"primaryKey" json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`

	FullName    string    `json:"full_name" gorm:"type:varchar(255);not null"`
	Nim         string    `json:"nim" gorm:"type:varchar(255);uniqueIndex"`
	PhoneNumber string    `json:"phone_number" gorm:"type:varchar(255);not null;default:''"`
	Semester    *string   `json:"semester"`
	Devision    *Devision `json:"devision"`
	Motivation  *string   `json:"motivation"`
	Payment     *string   `json:"payment"`
	Status      *Status   `json:"status" gorm:"default:PENDING"`

	BatchId *int   `json:"batch_id"`
	Batch   *Batch `json:"batch"`
}

type UpdateStatusRequest struct {
	Status string `json:"status" binding:"required"`
}

// ======================
// BATCH
// ======================

type Batch struct {
	ID        int       `gorm:"primaryKey" json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`

	Name      string     `json:"name" gorm:"type:varchar(255);not null"`
	IsActive  *bool      `json:"is_active" gorm:"default:false"`
	StartDate *time.Time `json:"start_date"`
	EndDate   *time.Time `json:"end_date"`

	NewMembers []NewMember `json:"new_members" gorm:"foreignKey:BatchId"`
}

type UpdateActiveStatusRequest struct {
	IsActive bool `json:"is_active"`
}

type RegistrationTrend struct {
	Day         string `json:"day"`
	Programming int    `json:"programming"`
	Electronic  int    `json:"electronic"`
	Mechanic    int    `json:"mechanic"`
}

// ======================
// RESPONSE
// ======================

type SuccessResponse struct {
	Success bool        `json:"success"`
	Status  int         `json:"status"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
	Meta    interface{} `json:"meta,omitempty"`
}

type ErrorResponse struct {
	Success bool              `json:"success"`
	Status  int               `json:"status"`
	Message string            `json:"message"`
	Errors  map[string]string `json:"errors,omitempty"`
}
