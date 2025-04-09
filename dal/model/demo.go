package model

type Demo struct {
	Model
	Name string `gorm:"column:name;NOT NULL"`
}

func (d *Demo) TableName() string {
	return "demo"
}
