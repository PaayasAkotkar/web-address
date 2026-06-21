package example

import (
	webaddress "app/webaddress/core"
	"context"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
)

func PlayWebAddress() {
	app := gin.New()
	app.GET("/", func(ctx *gin.Context) {
		if ctx.GetHeader("Authorization") != "Bearer test-token" {
			ctx.Status(http.StatusUnauthorized)
			return
		}
		ctx.Data(http.StatusOK, "application/json", []byte(`{"status":"success"}`))
	})

	app.POST("/db", func(ctx *gin.Context) {
		log.Println("data recieved")
		ctx.Data(http.StatusOK, "application/json", []byte(`{"status":"success"}`))
	})

	go func() {
		if err := app.Run(":3333"); err != nil && err != http.ErrServerClosed {
			log.Fatalf("server error: %v", err)
		}
	}()

	time.Sleep(100 * time.Millisecond)

	url := "http://localhost:3333/"
	client := webaddress.New(url)

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	var wg sync.WaitGroup
	wg.Add(1)

	client.Request().GoMonitor(ctx, func(res *webaddress.Result) {
		defer wg.Done()
		if res.Error != nil {
			panic(res.Error)
		}
		if string(res.Result) != `{"status":"success"}` {
			panic(string(res.Result))
		}
		log.Println("success:", string(res.Result))
	}, func() {
		log.Println("yo")

		client.Request().
			Add("key1", "GET").
			SetHeader("Authorization", "Bearer test-token")
		client.Request().SetBase("http://localhost:3333/db").Add("key12", "POST")
	})

	wg.Wait()
}
