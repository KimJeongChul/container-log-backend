package main

import (
	"container-log-backend/config"
	"container-log-backend/router"
	"container-log-backend/runner"
	"log"
	"net/http"

	"github.com/docker/docker/client"
	"github.com/gorilla/websocket"
)

var (
	dockerClient *client.Client
	upgrader     = websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool {
			return true // CORS 설정, 모든 오리진 허용
		},
	}
)

type ContainerInfo struct {
	ID     string `json:"id"`
	Name   string `json:"name"`
	Status string `json:"status"`
}

func init() {
	var err error
	dockerClient, err = client.NewClientWithOpts(client.FromEnv, client.WithVersion("1.40"))
	if err != nil {
		log.Fatalf("Error creating Docker client: %v", err)
	}
}

func main() {
	// Get OS Environments
	serverConfiguration, errConfig := config.Get()
	if errConfig != nil {
		panic(errConfig)
	}
	log.Println(serverConfiguration)
	
	// Create Gin Router
	engine := router.New(dockerClient)

	// Start runner
	if errRun := runner.Run(engine, serverConfiguration); errRun != nil {
		panic(errRun)
	}
}
