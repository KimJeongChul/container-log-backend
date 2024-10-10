package v1

import (
	"container-log-backend/model"

	"context"
	"net/http"

	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
	"github.com/gin-gonic/gin"
)

type ContaienrAPI struct {
	DockerClient *client.Client
}

func (api *ContaienrAPI) ListDockerContainers(c *gin.Context) {
	var w http.ResponseWriter = c.Writer

	// CORS
	setWriterAccessControlAllow(&w)

	containers, err := api.DockerClient.ContainerList(context.Background(), container.ListOptions{All: true})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	var containerInfo []model.ContainerInfo
	for _, container := range containers {
		containerInfo = append(containerInfo, model.ContainerInfo{
			ID:     container.ID,
			Name:   container.Names[0],
			Status: container.State,
		})
	}

	c.JSON(http.StatusOK, containerInfo)
}
