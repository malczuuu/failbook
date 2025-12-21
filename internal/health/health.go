// Copyright (c) 2025 Damian Malczewski
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in all
// copies or substantial portions of the Software.
//
// SPDX-License-Identifier: MIT

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
