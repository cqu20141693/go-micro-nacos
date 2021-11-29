package handler

import (
	"github.com/cqu20141693/sip-server/client"
	"github.com/cqu20141693/sip-server/common"
	"github.com/gin-gonic/gin"
)

type sipApi interface {
	health(c *gin.Context)
	GetCameraInfo(c *gin.Context)
	selfHealth(c *gin.Context)
}

type SipService struct {
	sipClient *client.SipClient
}

func (s *SipService) GetCameraInfo(c *gin.Context) {
	panic("implement me")
}

func NewSipService(sipClient *client.SipClient) *SipService {
	return &SipService{sipClient: sipClient}
}

func (s *SipService) InitRouteMapper(router *gin.Engine) {
	router.POST(common.HealthPath, s.health)
	router.POST(common.GetCameraInfoPath, s.GetCameraInfo)
	router.POST(common.SelfHealth, s.selfHealth)
}

func (s *SipService) health(c *gin.Context) {
	c.JSON(200, common.ResultUtils.Success(map[string]string{"status": "up"}))
}

func (s *SipService) selfHealth(c *gin.Context) {
	health := s.sipClient.Health()
	c.JSON(200, health)
}
