package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
)

type TestStruct struct {
	Message *string
}

func main() {
	r := gin.Default()

	r.POST("/post", func(c *gin.Context) {
		bodyBytes, err := io.ReadAll(c.Request.Body)
		if err != nil {
			c.JSON(400, gin.H{"error": "unable to read body"})
			return
		}
		defer c.Request.Body.Close()

		log.Println("Received body")
		fmt.Println(string(bodyBytes))
		log.Printf("Headers: %v", c.Request.Header)
		for key, values := range c.Request.Header {
			for _, value := range values {
				fmt.Printf("Header %s: %s\n", key, value)
			}
		}
		c.String(200, "got it\n")
	})

	type LoadTest struct {
		BlockCount int `json:"block_count"`
	}

	r.POST("/post-heavy", func(c *gin.Context) {
		bodyBytes, err := io.ReadAll(c.Request.Body)
		if err != nil {
			c.JSON(400, gin.H{"error": "unable to read body"})
			return
		}
		defer c.Request.Body.Close()
		var loadTest LoadTest
		err = json.Unmarshal(bodyBytes, &loadTest)
		if err != nil {
			c.JSON(400, gin.H{"error": "invalid JSON"})
			return
		}
		if loadTest.BlockCount <= 0 {
			c.JSON(400, gin.H{"error": "block_count must be greater than 0"})
			return
		}
		blockCount := loadTest.BlockCount
		const blockSize = 10 * 1024 * 1024 // 10 MB

		blocks := make([][]byte, 0, blockCount)
		for i := 0; i < blockCount; i++ {
			block := make([]byte, blockSize)
			for j := range block {
				block[j] = byte((i + j) % 256)
			}
			blocks = append(blocks, block)
		}

		sum := 0
		for _, block := range blocks {
			for _, b := range block {
				sum += int(b)
			}
		}
		fmt.Println("Simulated sum:", sum)

		c.String(200, "got it\n")
	})

	r.GET("/get", func(c *gin.Context) {
		query := c.Query("query")
		log.Printf("Received query: %s", query)
		log.Printf("Headers: %v", c.Request.Header)
		c.String(200, "got it\n")
	})

	var st TestStruct
	r.POST("/post-panic", func(c *gin.Context) {
		fmt.Println(*st.Message) // panic
	})

	// Create the server manually
	srv := &http.Server{
		Addr:    ":8080",
		Handler: r,
	}

	// Add shutdown route
	r.POST("/post-exit", func(c *gin.Context) {
		go func() {
			// Delay a bit so the response is sent before shutdown
			os.Exit(1)
		}()
		c.String(200, "server shutting down...\n")
	})

	r.POST("/post-long-time", func(c *gin.Context) {
		time.Sleep(30 * time.Second)
		c.String(200, "server Response after 30s\n")
	})

	r.POST("/post/delay/:seconds", func(c *gin.Context) {
		secondsStr := c.Param("seconds")
		seconds, err := time.ParseDuration(secondsStr + "s")
		if err != nil {
			c.JSON(400, gin.H{"error": "invalid seconds parameter"})
			return
		}
		time.Sleep(seconds)
		c.String(200, fmt.Sprintf("Response after %s\n", secondsStr))
	})

	// Handle CTRL+C too
	go func() {
		quit := make(chan os.Signal, 1)
		signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
		<-quit
		log.Println("Received termination signal, shutting down...")
		if err := srv.Shutdown(context.Background()); err != nil {
			log.Printf("Server forced to shutdown: %v", err)
		}
	}()

	log.Println("listening on :8080")
	if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatalf("listen: %s\n", err)
	}
}
