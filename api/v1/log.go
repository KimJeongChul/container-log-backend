package v1

import (
	"bufio"
	"container-log-backend/model"
	"context"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)


type LogSocketAPI struct {
	DockerClient *client.Client
}

func (api *LogSocketAPI) LogContainer(c *gin.Context) {
	containerID := c.Param("container_id")

	upgrader := websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool {
			return true // CORS 설정, 모든 오리진 허용
		},
	}

	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Println("Error upgrading connection:", err)
		return
	}
	defer conn.Close()

	_, err = api.DockerClient.ContainerInspect(context.Background(), containerID)
	if err != nil {
		conn.WriteMessage(websocket.TextMessage, []byte(fmt.Sprintf("Container with ID '%s' not found.", containerID)))
		return
	}

	now := time.Now()
	lastTimestamp := now.Unix()
	log.Println("lastTimestamp:", now)
	for {
		reader, err := api.DockerClient.ContainerLogs(context.Background(), containerID, container.LogsOptions{ShowStdout: true, ShowStderr: true, Follow: true, Since: strconv.FormatInt(lastTimestamp, 10)})
		if err != nil {
			time.Sleep(2 * time.Second)
			continue
		}

		scanner := bufio.NewScanner(reader)
		for scanner.Scan() {
			line := scanner.Text()
			resMsg := model.ResMsg{
				Msg: line,
			}
			if err := conn.WriteJSON(resMsg); err != nil {
				log.Println("Error sending message:", err)
				break
			}
			fmt.Println(line)
			lastTimestamp = time.Now().Unix()
		}

		if err := scanner.Err(); err != nil {
			log.Fatalf("Error reading logs: %v", err)
			time.Sleep(2 * time.Second)
		} else {
			log.Println("Log stream closed normally, reconnecting...")
			time.Sleep(1 * time.Second)
		}
	}
}