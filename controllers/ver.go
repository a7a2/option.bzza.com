package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	//"github.com/monoculum/formam"
)

type Ver struct { //{"commandCode":"2069","terminalVersion":"1.11.3.3"}
	CommandCode     string `form:"commandCode";json:"commandCode";binding:"required"`
	TerminalVersion string `form:"terminalVersion";json:"terminalVersion";binding:"required"`
}

func PostVer(c *gin.Context) { //{"commandCode":"2069","terminalVersion":"1.11.3.3"}
	//if c.Request.PostForm.Encode() == "%7B%22commandCode%22%3A%222069%22%2C%22terminalVersion%22%3A%221.11.3.3%22%7D=" {
	c.JSON(http.StatusOK, gin.H{"msgResult": "Success"})

	//	var ver models.Ver
	//	if err := c.ShouldBind(&ver); err == nil {
	//		if ver.CommandCode == "2069" && ver.TerminalVersion == "1.11.3.3" {
	//			c.JSON(http.StatusOK, gin.H{"msgResult": "Success"})
	//		} else {
	//			c.JSON(http.StatusUnauthorized, gin.H{"status": "unauthorized"})
	//		}
	//	} else {
	//		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	//	}

}
