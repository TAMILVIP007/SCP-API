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
	Gin.POST("/unban", RemoveBan)
	Gin.POST("/gentoken", GetToken)
}

func IndexPage(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "SCP-FoundationHQ API Is Running",
		"version": Version,
		"success": true,
	})
}

func AddBan(c *gin.Context) {
	var req *BanRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}
	role := FetchTokenRole(req.Token)
	if role == nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"message": "Invalid token",
		})
		return
	}
	if !strings.Contains(role.Role, "rsr") && !strings.Contains(role.Role, "cs") && !strings.Contains(role.Role, "mtf") {
		c.JSON(http.StatusUnauthorized, gin.H{
			"message": "You do not have the required permissions to perform this action",
		})
		return
	}
	if req.UserId == "" || req.Reason == "" || req.From == "" || req.BanClass == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Missing required parameters",
		})
		return
	}
	if strings.Contains(role.Role, "mtf") {
		if ok := SendRequest(req); ok {
			c.JSON(http.StatusAccepted, gin.H{
				"message": "Your ban request has been sent to the admins",
			})
		}
		return
	}
	if GetBanReason(req.UserId) != "" {
		c.JSON(http.StatusConflict, gin.H{
			"message": "User is already banned",
			"status":  "reason updated",
			"success": true,
		})
		return
	}
	AddNewBan(&BannedInfo{
		BanRequest: *req,
		TimeStamp:  time.Now().Format("2006-01-02 15:04:05"),
		BanId:      GenBanId(req.UserId, req.Reason),
	})
	c.JSON(http.StatusOK, GeneralResponse{
		Message: fmt.Sprintf("User %s banned successfully", req.UserId),
		Success: true,
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
	if req.UserID == "" || req.Role == "" || req.Token == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Missing required parameters",
		})
		return
	}
	role := FetchTokenRole(req.Token)
	if role == nil || role.Role != "rsr" {
		c.JSON(http.StatusUnauthorized, gin.H{
			"message": "You do not have the required permissions to perform this action",
		})
		return
	}
	token := GenToken(req.UserID, req.Role)
	err := AddNewToken(&TokensInfo{
		UserId:     req.UserID,
		Time:       time.Now().Format("2006-01-02 15:04:05"),
		Token:      token,
		Role:       req.Role,
		AssignedBy: role.UserId,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": fmt.Sprintf("Error generating token for user %s: %s", req.UserID, err.Error()),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": fmt.Sprintf("Token for user %s generated successfully", req.UserID),
		"token":   token,
		"success": true,
	})
}

func RemoveBan(c *gin.Context) {
	var req *UnbanRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}
	role := FetchTokenRole(req.Token)
	if role == nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"message": "Invalid token",
		})
		return
	}
	if !strings.Contains(role.Role, "rsr") && !strings.Contains(role.Role, "cs") && !strings.Contains(role.Role, "mtf") {
		c.JSON(http.StatusUnauthorized, gin.H{
			"message": "You do not have the required permissions to perform this action",
		})
		return
	}
	if req.UserId == "" || req.From == "" || req.Reason == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Missing required parameters",
		})
		return
	}
	if strings.Contains(role.Role, "mtf") {
		if ok := SendRequest(req); ok {
			c.JSON(http.StatusAccepted, gin.H{
				"message": "Your unban request has been sent to the admins",
				"success": true,
			})
		}
		return
	}
	if GetBanReason(req.UserId) == "" {
		c.JSON(http.StatusConflict, gin.H{
			"message": fmt.Sprintf("User %s is not banned", req.UserId),
			"success": true,
		})
		return
	}
	RemoveBanById(req.UserId)
	c.JSON(http.StatusOK, gin.H{
		"message": fmt.Sprintf("User %s unbanned successfully", req.UserId),
		"success": true,
	})
}

