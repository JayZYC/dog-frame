package request

type DemoCreate struct {
	Name string `json:"name" binding:"required"`
}
