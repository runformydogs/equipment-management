package dto

type CreateDeviceRequest struct {
	Type          string `json:"type" binding:"required"`
	Vendor        string `json:"vendor" binding:"required"`
	Model         string `json:"model" binding:"required"`
	Serial        string `json:"serial" binding:"required"`
	Location      string `json:"location" binding:"required"`
	NetworkNodeID *uint  `json:"network_node_id"`
}

type UpdateDeviceRequest struct {
	Type          string `json:"type"`
	Vendor        string `json:"vendor"`
	Model         string `json:"model"`
	Serial        string `json:"serial"`
	Location      string `json:"location"`
	Status        string `json:"status"`
	NetworkNodeID *uint  `json:"network_node_id"`
}

type DeviceResponse struct {
	ID            uint   `json:"id"`
	Type          string `json:"type"`
	Vendor        string `json:"vendor"`
	Model         string `json:"model"`
	Serial        string `json:"serial"`
	Location      string `json:"location"`
	Status        string `json:"status"`
	NetworkNodeID *uint  `json:"network_node_id,omitempty"`
	CreatedAt     string `json:"created_at"`
	UpdatedAt     string `json:"updated_at"`
}
