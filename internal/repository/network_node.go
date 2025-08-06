package repository

import (
	"equipment-management/internal/models"
	"gorm.io/gorm"
)

type NetworkNodeRepository struct {
	db *gorm.DB
}

func NewNetworkNodeRepository(db *gorm.DB) *NetworkNodeRepository {
	return &NetworkNodeRepository{db: db}
}

func (r *NetworkNodeRepository) Create(node *models.NetworkNode) error {
	return r.db.Create(node).Error
}

func (r *NetworkNodeRepository) GetByID(id uint) (*models.NetworkNode, error) {
	var node models.NetworkNode
	if err := r.db.Preload("Devices").Preload("Children").First(&node, id).Error; err != nil {
		return nil, err
	}
	return &node, nil
}

func (r *NetworkNodeRepository) Update(id uint, updateData *models.NetworkNode) (*models.NetworkNode, error) {
	var node models.NetworkNode
	if err := r.db.First(&node, id).Error; err != nil {
		return nil, err
	}

	if err := r.db.Model(&node).Updates(updateData).Error; err != nil {
		return nil, err
	}

	return &node, nil
}

func (r *NetworkNodeRepository) Delete(id uint) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Model(&models.Device{}).Where("network_node_id = ?", id).Update("network_node_id", nil).Error; err != nil {
			return err
		}
		if err := tx.Model(&models.NetworkNode{}).Where("parent_id = ?", id).Update("parent_id", nil).Error; err != nil {
			return err
		}
		return tx.Delete(&models.NetworkNode{}, id).Error
	})
}

func (r *NetworkNodeRepository) GetAll() ([]models.NetworkNode, error) {
	var nodes []models.NetworkNode
	if err := r.db.Find(&nodes).Error; err != nil {
		return nil, err
	}
	return nodes, nil
}

func (r *NetworkNodeRepository) GetFullTree() ([]models.NetworkNode, error) {
	var nodes []models.NetworkNode
	err := r.db.
		Preload("Devices").
		Preload("Children").
		Where("parent_id IS NULL").
		Find(&nodes).Error

	if err != nil {
		return nil, err
	}

	for i := range nodes {
		if err := r.loadChildrenRecursive(&nodes[i]); err != nil {
			return nil, err
		}
	}

	return nodes, nil
}

func (r *NetworkNodeRepository) loadChildrenRecursive(node *models.NetworkNode) error {
	if err := r.db.
		Preload("Devices").
		Preload("Children").
		Where("parent_id = ?", node.ID).
		Find(&node.Children).Error; err != nil {
		return err
	}

	for i := range node.Children {
		if err := r.loadChildrenRecursive(&node.Children[i]); err != nil {
			return err
		}
	}

	return nil
}
