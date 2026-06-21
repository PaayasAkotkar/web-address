package example

import (
	webaddress "app/webaddress/core"
	"context"
	"io"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
)

// PlayWebAddress showcases the usage of the webaddress package
// the server starts -> push the request fetch the result
func PlayWebAddress() {

	// daily stuff
	app := gin.New()
	app.GET("/", func(ctx *gin.Context) {
		if ctx.GetHeader("Authorization") != "Bearer test-token" {
			ctx.Status(http.StatusUnauthorized)
			return
		}
		ctx.Data(http.StatusOK, "application/json", []byte(`{"status":"success 😄"}`))
	})

	app.POST("/db", func(ctx *gin.Context) {
		defer ctx.Request.Body.Close()
		b, err := io.ReadAll(ctx.Request.Body)

		if err != nil {
			ctx.JSON(http.StatusBadRequest, []byte(`{"error":"not able to read"}`))
		}
		log.Println("[data recieved 🤗]", string(b))
		ctx.Data(http.StatusOK, "application/json", []byte(`{"status":"success 😉"}`))
	})

	go func() {
		if err := app.Run(":3333"); err != nil && err != http.ErrServerClosed {
			log.Fatalf("server error: %v", err)
		}
	}()

	// end

	time.Sleep(100 * time.Millisecond)

	// web-address
	url := "http://localhost:3333/"
	client := webaddress.New(url)

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	// this is one way to do use!!! there can be many 🤗
	var wg sync.WaitGroup
	wg.Add(1)

	client.Request().GoMonitor(ctx, func(res *webaddress.Result) {
		defer wg.Done()
		log.Println("success:", string(res.Result))
	}, func() {
		log.Println("[welcome welcome welcome]")

		client.Request().
			Add("key1", "GET", nil).
			SetHeader("Authorization", "Bearer test-token")
		client.Request().SetBase("http://localhost:3333/db").Add("key12", "POST", []byte(`{"name":"don"}`))
	})

	wg.Wait()
	// end
}
