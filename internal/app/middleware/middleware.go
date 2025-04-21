package middleware

import (
	"com.say.more.server/internal/app/constant"
	"com.say.more.server/internal/app/controllers"
	"com.say.more.server/internal/app/errors"
	token2 "com.say.more.server/internal/pkg/token"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"strconv"
	"strings"
	"time"
)

func SetAuth(r *gin.Engine) {
	r.Use(func(c *gin.Context) {
		var deviceType, device string
		deviceType = c.Request.Header.Get("Device-Type")
		device = c.Request.Header.Get("Device")
		c.Set("device_type", deviceType)
		c.Set("device", device)

		// ignore check method
		if _, ok := constant.IgnoreAuthMethodMap[c.FullPath()]; !ok {
			token := strings.TrimSpace(c.GetHeader("Authorization"))
			if token == "" {
				// get token from cookie
				authorization, err := c.Cookie("Authorization")
				if err != nil && !strings.Contains(err.Error(), "named cookie not present") {
					controllers.Unauthorized(c, err.Error())
					c.Abort()
					return
				}
				token = strings.TrimSpace(authorization)
			}
			if token == "" {
				token = strings.TrimSpace(c.Query("Authorization"))
			}
			if token == "" {
				controllers.Unauthorized(c, errors.ECTokenEmpty)
				c.Abort()
				return
			}

			studentId, _, err := token2.VerifyToken(token, constant.DefaultDeviceType)
			if err != nil {
				controllers.Unauthorized(c, errors.ECTokenInvalid)
				c.Abort()
				return
			}

			c.Set("student_id", strconv.FormatInt(studentId, 10))
		}
		c.Next()
	})
}

func SetCors(r *gin.Engine) {
	r.Use(cors.New(cors.Config{
		AllowMethods: []string{"GET", "POST", "PUT", "HEAD", "DELETE", "PATCH"},
		AllowHeaders: []string{"Origin", "Content-Length", "Content-Type", "signature",
			"Authorization", "X-Forwarded-For", "X-Real-Ip", "Device-Type",
			"X-Appengine-Remote-Addr", "Access-Control-Allow-Origin"},
		ExposeHeaders:    []string{"x-pagination", "Content-Disposition"},
		AllowCredentials: true,
		AllowAllOrigins:  true,
		MaxAge:           12 * time.Hour,
	}))
}

func Cors() gin.HandlerFunc {
	return func(c *gin.Context) {
		//origin := c.Request.Header.Get("Origin")
		//if origin != "" {
		//	c.Header("Access-Control-Allow-Origin", origin)
		//	// // Primarily sets Access-Control-Allow-Origin
		//	c.Header("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE, UPDATE")
		//	c.Header("Access-Control-Allow-Headers", "Origin, Device, Content-Length, Content-Type, signature, Authorization, X-Forwarded-For, X-Real-Ip, Device-Type, sec-ch-ua, sec-ch-ua-mobile, sec-ch-ua-platform, User-Agent")
		//	c.Header("Access-Control-Expose-Headers", "x-pagination, Content-Disposition, Content-Length, Access-Control-Allow-Origin, Access-Control-Allow-Headers, Cache-Control, Content-Language, Content-Type")
		//	c.Header("Access-Control-Allow-Credentials", "true")
		//}
		//if c.Request.Method == "OPTIONS" {
		//	c.AbortWithStatus(http.StatusNoContent)
		//	return
		//}
		c.Next()
	}
}
