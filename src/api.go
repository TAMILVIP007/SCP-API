package src

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

var (
	Gin = gin.Default()
)

func init() {
	Gin.GET("/", IndexPage)
	Gin.POST("/ban", AddBan)
	Gin.POST("/token", GetToken)
}

func IndexPage(c *gin.Context) {
	c.JSON(http.StatusOK, IndexResponse{
		Message: "SCP-FoundationHQ API Is Running",
		Version: Version,
	})
}

func AddBan(c *gin.Context) {
	var req BanRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	date := time.Now().Format("2006-01-02 15:04:05")
	if !CheckBanToken(req.Token) {
		c.JSON(http.StatusForbidden, gin.H{
			"message": "Invalid token",
		})
		return
	}

	if GetBanReason(req.ID) != "" {
		c.JSON(http.StatusConflict, gin.H{
			"message": fmt.Sprintf("User %s is already banned with reason %s. Updating ban...", req.ID, GetBanReason(req.ID)),
		})
		return
	}

	err := AddNewBan(&BannedInfo{
		UserId:   req.ID,
		Reason:   req.Reason,
		BannedBy: req.BannedBy,
		Date:     date,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": fmt.Sprintf("Error banning user %s: %s", req.ID, err.Error()),
		})
		return
	}

	c.JSON(http.StatusOK, BanResponse{
		Message: fmt.Sprintf("User %s banned successfully", req.ID),
	})
}

func GetToken(c *gin.Context) {
	var req TokenRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	if req.ID == "" || req.Rights == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Missing required parameters",
		})
		return
	}

	token := GenToken(req.ID, req.Rights)
	err := AddNewToken(&TokensInfo{
		UserId: req.ID,
		Time:   time.Now().Format("2006-01-02 15:04:05"),
		Token:  token,
		Rights: strings.Split(req.Rights, "+"),
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": fmt.Sprintf("Error generating token for user %s: %s", req.ID, err.Error()),
		})
		return
	}

	c.JSON(http.StatusOK, TokenResponse{
		Message: fmt.Sprintf("Token for user %s generated successfully", req.ID),
		Token:   token,
	})
}
