package main

import (
	"io"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.Default()

	router.Any("/*path", func(c *gin.Context) {
		target := "http://localhost:9999"

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

	// Запуск сервера
	router.Run(":8088")
}
