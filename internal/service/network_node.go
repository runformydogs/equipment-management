package service

import (
	"equipment-management/internal/dto"
	"equipment-management/internal/models"
	"equipment-management/internal/repository"
	"fmt"
	"time"
)

type NetworkNodeService struct {
	repo *repository.NetworkNodeRepository
}

func NewNetworkNodeService(repo *repository.NetworkNodeRepository) *NetworkNodeService {
	return &NetworkNodeService{repo: repo}
}

func (s *NetworkNodeService) CreateNode(req *dto.CreateNetworkNodeRequest) (*models.NetworkNode, error) {
	node := models.NetworkNode{
		Name:        req.Name,
		Description: req.Description,
		ParentID:    req.ParentID,
	}
	if err := s.repo.Create(&node); err != nil {
		return nil, err
	}
	return &node, nil
}

func (s *NetworkNodeService) GetNode(id uint) (*models.NetworkNode, error) {
	return s.repo.GetByID(id)
}

func (s *NetworkNodeService) UpdateNode(id uint, req *dto.UpdateNetworkNodeRequest) (*models.NetworkNode, error) {
	updateData := models.NetworkNode{
		Name:        req.Name,
		Description: req.Description,
		ParentID:    req.ParentID,
	}
	return s.repo.Update(id, &updateData)
}

func (s *NetworkNodeService) DeleteNode(id uint) error {
	return s.repo.Delete(id)
}

func (s *NetworkNodeService) GetAllNodes() ([]models.NetworkNode, error) {
	return s.repo.GetAll()
}

func (s *NetworkNodeService) ToNetworkNodeResponse(node *models.NetworkNode) dto.NetworkNodeResponse {
	return dto.NetworkNodeResponse{
		ID:          node.ID,
		Name:        node.Name,
		Description: node.Description,
		ParentID:    node.ParentID,
		CreatedAt:   node.CreatedAt.Format(time.RFC3339),
		UpdatedAt:   node.UpdatedAt.Format(time.RFC3339),
	}
}

func (s *NetworkNodeService) GetFullTree() ([]dto.TreeNode, error) {
	nodes, err := s.repo.GetFullTree()
	if err != nil {
		return nil, err
	}

	return s.convertToTree(nodes), nil
}

func (s *NetworkNodeService) convertToTree(nodes []models.NetworkNode) []dto.TreeNode {
	result := make([]dto.TreeNode, 0, len(nodes))

	for _, node := range nodes {
		treeNode := dto.TreeNode{
			ID:          node.ID,
			Name:        node.Name,
			Description: node.Description,
			Type:        "node",
			Children:    make([]dto.TreeNode, 0),
		}

		for _, device := range node.Devices {
			treeNode.Children = append(treeNode.Children, dto.TreeNode{
				ID:   device.ID,
				Name: fmt.Sprintf("%s: %s", device.Type, device.Model),
				Type: "device",
			})
		}

		if len(node.Children) > 0 {
			childNodes := s.convertToTree(node.Children)
			treeNode.Children = append(treeNode.Children, childNodes...)
		}

		result = append(result, treeNode)
	}

	return result
}
