package api

import (
	"api_orion/model"
	"api_orion/service"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

type UserAPI interface {
	Register(c *gin.Context)
	Login(c *gin.Context)
	Logout(c *gin.Context)
	GetUserProfile(c *gin.Context)
}

type userAPI struct {
	userService service.UserService
}

func NewUserAPI(userService service.UserService) *userAPI {
	return &userAPI{userService}
}

// ====================
// REGISTRATION
// ====================
func (u *userAPI) Register(c *gin.Context) {
	var req model.UserRegister

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, model.ErrorResponse{
			Success: false,
			Status:  http.StatusBadRequest,
			Message: "Validation failed",
			Errors: map[string]string{
				"body": "Invalid JSON format",
			},
		})
		return
	}

	// Minimal Validation
	errors := make(map[string]string)
	if req.Fullname == "" {
		errors["fullname"] = "Fullname is required"
	}
	if req.Email == "" {
		errors["email"] = "Email is required"
	}
	if req.Password == "" {
		errors["password"] = "Password is required"
	}

	if len(errors) > 0 {
		c.JSON(http.StatusBadRequest, model.ErrorResponse{
			Success: false,
			Status:  http.StatusBadRequest,
			Message: "Validation failed",
			Errors:  errors,
		})
		return
	}

	// Hash password before saving
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, model.ErrorResponse{
			Success: false,
			Status:  http.StatusInternalServerError,
			Message: "Failed to hash password",
			Errors: map[string]string{
				"bcrypt": err.Error(),
			},
		})
		return
	}

	user := model.User{
		Fullname: req.Fullname,
		Email:    req.Email,
		Password: string(hashedPassword),
	}

	if err := u.userService.Register(user); err != nil {
		c.JSON(http.StatusInternalServerError, model.ErrorResponse{
			Success: false,
			Status:  http.StatusInternalServerError,
			Message: "Registration failed",
			Errors: map[string]string{
				"server": err.Error(),
			},
		})
		return
	}

	c.JSON(http.StatusCreated, model.SuccessResponse{
		Success: true,
		Status:  http.StatusCreated,
		Message: "Registration successful",
		Data: gin.H{
			"fullname": req.Fullname,
			"email":    req.Email,
		},
	})
}

// ====================
// LOGIN
// ====================
func (u *userAPI) Login(c *gin.Context) {
	var req model.User

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, model.ErrorResponse{
			Success: false,
			Status:  http.StatusBadRequest,
			Message: "Validation failed",
			Errors: map[string]string{
				"body": "Invalid JSON format",
			},
		})
		return
	}

	// Minimal validation
	errors := make(map[string]string)
	if req.Email == "" {
		errors["email"] = "Email is required"
	}
	if req.Password == "" {
		errors["password"] = "Password is required"
	}
	if len(errors) > 0 {
		c.JSON(http.StatusBadRequest, model.ErrorResponse{
			Success: false,
			Status:  http.StatusBadRequest,
			Message: "Validation failed",
			Errors:  errors,
		})
		return
	}

	token, user, err := u.userService.Login(req)
	if err != nil {
		c.JSON(http.StatusUnauthorized, model.ErrorResponse{
			Success: false,
			Status:  http.StatusUnauthorized,
			Message: "Invalid email or password",
			Errors: map[string]string{
				"auth": err.Error(),
			},
		})
		return
	}

	// Save token to cookie
	c.SetCookie("session_token", *token, int((12 * time.Hour).Seconds()), "/", "", false, true)

	c.JSON(http.StatusOK, model.SuccessResponse{
		Success: true,
		Status:  http.StatusOK,
		Message: "Login successful",
		Data: gin.H{
			"user_id": user.ID,
			"email":   user.Email,
			"token":   token,
		},
	})
}

// ====================
// LOGOUT
// ====================
func (u *userAPI) Logout(c *gin.Context) {

	c.SetCookie(
		"session_token",
		"",
		-1,
		"/",
		"",
		true,
		true,
	)

	c.JSON(http.StatusOK, model.SuccessResponse{
		Success: true,
		Status:  http.StatusOK,
		Message: "Logout successful",
	})
}

// ====================
// GET USER PROFILE
// ====================
func (u *userAPI) GetUserProfile(c *gin.Context) {
	userID, exists := c.Get("id")
	if !exists {
		c.JSON(http.StatusUnauthorized, model.ErrorResponse{
			Success: false,
			Status:  http.StatusUnauthorized,
			Message: "Unauthorized",
			Errors: map[string]string{
				"auth": "User not authenticated",
			},
		})
		return
	}

	id, ok := userID.(int)
	if !ok {
		c.JSON(http.StatusInternalServerError, model.ErrorResponse{
			Success: false,
			Status:  http.StatusInternalServerError,
			Message: "Internal server error",
			Errors: map[string]string{
				"server": "Failed to parse user ID",
			},
		})
		return
	}

	user, err := u.userService.GetProfile(id)
	if err != nil {
		c.JSON(http.StatusNotFound, model.ErrorResponse{
			Success: false,
			Status:  http.StatusNotFound,
			Message: "User not found",
			Errors: map[string]string{
				"user": err.Error(),
			},
		})
		return
	}

	c.JSON(http.StatusOK, model.SuccessResponse{
		Success: true,
		Status:  http.StatusOK,
		Message: "Profile retrieved successfully",
		Data: gin.H{
			"user_id":    user.ID,
			"fullname":   user.Fullname,
			"email":      user.Email,
			"created_at": user.CreatedAt,
			"updated_at": user.UpdatedAt,
		},
	})
}
