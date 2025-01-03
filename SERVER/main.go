package main

import (
	"io"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.Default()

	router.Any("/proxy/*path", func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if !validateToken(authHeader) {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
			return
		}

		target := "http://localhost:8081"

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
	router.Run(":8080") // Слушаем порт 8080
}

func validateToken(authHeader string) bool {
	const validToken = "8ebdb575c185d04674b93e23217b6589cad87e4d4c715520f3222164fa39469b_kitaici-pidorasini"
	if !strings.HasPrefix(authHeader, "Bearer ") {
		return false
	}
	token := strings.TrimPrefix(authHeader, "Bearer ")
	return token == validToken
}
