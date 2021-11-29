package client

import (
	"context"
	"encoding/json"
	ccMicro "github.com/cqu20141693/sip-server"
	"github.com/cqu20141693/sip-server/common"
	"go-micro.dev/v4/client"
	"go-micro.dev/v4/logger"
)

const ServerName = "sip-link"

type SipClient struct {
	client client.Client
}

func NewSipClient() *SipClient {

	return &SipClient{ccMicro.CreateClient()}
}

func (receiver *SipClient) Health() *common.ResultCommon {
	req := receiver.client.NewRequest(ServerName, common.HealthPath, map[string]string{})
	var result map[string]interface{}
	err := receiver.client.Call(context.Background(), req, &result)
	if err != nil {
		logger.Info(err)
		return common.ResultUtils.Fail("0000", "call failed")
	}
	marshal, _ := json.Marshal(result)
	var resp common.ResultCommon
	_ = json.Unmarshal(marshal, &resp)
	return &resp
}
