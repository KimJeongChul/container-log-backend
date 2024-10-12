package v1

import (
	"bufio"
	"container-log-backend/model"
	"context"
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
	log.Println("container_id: ", containerID)

	upgrader := websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}

	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Println("Error upgrading connection:", err)
		return
	}
	defer conn.Close()

	// Connection timeout
	conn.SetReadLimit(maxMessageSize)
	conn.SetReadDeadline(time.Now().Add(pongWait))
	conn.SetPongHandler(func(appData string) error {
		conn.SetReadDeadline(time.Now().Add(pongWait))
		return nil
	})

	// Ping period
	ticker := time.NewTicker(pingPeriod)
	defer ticker.Stop()

	_, err = api.DockerClient.ContainerInspect(context.Background(), containerID)
	if err != nil {
		log.Printf("Error inspecting from container: %v\n", err)
		return
	}
	now := time.Now()
	lastTimestamp := now.Unix()

	// Get Docker Logstream
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	logOptions := container.LogsOptions{ShowStdout: true, ShowStderr: true, Follow: true, Since: strconv.FormatInt(lastTimestamp, 10)}

	logStream, err := api.DockerClient.ContainerLogs(ctx, containerID, logOptions)
	if err != nil {
		log.Printf("Error fetching logs from container: %v\n", err)
		return
	}
	defer logStream.Close()

	scanner := bufio.NewScanner(logStream)
	for scanner.Scan() {
		select {
		case <-ticker.C:
			// Sned ping message
			conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				log.Printf("Ping failed: %v\n", err)
				return
			}
		default:
			// Send log message
			logLine := scanner.Text()
			resMsg := model.Msg{
				Msg: logLine,
			}
			conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := conn.WriteJSON(resMsg); err != nil {
				log.Println("Error sending message:", err)
				return
			}
		}
	}

	if err := scanner.Err(); err != nil {
		log.Printf("Scann failed: %v\n", err)
		time.Sleep(1 * time.Second)
		return
	}
}
