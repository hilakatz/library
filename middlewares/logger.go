package middlewares

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"strings"
	"time"
)

func RequestLogger(c *gin.Context) {
	var params []string
	for _, val := range c.Request.URL.Query() {
		for _, v := range val {
			params = append(params, v)
		}
	}
	// Print the request details to the log
	fmt.Printf("%s - %s::%s::%v\n",
		time.Now().Format(time.RFC3339),
		c.Request.Method,
		c.Request.URL.Path,
		strings.Join(params, `,`),
	)
	// Continue with the request
	c.Next()
}
