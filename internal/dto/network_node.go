package dto

type CreateNetworkNodeRequest struct {
	Name        string `json:"name" binding:"required"`
	Description string `json:"description"`
	ParentID    *uint  `json:"parent_id"`
}

type UpdateNetworkNodeRequest struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	ParentID    *uint  `json:"parent_id"`
}

type NetworkNodeResponse struct {
	ID          uint                  `json:"id"`
	Name        string                `json:"name"`
	Description string                `json:"description"`
	ParentID    *uint                 `json:"parent_id,omitempty"`
	Children    []NetworkNodeResponse `json:"children,omitempty"`
	Devices     []DeviceResponse      `json:"devices,omitempty"`
	CreatedAt   string                `json:"created_at,omitempty"`
	UpdatedAt   string                `json:"updated_at,omitempty"`
}
