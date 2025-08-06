package controller

import (
	"net/http"
	"strconv"

	"equipment-management/internal/dto"
	"equipment-management/internal/service"
	"github.com/gin-gonic/gin"
)

type NetworkNodeController struct {
	service *service.NetworkNodeService
}

func NewNetworkNodeController(service *service.NetworkNodeService) *NetworkNodeController {
	return &NetworkNodeController{service: service}
}

func (c *NetworkNodeController) CreateNode(ctx *gin.Context) {
	var req dto.CreateNetworkNodeRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request data"})
		return
	}

	node, err := c.service.CreateNode(&req)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create network node"})
		return
	}

	response := c.service.ToNetworkNodeResponse(node)
	ctx.JSON(http.StatusCreated, response)
}

func (c *NetworkNodeController) GetNode(ctx *gin.Context) {
	id, err := strconv.ParseUint(ctx.Param("id"), 10, 32)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid node ID"})
		return
	}

	node, err := c.service.GetNode(uint(id))
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "Network node not found"})
		return
	}

	response := c.service.ToNetworkNodeResponse(node)
	ctx.JSON(http.StatusOK, response)
}

func (c *NetworkNodeController) UpdateNode(ctx *gin.Context) {
	id, err := strconv.ParseUint(ctx.Param("id"), 10, 32)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid node ID"})
		return
	}

	var req dto.UpdateNetworkNodeRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request data"})
		return
	}

	node, err := c.service.UpdateNode(uint(id), &req)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update network node"})
		return
	}

	response := c.service.ToNetworkNodeResponse(node)
	ctx.JSON(http.StatusOK, response)
}

func (c *NetworkNodeController) DeleteNode(ctx *gin.Context) {
	id, err := strconv.ParseUint(ctx.Param("id"), 10, 32)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid node ID"})
		return
	}

	if err := c.service.DeleteNode(uint(id)); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete network node"})
		return
	}

	ctx.Status(http.StatusNoContent)
}

func (c *NetworkNodeController) GetAllNodes(ctx *gin.Context) {
	nodes, err := c.service.GetAllNodes()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get network nodes"})
		return
	}

	response := make([]dto.NetworkNodeResponse, len(nodes))
	for i, node := range nodes {
		response[i] = c.service.ToNetworkNodeResponse(&node)
	}

	ctx.JSON(http.StatusOK, response)
}

func (c *NetworkNodeController) GetFullTree(ctx *gin.Context) {
	tree, err := c.service.GetFullTree()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get tree"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"tree": tree})
}
