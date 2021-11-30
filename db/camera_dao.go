package db

import "time"

type CameraDO struct {
	Id        uint64    `gorm:"primarykey"`
	CreatedAt time.Time `gorm:"column:create_at"`
	UpdatedAt time.Time `gorm:"column:update_at"`
	GroupKey  string    `gorm:"column:group_key"`
	Sn        string
	CameraId  string `gorm:"column:camera_id"`
}

func (p CameraDO) TableName() string {
	return "tb_camera"
}
