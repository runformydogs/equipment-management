package dto

type TreeNode struct {
	ID          uint       `json:"id"`
	Name        string     `json:"name"`
	Description string     `json:"description,omitempty"`
	Type        string     `json:"type"`
	Children    []TreeNode `json:"children,omitempty"`
}
