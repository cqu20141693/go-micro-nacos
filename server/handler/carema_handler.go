package handler

import (
	"errors"
	"fmt"
	ccMicro "github.com/cqu20141693/sip-server"
	"github.com/cqu20141693/sip-server/common"
	"github.com/cqu20141693/sip-server/db"
	"github.com/cqu20141693/sip-server/domain"
	"github.com/gin-gonic/gin"
	"go-micro.dev/v4/logger"
	"gorm.io/gorm"
	"net/http"
	"strconv"
)

type cameraApi interface {
	Insert(c *gin.Context)
	GetByCameraId(c *gin.Context)
	GetByGKSN(c *gin.Context)
	UpdateByCameraId(c *gin.Context)
	UpdateByGKSN(c *gin.Context)
	DeleteById(c *gin.Context)
	DeleteByGKSN(c *gin.Context)
	DeleteByCameraId(c *gin.Context)
}

func (c2 *CameraService) InitRouteMapper(router *gin.Engine) {
	router.POST("api/camera/insert", c2.Insert)
	router.POST("api/camera/getByCameraId", c2.GetByCameraId)
	router.POST("api/camera/getByGKSN", c2.GetByGKSN)
	router.POST("api/camera/updateByCameraId", c2.UpdateByCameraId)
	router.POST("api/camera/updateByGKSN", c2.UpdateByGKSN)
	router.POST("api/camera/deleteById", c2.DeleteById)
	router.POST("api/camera/deleteByGKSN", c2.DeleteByGKSN)
	router.POST("api/camera/deleteByCameraId", c2.DeleteByCameraId)
}

type CameraService struct {
}

func (c2 *CameraService) Insert(c *gin.Context) {
	var cameraVo domain.CameraVO
	err := c.ShouldBindJSON(&cameraVo)
	if err != nil {
		logger.Info("binding failed")
		c.JSON(http.StatusOK, common.ResultUtils.Fail("0000", err.Error()))
		return
	}
	do := domain.CameraDO{
		GroupKey: cameraVo.GroupKey,
		Sn:       cameraVo.Sn,
		CameraId: cameraVo.CameraId,
	}
	results := db.MysqlDB.Create(&do)
	if results.Error != nil {
		logger.Info("insert failed")
		c.JSON(http.StatusOK, common.ResultUtils.Fail("0000", results.Error.Error()))
		return
	}
	c.JSON(http.StatusOK, common.ResultUtils.Success(do.Id))
}

func (c2 *CameraService) GetByCameraId(c *gin.Context) {
	defer ccMicro.PanicFunc()
	cameraId, done := GetQueryParamStr(c, "cameraId", "required,max=64")
	if done {
		return
	}
	var cameraDo domain.CameraDO
	if results := db.MysqlDB.Where("camera_id=?", cameraId).First(&cameraDo); results.Error != nil {
		logger.Info("select failed")
		if errors.Is(results.Error, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusOK, common.ResultUtils.Success(nil))
		} else {
			c.JSON(http.StatusOK, common.ResultUtils.Fail("0000", results.Error.Error()))
		}
		return
	}

	c.JSON(http.StatusOK, common.ResultUtils.Success(domain.ConvertCameraDTO(&cameraDo)))

}

func GetQueryParamStr(c *gin.Context, name, tag string) (string, bool) {
	param := c.Query(name)

	err := ccMicro.Validate.Var(param, tag)
	if err != nil {
		logger.Info(fmt.Sprintf("param %s validate failed", name))
		c.JSON(http.StatusOK, common.ResultUtils.Fail("0000", err.Error()))
		return "", true
	}
	return param, false
}
func GetQueryParamUint64(c *gin.Context, name, tag string) (uint64, bool) {
	param := c.Query(name)

	err := ccMicro.Validate.Var(param, tag)
	if err != nil {
		logger.Info(fmt.Sprintf("param %s validate failed", name))
		c.JSON(http.StatusOK, common.ResultUtils.Fail("0000", err.Error()))
		return 0, true
	}
	parseInt, err := strconv.ParseUint(param, 10, 64)
	if err != nil {
		return 0, true
	}
	return parseInt, false
}

func (c2 *CameraService) GetByGKSN(c *gin.Context) {

	defer ccMicro.PanicFunc()
	groupKey := c.Query("groupKey")
	sn := c.Query("sn")

	err := ccMicro.Validate.Var(groupKey, "required,max=16,min=16")
	if err != nil {
		logger.Info("param groupKey validate failed")
		c.JSON(http.StatusOK, common.ResultUtils.Fail("0000", err.Error()))
		return
	}
	err = ccMicro.Validate.Var(sn, "required,max=32")
	if err != nil {
		logger.Info("param sn validate failed")
		c.JSON(http.StatusOK, common.ResultUtils.Fail("0000", err.Error()))
		return
	}
	var cameraDo domain.CameraDO
	if results := db.MysqlDB.Where("group_key=? and sn=?", groupKey, sn).First(&cameraDo); results.Error != nil {
		logger.Info("select failed")
		if errors.Is(results.Error, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusOK, common.ResultUtils.Success(nil))
		} else {
			c.JSON(http.StatusOK, common.ResultUtils.Fail("0000", results.Error.Error()))
		}
		return
	}

	c.JSON(http.StatusOK, common.ResultUtils.Success(domain.ConvertCameraDTO(&cameraDo)))

}

func (c2 *CameraService) UpdateByCameraId(c *gin.Context) {
	var req domain.UpdateReq
	err := c.ShouldBindJSON(&req)
	if err != nil {
		logger.Info("binding failed")
		c.JSON(http.StatusOK, common.ResultUtils.Fail("0000", err.Error()))
		return
	}
	do := domain.ConvertCameraDO(&req)
	results := db.MysqlDB.Model(&domain.CameraDO{}).Where("camera_id=?", req.CameraId).Updates(do)
	if results.Error != nil {
		logger.Info("update failed")
		c.JSON(http.StatusOK, common.ResultUtils.Fail("0000", results.Error.Error()))
	}
	c.JSON(http.StatusOK, common.ResultUtils.Success(results.RowsAffected > 0))
}

func (c2 *CameraService) UpdateByGKSN(c *gin.Context) {
	var req domain.UpdateReq
	err := c.ShouldBindJSON(&req)
	if err != nil {
		logger.Info("binding failed")
		c.JSON(http.StatusOK, common.ResultUtils.Fail("0000", err.Error()))
		return
	}
	results := db.MysqlDB.Model(&domain.CameraDO{}).Where("group_key=? and sn=?", req.GroupKey, req.Sn).Update("camera_id", req.CameraId)
	if results.Error != nil {
		logger.Info("update failed")
		c.JSON(http.StatusOK, common.ResultUtils.Fail("0000", results.Error.Error()))
	}
	c.JSON(http.StatusOK, common.ResultUtils.Success(results.RowsAffected > 0))
}

func (c2 *CameraService) DeleteById(c *gin.Context) {
	id, over := GetQueryParamUint64(c, "Id", "required,min=1")
	if over {
		return
	}

	tx := db.MysqlDB.Delete(&domain.CameraDO{Id: id})
	if tx.Error != nil {
		c.JSON(http.StatusOK, common.ResultUtils.Fail("0000", tx.Error.Error()))
		return
	}
	c.JSON(http.StatusOK, common.ResultUtils.Success(true))
}

func (c2 *CameraService) DeleteByGKSN(c *gin.Context) {
	groupKey, over := GetQueryParamStr(c, "groupKey", "required,min=16,max=16")
	if over {
		return
	}
	sn, over := GetQueryParamStr(c, "sn", "required,max=32")
	if over {
		return
	}
	tx := db.MysqlDB.Where("group_key=? and sn=?", groupKey, sn).Delete(&domain.CameraDO{})
	if tx.Error != nil {
		c.JSON(http.StatusOK, common.ResultUtils.Fail("0000", tx.Error.Error()))
		return
	}
	c.JSON(http.StatusOK, common.ResultUtils.Success(true))
}

func (c2 *CameraService) DeleteByCameraId(c *gin.Context) {
	cameraId, over := GetQueryParamStr(c, "cameraId", "required,max=64")
	if over {
		return
	}
	tx := db.MysqlDB.Where("camera_id=?", cameraId).Delete(&domain.CameraDO{})
	if tx.Error != nil {
		c.JSON(http.StatusOK, common.ResultUtils.Fail("0000", tx.Error.Error()))
		return
	}
	c.JSON(http.StatusOK, common.ResultUtils.Success(true))
}
