package domain

import "time"

//struct

type CameraDO struct {
	Id        uint64    `gorm:"primarykey" json:"id"`
	CreatedAt time.Time `gorm:"column:create_at" json:"createdAt"`
	UpdatedAt time.Time `gorm:"column:update_at" json:"updatedAt"`
	GroupKey  string    `gorm:"column:group_key" json:"groupKey"`
	Sn        string
	CameraId  string `gorm:"column:camera_id" json:"cameraId"`
}

func (p CameraDO) TableName() string {
	return "tb_camera"
}

type CameraVO struct {
	GroupKey string `json:"groupKey" binding:"required,max=16,min=16"`
	Sn       string `json:"sn" binding:"required,max=32"`
	CameraId string `json:"cameraId" binding:"required,max=64"`
}

type CameraDTO struct {
	*CameraDO
}

type UpdateReq struct {
	GroupKey string `json:"groupKey" binding:"required,max=16,min=16"`
	Sn       string `json:"sn" binding:"required,max=32"`
	CameraId string `json:"cameraId" binding:"required,max=64"`
}

//func

func ConvertCameraDTO(do *CameraDO) *CameraDTO {
	return &CameraDTO{
		do,
	}
}

func ConvertCameraDO(req *UpdateReq) *CameraDO {
	return &CameraDO{
		GroupKey: req.GroupKey,
		Sn:       req.Sn,
		CameraId: req.CameraId,
	}
}
