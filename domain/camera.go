package domain

import "time"

/*
CREATE TABLE "tb_camera" (
  "id" bigint NOT NULL AUTO_INCREMENT,
  "create_at" datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
  "update_at" datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  "group_key" varchar(16) CHARACTER SET utf8 COLLATE utf8_bin NOT NULL,
  "sn" varchar(32) CHARACTER SET utf8 COLLATE utf8_bin NOT NULL,
  "camera_id" varchar(64) CHARACTER SET utf8 COLLATE utf8_bin NOT NULL,
  "token" varchar(32) COLLATE utf8mb4_general_ci NOT NULL,
  PRIMARY KEY ("id"),
  UNIQUE KEY "tb_camera_UN" ("camera_id")
) ENGINE=InnoDB AUTO_INCREMENT=2 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci
*/
//struct

type CameraDO struct {
	Id        uint64    `gorm:"primarykey" json:"id"`
	CreatedAt time.Time `gorm:"column:create_at" json:"createdAt"`
	UpdatedAt time.Time `gorm:"column:update_at" json:"updatedAt"`
	GroupKey  string    `gorm:"column:group_key" json:"groupKey"`
	Sn        string
	CameraId  string `gorm:"column:camera_id" json:"cameraId"`
	Token     string
}

func (p CameraDO) TableName() string {
	return "tb_camera"
}

type CameraVO struct {
	GroupKey string `json:"groupKey" binding:"required,max=16,min=16"`
	Sn       string `binding:"required,max=32"`
	CameraId string `json:"cameraId" binding:"required,max=64"`
	Token    string `binding:"required"`
}

type CameraDTO struct {
	*CameraDO
}

type UpdateReq struct {
	GroupKey string `json:"groupKey" binding:"required,max=16,min=16"`
	Sn       string `binding:"required,max=32"`
	CameraId string `json:"cameraId" binding:"required,max=64"`
	Token    string `binding:"required"`
}

//func

func ConvertCameraDTO(do *CameraDO) *CameraDTO {
	return &CameraDTO{
		do,
	}
}
