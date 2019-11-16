package middleware

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"spm/api"
	"spm/service"
	"strconv"
)

func Auth() gin.HandlerFunc {
	return func(c *gin.Context) {
		// before request

		id, err1 := c.Cookie("id")
		token, err2 := c.Cookie("token")

		if err1 != nil || err2 != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, api.ErrorResponse("cookie error"))
			return
		}

		uid := validateToken(id, token)
		//log.Println(id, token, accountType)
		if uid == 0 { // abort the operation
			c.AbortWithStatusJSON(http.StatusUnauthorized, api.ErrorResponse("unauthorized"))
			return
		}

		if c.Keys == nil {
			c.Keys = make(map[string]interface{})
		}
		c.Keys["id"] = uid

		c.Next()

		// after request
	}
}

func validateToken(id string, token string) uint {

	val, err := service.Cache.Get(id).Result()

	uid, ok := strconv.Atoi(id)

	if err != nil || val != token || ok != nil {
		return 0
	}

	return uint(uid)
}
