package controller

import (
	"net/http"
	"strconv"

	"equipment-management/internal/dto"
	"equipment-management/internal/service"
	"github.com/gin-gonic/gin"
)

type DeviceController struct {
	service *service.DeviceService
}

func NewDeviceController(service *service.DeviceService) *DeviceController {
	return &DeviceController{service: service}
}

func (c *DeviceController) CreateDevice(ctx *gin.Context) {
	var req dto.CreateDeviceRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request data"})
		return
	}

	device, err := c.service.CreateDevice(&req)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create device"})
		return
	}

	response := c.service.ToDeviceResponse(device)
	ctx.JSON(http.StatusCreated, response)
}

func (c *DeviceController) GetDevice(ctx *gin.Context) {
	id, err := strconv.ParseUint(ctx.Param("id"), 10, 32)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid device ID"})
		return
	}

	device, err := c.service.GetDevice(uint(id))
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "Device not found"})
		return
	}

	response := c.service.ToDeviceResponse(device)
	ctx.JSON(http.StatusOK, response)
}

func (c *DeviceController) UpdateDevice(ctx *gin.Context) {
	id, err := strconv.ParseUint(ctx.Param("id"), 10, 32)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid device ID"})
		return
	}

	var req dto.UpdateDeviceRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request data"})
		return
	}

	device, err := c.service.UpdateDevice(uint(id), &req)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update device"})
		return
	}

	response := c.service.ToDeviceResponse(device)
	ctx.JSON(http.StatusOK, response)
}

func (c *DeviceController) DeleteDevice(ctx *gin.Context) {
	id, err := strconv.ParseUint(ctx.Param("id"), 10, 32)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid device ID"})
		return
	}

	if err := c.service.DeleteDevice(uint(id)); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete device"})
		return
	}

	ctx.Status(http.StatusNoContent)
}

func (c *DeviceController) GetAllDevices(ctx *gin.Context) {
	devices, err := c.service.GetAllDevices()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get devices"})
		return
	}

	response := make([]dto.DeviceResponse, len(devices))
	for i, device := range devices {
		response[i] = c.service.ToDeviceResponse(&device)
	}

	ctx.JSON(http.StatusOK, response)
}
