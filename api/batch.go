package api

import (
	"api_orion/model"
	"api_orion/service"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type BatchAPI interface {
	CreateBatch(c *gin.Context)
	GetAllBatch(c *gin.Context)
	GetBatchByID(c *gin.Context)
	GetActiveBatch(c *gin.Context)
	Update(c *gin.Context)
	UpdateActiveStatus(c *gin.Context)
	Delete(c *gin.Context)
}

type BatchAPIImpl struct {
	batchService service.BatchService
}

func NewBatchAPI(batchService service.BatchService) BatchAPI {
	return &BatchAPIImpl{batchService: batchService}
}

func (api *BatchAPIImpl) CreateBatch(c *gin.Context) {
	var batch model.Batch
	if err := c.ShouldBindJSON(&batch); err != nil {
		c.JSON(http.StatusBadRequest, model.ErrorResponse{
			Success: false,
			Status:  http.StatusBadRequest,
			Message: "Invalid request body",
			Errors:  map[string]string{"error": err.Error()},
		})
		return
	}

	if err := api.batchService.CreateBatch(&batch); err != nil {
		c.JSON(http.StatusInternalServerError, model.ErrorResponse{
			Success: false,
			Status:  http.StatusInternalServerError,
			Message: "Failed to create batch",
			Errors:  map[string]string{"error": err.Error()},
		})
		return
	}

	c.JSON(http.StatusOK, model.SuccessResponse{
		Success: true,
		Status:  http.StatusOK,
		Message: "Batch created successfully",
		Data: gin.H{
			"batch": batch,
		},
	})
}

func (api *BatchAPIImpl) GetAllBatch(c *gin.Context) {
	batches, err := api.batchService.GetAllBatch()
	if err != nil {
		c.JSON(http.StatusInternalServerError, model.ErrorResponse{
			Success: false,
			Status:  http.StatusInternalServerError,
			Message: "Failed to get all batches",
			Errors:  map[string]string{"error": err.Error()},
		})
		return
	}

	c.JSON(http.StatusOK, model.SuccessResponse{
		Success: true,
		Status:  http.StatusOK,
		Message: "All batches retrieved successfully",
		Data: gin.H{
			"batches": batches,
		},
	})
}

func (api *BatchAPIImpl) GetBatchByID(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, model.ErrorResponse{
			Success: false,
			Status:  http.StatusBadRequest,
			Message: "Invalid request body",
			Errors:  map[string]string{"error": err.Error()},
		})
		return
	}
	batch, err := api.batchService.GetBatchByID(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, model.ErrorResponse{
			Success: false,
			Status:  http.StatusInternalServerError,
			Message: "Failed to get batch",
			Errors:  map[string]string{"error": err.Error()},
		})
		return
	}

	c.JSON(http.StatusOK, model.SuccessResponse{
		Success: true,
		Status:  http.StatusOK,
		Message: "Batch retrieved successfully",
		Data: gin.H{
			"batch": batch,
		},
	})
}

func (api *BatchAPIImpl) GetActiveBatch(c *gin.Context) {
	batch, err := api.batchService.GetActiveBatch()
	if err != nil {
		c.JSON(http.StatusInternalServerError, model.ErrorResponse{
			Success: false,
			Status:  http.StatusInternalServerError,
			Message: "Failed to get active batch",
			Errors:  map[string]string{"error": err.Error()},
		})
		return
	}

	c.JSON(http.StatusOK, model.SuccessResponse{
		Success: true,
		Status:  http.StatusOK,
		Message: "Active batch retrieved successfully",
		Data: gin.H{
			"batch": batch,
		},
	})
}

func (api *BatchAPIImpl) Update(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, model.ErrorResponse{
			Success: false,
			Status:  http.StatusBadRequest,
			Message: "Invalid request body",
			Errors:  map[string]string{"error": err.Error()},
		})
		return
	}
	var batch model.Batch
	if err := c.ShouldBindJSON(&batch); err != nil {
		c.JSON(http.StatusBadRequest, model.ErrorResponse{
			Success: false,
			Status:  http.StatusBadRequest,
			Message: "Invalid request body",
			Errors:  map[string]string{"error": err.Error()},
		})
		return
	}

	if err := api.batchService.Update(id, &batch); err != nil {
		c.JSON(http.StatusInternalServerError, model.ErrorResponse{
			Success: false,
			Status:  http.StatusInternalServerError,
			Message: "Failed to update batch",
			Errors:  map[string]string{"error": err.Error()},
		})
		return
	}

	c.JSON(http.StatusOK, model.SuccessResponse{
		Success: true,
		Status:  http.StatusOK,
		Message: "Batch updated successfully",
	})
}

func (api *BatchAPIImpl) UpdateActiveStatus(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, model.ErrorResponse{
			Success: false,
			Status:  http.StatusBadRequest,
			Message: "Invalid batch ID",
			Errors:  map[string]string{"id": "Invalid batch ID"},
		})
		return
	}

	var req model.UpdateActiveStatusRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, model.ErrorResponse{
			Success: false,
			Status:  http.StatusBadRequest,
			Message: "Invalid request body",
			Errors:  map[string]string{"body": "is_active field is required"},
		})
		return
	}

	if err := api.batchService.UpdateActiveStatus(id, req.IsActive); err != nil {
		c.JSON(http.StatusInternalServerError, model.ErrorResponse{
			Success: false,
			Status:  http.StatusInternalServerError,
			Message: "Failed to update batch active status",
			Errors:  map[string]string{"error": err.Error()},
		})
		return
	}

	c.JSON(http.StatusOK, model.SuccessResponse{
		Success: true,
		Status:  http.StatusOK,
		Message: "Batch active status updated successfully",
		Data: gin.H{
			"id":        id,
			"is_active": req.IsActive,
		},
	})
}

func (api *BatchAPIImpl) Delete(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, model.ErrorResponse{
			Success: false,
			Status:  http.StatusBadRequest,
			Message: "Invalid request body",
			Errors:  map[string]string{"error": err.Error()},
		})
		return
	}
	if err := api.batchService.Delete(id); err != nil {
		c.JSON(http.StatusInternalServerError, model.ErrorResponse{
			Success: false,
			Status:  http.StatusInternalServerError,
			Message: "Failed to delete batch",
			Errors:  map[string]string{"error": err.Error()},
		})
		return
	}

	c.JSON(http.StatusOK, model.SuccessResponse{
		Success: true,
		Status:  http.StatusOK,
		Message: "Batch deleted successfully",
	})
}
