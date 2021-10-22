package test

import (
	"context"
	"github.com/gin-gonic/gin"
	"github.com/msalvati1997/b1multicasting/handler"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func Test_1(t *testing.T) {
	router := gin.Default()
	routerGroup := router.Group(":90")
	routerGroup.GET("/groups", handler.GetGroups)
	routerGroup.POST("/groups", handler.CreateGroup)
	handler.GrpcPort = 90
	req, _ := http.NewRequestWithContext(context.Background(), "POST", "http://localhost:8080:/90/groups", strings.NewReader(`{"multicast_id": "m1","multicast_type": "BMULTICAST"}`))
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	t.Logf("status: %d", w.Code)
	t.Logf("response: %s", w.Body.String())
}
