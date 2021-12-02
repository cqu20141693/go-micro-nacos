package common

import "github.com/gin-gonic/gin"

type BaseRestController struct {
}

func (c *BaseRestController) ResponseSuccess(g *gin.Context) {
	g.JSON(200, gin.H{
		"code": 200,
		"data": "",
		"msg":  "操作成功",
	})
}

func (c *BaseRestController) ResponseData(g *gin.Context, data interface{}) {
	g.JSON(200, gin.H{
		"code": 200,
		"data": data,
		"msg":  "操作成功",
	})
}
func (c *BaseRestController) ResponseFailureForAuth(g *gin.Context, err interface{}) {
	g.JSON(401, gin.H{
		"code": 401,
		"data": "",
		"msg":  err,
	})
}
func (c *BaseRestController) ResponseFailureForParameter(g *gin.Context, err interface{}) {
	g.JSON(403, gin.H{
		"code": 403,
		"data": "",
		"msg":  err,
	})
}

func (c *BaseRestController) ResponseFailureForFuncErr(g *gin.Context, err interface{}) {
	g.JSON(500, gin.H{
		"code": 500,
		"data": "",
		"msg":  err,
	})
}

func (c *BaseRestController) ResponseFailure(g *gin.Context, httpCode, code int, err interface{}) {
	g.JSON(httpCode, gin.H{
		"code": code,
		"data": "",
		"msg":  err,
	})
}
