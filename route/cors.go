package route

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// make cross-origin request
func Cors() gin.HandlerFunc {
	return func(c *gin.Context) {
		method := c.Request.Method
		origin := c.Request.Header.Get("Origin")
		if origin != "" {
			// Whether the response can be shared with requesting code from the given origin
			c.Writer.Header().Set("Access-Control-Allow-Origin", origin)
			// One or more HTTP request methods allowed when accessing a resource in response to preflight request
			c.Header("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE,UPDATE")
			// Used in response to a preflight request to indicate the HTTP headers that can be used during the actual request
			c.Header("Access-Control-Allow-Headers", "Content-Type, Content-Length, Authorization, Accept, X-Requested-With")
			// Allows a server to indicate which response headers should be made avaliable to scripts running in the browser in response to a cross-origin request
			c.Header("Access-Control-Expose-Headers", "Content-Length, Access-Control-Allow-Origin, Access-Control-Allow-Headers")
			//Indicates how long the results of a preflight request can be cached
			c.Header("Access-Control-Max-Age", "172800")
			// Tells browser whether the server allows credentials to be included in cross-origin HTTP requests
			c.Header("Access-Control-Allow-Credentials", "true")
		}
		if method == "OPTIONS" {
			c.AbortWithStatus(http.StatusOK)
		}
		c.Next()
	}
}
