package provider

// Error API错误结构
type Error struct {
	Message string `json:"message"`
	Type    string `json:"type,omitempty"`
	Code    int    `json:"code,omitempty"`
}
