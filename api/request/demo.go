package request

type DemoCreate struct {
	Name int64 `json:"name" binding:"required"`
}
