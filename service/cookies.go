package service

import "github.com/gin-gonic/gin"

const Cookies10days = 864000

var domain = cfg["domain"]

func SetCookie(c *gin.Context, key, value string, expire int) {
	c.SetCookie(key, value, expire, "/", domain, false, true)
}

func SetCookies(c *gin.Context, fields map[string]string, expire int) {
	for key, val := range fields {
		c.SetCookie(key, val, expire, "/", domain, false, true)
	}
}
