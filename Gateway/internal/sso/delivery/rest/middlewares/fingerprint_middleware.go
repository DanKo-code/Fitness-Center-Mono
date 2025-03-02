package middlewares

import (
	"crypto/md5"
	"encoding/hex"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
)

func FingerprintMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {

		userAgent := c.GetHeader("User-Agent")
		acceptHeaders := c.GetHeader("Accept")

		clientIP := c.ClientIP()

		fingerprintSource := strings.Join([]string{userAgent, acceptHeaders, clientIP}, "")

		hash := md5.New()
		hash.Write([]byte(fingerprintSource))
		fingerprint := hex.EncodeToString(hash.Sum(nil))

		c.Set(os.Getenv("APP_FINGERPRINT_REQUEST_KEY"), fingerprint)

		c.Next()
	}
}
