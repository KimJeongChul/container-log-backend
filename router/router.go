package router

import (
	v1 "container-log-backend/api/v1"

	"github.com/gin-contrib/cors"
	"github.com/gin-contrib/secure"
	"github.com/gin-gonic/gin"

	"github.com/docker/docker/client"
)

func New(dockerClient *client.Client) *gin.Engine {
	router := gin.Default()
	router.Use(secure.New(secure.Config{
		FrameDeny:          true,
		ContentTypeNosniff: true,
		BrowserXssFilter:   true,
		IENoOpen:           true,
		ReferrerPolicy:     "strict-origin-when-cross-origin",
	}))
	router.Use(gin.Logger())
	router.Use(gin.Recovery())

	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},                                       // 허용할 오리진 (모든 오리진 허용)
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE"},            // 허용할 HTTP 메서드
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"}, // 허용할 헤더
		ExposeHeaders:    []string{"Content-Length"},                          // 클라이언트가 접근할 수 있는 헤더
		AllowCredentials: true,                                                // 쿠키와 인증 정보를 포함할 수 있도록 허용
	}))

	gin.SetMode(gin.ReleaseMode)

	containerAPI := v1.ContaienrAPI{DockerClient: dockerClient}
	logSocketAPI := v1.LogSocketAPI{DockerClient: dockerClient}
	
	router.GET("/containers", containerAPI.ListDockerContainers)
	router.GET("/logs/:container_id", logSocketAPI.LogContainer)
	return router
}
