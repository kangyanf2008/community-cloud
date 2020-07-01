package web

import "github.com/gin-gonic/gin"

func loginV1(c *gin.Context) {
	name := c.Param("name")
	pass := c.Param("pass")
	//输出json结果给调用方
	c.JSON(200, gin.H{
		"name": name,
		"pass": pass,
	})
}
