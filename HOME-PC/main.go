package main

import (
	"io"
	"net/http"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
)

func port() string {
	port := os.Getenv("PORT")
	if len(port) == 0 {
		port = "8088"
	}
	return ":" + port
}
func port_llm() string {
	port := os.Getenv("LLM_PORT_APP")
	if len(port) == 0 {
		port = "9999"
	}
	return ":" + port
}

func main() {
	log.SetFormatter(&log.JSONFormatter{})
	log.SetOutput(os.Stdout)

	router := gin.Default()
	router.Any("/*path", func(c *gin.Context) {
		target := "http://localhost:" + port_llm()

		// Построение целевого URL
		path := strings.TrimPrefix(c.Param("path"), "/")
		targetURL := target + "/" + path

		// Создание нового HTTP-запроса
		req, err := http.NewRequest(c.Request.Method, targetURL, c.Request.Body)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create request"})
			return
		}

		// Копируем заголовки запроса
		for key, values := range c.Request.Header {
			for _, value := range values {
				req.Header.Add(key, value)
			}
		}

		// Выполняем запрос
		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to forward request"})
			return
		}
		defer resp.Body.Close()

		// Читаем тело ответа
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to read response"})
			return
		}

		// Передаем ответ клиенту
		for key, values := range resp.Header {
			for _, value := range values {
				c.Header(key, value)
			}
		}
		c.Status(resp.StatusCode)
		c.Writer.Write(body)
	})

	log.Info("Starting pc-proxy on port " + port())
	// Запуск сервера
	router.Run(port())
}
