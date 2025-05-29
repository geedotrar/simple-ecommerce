package handlers

import (
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strconv"

	"product-service/internal/middleware"
	"product-service/internal/models"
	"product-service/internal/service"
	"product-service/pkg/helpers"

	"github.com/gin-gonic/gin"
)

type ProductHandler interface {
	GetAllProducts(ctx *gin.Context)
	CreateProduct(ctx *gin.Context)
	UpdateProductStatus(ctx *gin.Context)
	DeleteProduct(ctx *gin.Context)
	UpdateProduct(ctx *gin.Context)
}

type productHandlerImpl struct {
	service   service.ProductService
	uploadDir string
}

func NewproductHandler(service service.ProductService, uploadDir string) *productHandlerImpl {
	helpers.InitValidator()
	return &productHandlerImpl{service, uploadDir}
}

func (h *productHandlerImpl) GetAllProducts(c *gin.Context) {
	ctx := c.Request.Context()

	claimsRaw, exists := c.Get("claims")
	if !exists {
		c.JSON(http.StatusUnauthorized, models.ErrorResponse(http.StatusUnauthorized, "Unauthorized", "Missing claims"))
		return
	}

	claims, ok := claimsRaw.(struct {
		middleware.UserClaims
		Permissions []string
	})
	if !ok {
		c.JSON(http.StatusUnauthorized, models.ErrorResponse(http.StatusUnauthorized, "Unauthorized", "Invalid claims format"))
		return
	}

	// log.Printf("User permissions: %v", claims.Permissions)

	pageStr := c.DefaultQuery("page", "1")
	search := c.DefaultQuery("search", "")
	statusStr := c.Query("status")

	page, _ := strconv.Atoi(pageStr)
	if page < 1 {
		page = 1
	}
	limit := 15
	offset := (page - 1) * limit

	var (
		products []models.Product
		total    int64
		err      error
	)

	var status *int
	if statusStr != "" {
		statusVal, err := strconv.Atoi(statusStr)
		if err != nil || (statusVal != 0 && statusVal != 1) {
			c.JSON(http.StatusBadRequest, models.ErrorResponse(http.StatusBadRequest, "Invalid status filter", "Status must be 0 or 1"))
			return
		}
		status = &statusVal
	}

	if helpers.Contains(claims.Permissions, "view_all_products") {
		products, total, err = h.service.GetAll(ctx, limit, offset, search, status)
	} else if helpers.Contains(claims.Permissions, "view_active_products") {
		products, total, err = h.service.GetByStatusActive(ctx, limit, offset, search)
	} else {
		c.JSON(http.StatusForbidden, models.ErrorResponse(http.StatusForbidden, "Permission denied", nil))
		return
	}

	if err != nil {
		log.Printf("Error getting products: %v", err)
		c.JSON(http.StatusInternalServerError, models.ErrorResponse(http.StatusInternalServerError, "Failed to get products", err.Error()))
		return
	}

	totalPages := int((total + int64(limit) - 1) / int64(limit))
	c.JSON(http.StatusOK, models.PaginatedResponse{
		Status:      http.StatusOK,
		Message:     "Successfully Get Products",
		Data:        products,
		Total:       total,
		CurrentPage: page,
		PerPage:     limit,
		TotalPages:  totalPages,
		Error:       false,
	})
}

func (h *productHandlerImpl) CreateProduct(c *gin.Context) {
	ctx := c.Request.Context()

	var input models.CreateProductInput
	if err := c.ShouldBind(&input); err != nil {
		validationErrors := helpers.ParseValidationErrors(err)
		c.JSON(http.StatusBadRequest, models.ErrorResponse(http.StatusBadRequest, "Validation failed", validationErrors))
		return
	}

	product := models.Product{
		Name:        input.Name,
		Description: input.Description,
		Price:       input.Price,
		Quantity:    input.Quantity,
		Status:      1,
	}

	file, err := c.FormFile("image")
	if err == nil {
		s3Uploader, err := helpers.NewS3Uploader(
			os.Getenv("AWS_S3_BUCKET"),
			os.Getenv("AWS_S3_FOLDER"),
		)
		if err != nil {
			c.JSON(http.StatusInternalServerError, models.ErrorResponse(http.StatusInternalServerError, "Failed to initialize S3 uploader", err.Error()))
			return
		}

		imageURL, err := s3Uploader.UploadWithFallback(ctx, file, h.uploadDir)
		if err != nil {
			c.JSON(http.StatusInternalServerError, models.ErrorResponse(http.StatusInternalServerError, "Failed to upload image", err.Error()))
			return
		}

		product.ImageURL = imageURL
	}

	if err := h.service.Create(ctx, &product); err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse(http.StatusInternalServerError, "Failed to create product", err.Error()))
		return
	}

	c.JSON(http.StatusCreated, models.SuccessResponse(http.StatusCreated, "Product created successfully", product))
}

func (h *productHandlerImpl) UpdateProductStatus(c *gin.Context) {
	ctx := c.Request.Context()

	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse(http.StatusBadRequest, "Invalid product ID", "Product ID must be a number"))
		return
	}

	product, err := h.service.GetByID(ctx, uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, models.ErrorResponse(http.StatusNotFound, "Product not found", err.Error()))
		return
	}

	product.Status = 1 - product.Status

	if err := h.service.Update(ctx, uint(id), product); err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse(http.StatusInternalServerError, "Failed to update product status", err.Error()))
		return
	}

	c.JSON(http.StatusOK, models.SuccessResponse(http.StatusOK, "Product status updated successfully", product))
}

func (h *productHandlerImpl) UpdateProduct(c *gin.Context) {
	ctx := c.Request.Context()

	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse(http.StatusBadRequest, "Invalid product ID", "Product ID must be a number"))
		return
	}

	var input models.UpdateProductInput
	if err := c.ShouldBind(&input); err != nil {
		validationErrors := helpers.ParseValidationErrors(err)
		c.JSON(http.StatusBadRequest, models.ErrorResponse(http.StatusBadRequest, "Validation failed", validationErrors))
		return
	}

	product, err := h.service.GetByID(ctx, uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, models.ErrorResponse(http.StatusNotFound, "Product not found", err.Error()))
		return
	}

	product.Name = input.Name
	product.Description = input.Description
	product.Price = input.Price
	product.Quantity = input.Quantity

	oldImageURL := product.ImageURL

	file, err := c.FormFile("image")
	if err == nil {
		s3Uploader, err := helpers.NewS3Uploader(
			os.Getenv("AWS_S3_BUCKET"),
			os.Getenv("AWS_S3_FOLDER"),
		)
		if err != nil {
			c.JSON(http.StatusInternalServerError, models.ErrorResponse(http.StatusInternalServerError, "Failed to initialize S3 uploader", err.Error()))
			return
		}

		newImageURL, err := s3Uploader.UploadWithFallback(ctx, file, h.uploadDir)
		if err != nil {
			c.JSON(http.StatusInternalServerError, models.ErrorResponse(http.StatusInternalServerError, "Failed to upload image", err.Error()))
			return
		}

		product.ImageURL = newImageURL

		if oldImageURL != "" {
			if helpers.IsS3URL(oldImageURL) {
				if err := s3Uploader.DeleteFileFromS3(ctx, oldImageURL); err != nil {
					log.Printf("Failed to delete old image from S3: %v", err)
				}
			} else {
				oldFileName := helpers.GetFileNameFromURL(oldImageURL)
				oldFilePath := filepath.Join(h.uploadDir, oldFileName)
				if err := os.Remove(oldFilePath); err != nil && !os.IsNotExist(err) {
					log.Printf("Failed to delete old local image file: %v", err)
				}
			}
		}
	}

	if err := h.service.Update(ctx, uint(id), product); err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse(http.StatusInternalServerError, "Failed to update product", err.Error()))
		return
	}

	c.JSON(http.StatusOK, models.SuccessResponse(http.StatusOK, "Product updated successfully", product))
}

func (h *productHandlerImpl) DeleteProduct(c *gin.Context) {
	ctx := c.Request.Context()

	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse(http.StatusBadRequest, "Invalid product ID", "Product ID must be a number"))
		return
	}

	err = h.service.Delete(ctx, uint(id))
	if err != nil {
		if err.Error() == "product already deleted" {
			c.JSON(http.StatusBadRequest, models.ErrorResponse(http.StatusBadRequest, "Product already deleted", nil))
			return
		}

		if err.Error() == "record not found" || err.Error() == "gorm: record not found" {
			c.JSON(http.StatusNotFound, models.ErrorResponse(http.StatusNotFound, "Product not found", nil))
			return
		}

		c.JSON(http.StatusInternalServerError, models.ErrorResponse(http.StatusInternalServerError, "Failed to delete product", err.Error()))
		return
	}

	c.JSON(http.StatusOK, models.SuccessResponse(http.StatusOK, "Product deleted successfully", nil))
}
