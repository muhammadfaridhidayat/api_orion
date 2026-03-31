package api

import (
	"api_orion/model"
	"api_orion/service"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type MemberAPI interface {
	CreateMember(c *gin.Context)
	GetAllMember(c *gin.Context)
	GetMemberByID(c *gin.Context)
	GetMemberByNim(c *gin.Context)
	Update(c *gin.Context)
	UpdateStatus(c *gin.Context)
	GetRegistrationTrend(c *gin.Context)
	Delete(c *gin.Context)
}

type memberAPI struct {
	memberService service.MemberService
	batchService  service.BatchService
}

func NewMemberAPI(memberService service.MemberService, batchService service.BatchService) *memberAPI {
	return &memberAPI{memberService: memberService, batchService: batchService}
}

func (m *memberAPI) CreateMember(c *gin.Context) {
	var req model.NewMember

	// Parse text fields from multipart form
	req.FullName = c.PostForm("full_name")
	req.Nim = c.PostForm("nim")
	req.PhoneNumber = c.PostForm("phone_number")

	if semester := c.PostForm("semester"); semester != "" {
		req.Semester = &semester
	}
	if devision := c.PostForm("devision"); devision != "" {
		d := model.Devision(devision)
		req.Devision = &d
	}
	if motivation := c.PostForm("motivation"); motivation != "" {
		req.Motivation = &motivation
	}
	if batchIDStr := c.PostForm("batch_id"); batchIDStr != "" {
		batchID, err := strconv.Atoi(batchIDStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, model.ErrorResponse{
				Success: false,
				Status:  http.StatusBadRequest,
				Message: "Validation failed",
				Errors: map[string]string{
					"batch_id": "Invalid batch ID",
				},
			})
			return
		}
		req.BatchId = &batchID
	}

	// Validate required fields
	if req.FullName == "" || req.Nim == "" || req.PhoneNumber == "" {
		c.JSON(http.StatusBadRequest, model.ErrorResponse{
			Success: false,
			Status:  http.StatusBadRequest,
			Message: "Validation failed",
			Errors: map[string]string{
				"body": "full_name, nim, and phone_number are required",
			},
		})
		return
	}

	// Check for duplicate NIM before saving the file
	existingMember, _ := m.memberService.GetMemberByNim(req.Nim)
	if existingMember != nil {
		c.JSON(http.StatusConflict, model.ErrorResponse{
			Success: false,
			Status:  http.StatusConflict,
			Message: "Member with this NIM already exists",
			Errors: map[string]string{
				"nim": "NIM already registered",
			},
		})
		return
	}

	// Handle file upload for payment proof
	var savePath string
	file, err := c.FormFile("payment")
	if err == nil {
		// Ensure upload directory exists
		uploadDir := "./uploads/payments"
		if err := os.MkdirAll(uploadDir, os.ModePerm); err != nil {
			c.JSON(http.StatusInternalServerError, model.ErrorResponse{
				Success: false,
				Status:  http.StatusInternalServerError,
				Message: "Failed to create upload directory",
				Errors: map[string]string{
					"server": err.Error(),
				},
			})
			return
		}

		// Generate unique filename
		ext := strings.ToLower(filepath.Ext(file.Filename))
		uniqueName := fmt.Sprintf("%s%s", uuid.New().String(), ext)
		savePath = filepath.Join(uploadDir, uniqueName)

		if err := c.SaveUploadedFile(file, savePath); err != nil {
			c.JSON(http.StatusInternalServerError, model.ErrorResponse{
				Success: false,
				Status:  http.StatusInternalServerError,
				Message: "Failed to save payment file",
				Errors: map[string]string{
					"payment": err.Error(),
				},
			})
			return
		}

		paymentURL := fmt.Sprintf("/uploads/payments/%s", uniqueName)
		req.Payment = &paymentURL
	}

	if err := m.memberService.CreateMember(&req); err != nil {
		// Clean up uploaded file if DB insert fails
		if savePath != "" {
			os.Remove(savePath)
		}
		c.JSON(http.StatusInternalServerError, model.ErrorResponse{
			Success: false,
			Status:  http.StatusInternalServerError,
			Message: "Failed to create member",
			Errors: map[string]string{
				"server": err.Error(),
			},
		})
		return
	}

	c.JSON(http.StatusCreated, model.SuccessResponse{
		Success: true,
		Status:  http.StatusCreated,
		Message: "Member created successfully",
		Data: gin.H{
			"member": req,
		},
	})
}

func (m *memberAPI) GetAllMember(c *gin.Context) {
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	query := c.Query("query")

	members, err := m.memberService.GetAllMember(limit, page, query)
	if err != nil {
		c.JSON(http.StatusInternalServerError, model.ErrorResponse{
			Success: false,
			Status:  http.StatusInternalServerError,
			Message: "Failed to get all members",
			Errors: map[string]string{
				"server": err.Error(),
			},
		})
		return
	}

	c.JSON(http.StatusOK, model.SuccessResponse{
		Success: true,
		Status:  http.StatusOK,
		Message: "All members retrieved successfully",
		Data: gin.H{
			"members": members,
		},
	})
}

func (m *memberAPI) GetMemberByID(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, model.ErrorResponse{
			Success: false,
			Status:  http.StatusBadRequest,
			Message: "Validation failed",
			Errors: map[string]string{
				"id": "Invalid member ID",
			},
		})
		return
	}

	member, err := m.memberService.GetMemberByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, model.ErrorResponse{
			Success: false,
			Status:  http.StatusNotFound,
			Message: "Member not found",
			Errors: map[string]string{
				"member": err.Error(),
			},
		})
		return
	}

	c.JSON(http.StatusOK, model.SuccessResponse{
		Success: true,
		Status:  http.StatusOK,
		Message: "Member retrieved successfully",
		Data: gin.H{
			"member": member,
		},
	})
}

func (m *memberAPI) GetMemberByNim(c *gin.Context) {
	nim := c.Param("nim")

	member, err := m.memberService.GetMemberByNim(nim)
	if err != nil {
		c.JSON(http.StatusNotFound, model.ErrorResponse{
			Success: false,
			Status:  http.StatusNotFound,
			Message: "Member not found",
			Errors: map[string]string{
				"member": err.Error(),
			},
		})
		return
	}

	c.JSON(http.StatusOK, model.SuccessResponse{
		Success: true,
		Status:  http.StatusOK,
		Message: "Member retrieved successfully",
		Data: gin.H{
			"member": member,
		},
	})
}

func (m *memberAPI) Update(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, model.ErrorResponse{
			Success: false,
			Status:  http.StatusBadRequest,
			Message: "Validation failed",
			Errors: map[string]string{
				"id": "Invalid member ID",
			},
		})
		return
	}

	var req model.NewMember
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

	if err := m.memberService.Update(id, &req); err != nil {
		c.JSON(http.StatusInternalServerError, model.ErrorResponse{
			Success: false,
			Status:  http.StatusInternalServerError,
			Message: "Failed to update member",
			Errors: map[string]string{
				"server": err.Error(),
			},
		})
		return
	}

	c.JSON(http.StatusOK, model.SuccessResponse{
		Success: true,
		Status:  http.StatusOK,
		Message: "Member updated successfully",
		Data: gin.H{
			"member": req,
		},
	})
}

func (m *memberAPI) UpdateStatus(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, model.ErrorResponse{
			Success: false,
			Status:  http.StatusBadRequest,
			Message: "Validation failed",
			Errors: map[string]string{
				"id": "Invalid member ID",
			},
		})
		return
	}

	var req model.UpdateStatusRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, model.ErrorResponse{
			Success: false,
			Status:  http.StatusBadRequest,
			Message: "Validation failed",
			Errors: map[string]string{
				"body": "status is required",
			},
		})
		return
	}

	// Validate status value
	status := model.Status(req.Status)
	if status != model.Pending && status != model.Verified && status != model.Rejected {
		c.JSON(http.StatusBadRequest, model.ErrorResponse{
			Success: false,
			Status:  http.StatusBadRequest,
			Message: "Validation failed",
			Errors: map[string]string{
				"status": "Status must be one of: PENDING, VERIFIED, REJECTED",
			},
		})
		return
	}

	if err := m.memberService.UpdateStatus(id, req.Status); err != nil {
		c.JSON(http.StatusInternalServerError, model.ErrorResponse{
			Success: false,
			Status:  http.StatusInternalServerError,
			Message: "Failed to update registration status",
			Errors: map[string]string{
				"server": err.Error(),
			},
		})
		return
	}

	c.JSON(http.StatusOK, model.SuccessResponse{
		Success: true,
		Status:  http.StatusOK,
		Message: "Registration status updated successfully",
		Data: gin.H{
			"id":     id,
			"status": req.Status,
		},
	})
}

func (m *memberAPI) GetRegistrationTrend(c *gin.Context) {
	// Get the active batch
	activeBatch, err := m.batchService.GetActiveBatch()
	if err != nil {
		c.JSON(http.StatusNotFound, model.ErrorResponse{
			Success: false,
			Status:  http.StatusNotFound,
			Message: "No active batch found",
			Errors: map[string]string{
				"batch": "No active batch is currently set",
			},
		})
		return
	}

	trends, err := m.memberService.GetRegistrationTrend(activeBatch.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, model.ErrorResponse{
			Success: false,
			Status:  http.StatusInternalServerError,
			Message: "Failed to get registration trend",
			Errors: map[string]string{
				"server": err.Error(),
			},
		})
		return
	}

	c.JSON(http.StatusOK, model.SuccessResponse{
		Success: true,
		Status:  http.StatusOK,
		Message: "Registration trend retrieved successfully",
		Data: gin.H{
			"trends": trends,
		},
	})
}

func (m *memberAPI) Delete(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, model.ErrorResponse{
			Success: false,
			Status:  http.StatusBadRequest,
			Message: "Validation failed",
			Errors: map[string]string{
				"id": "Invalid member ID",
			},
		})
		return
	}

	if err := m.memberService.Delete(id); err != nil {
		c.JSON(http.StatusInternalServerError, model.ErrorResponse{
			Success: false,
			Status:  http.StatusInternalServerError,
			Message: "Failed to delete member",
			Errors: map[string]string{
				"server": err.Error(),
			},
		})
		return
	}

	c.JSON(http.StatusOK, model.SuccessResponse{
		Success: true,
		Status:  http.StatusOK,
		Message: "Member deleted successfully",
	})
}
