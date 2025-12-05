package health

import (
	"net/http"
	"sync/atomic"

	"github.com/gin-gonic/gin"
)

type Status struct {
	ready int32
}

func NewStatus() *Status {
	return &Status{ready: 0}
}

func (s *Status) SetReady() {
	atomic.StoreInt32(&s.ready, 1)
}

func (s *Status) SetNotReady() {
	atomic.StoreInt32(&s.ready, 0)
}

func (s *Status) IsReady() bool {
	return atomic.LoadInt32(&s.ready) == 1
}

func LivenessHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status": "alive",
		})
	}
}

func ReadinessHandler(status *Status) gin.HandlerFunc {
	return func(c *gin.Context) {
		if status.IsReady() {
			c.JSON(http.StatusOK, gin.H{
				"status": "ready",
			})
		} else {
			c.JSON(http.StatusServiceUnavailable, gin.H{
				"status": "not ready",
			})
		}
	}
}
