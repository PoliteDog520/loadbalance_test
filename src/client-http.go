package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"io/ioutil"
	"log"
	"net/http"
	"time"
)

var counter int

const (
	// disable Keep-Alive to act as massive independent clients
	disableKeepalive = true

	timeout      = 5 * time.Second
	waitDuration = 1 * time.Second
)

func main() {
	initRoute()
}

func initRoute() {
	router := gin.New()

	// 根目錄回傳200
	router.GET("/", func(c *gin.Context) {
		client := &http.Client{
			Transport: &http.Transport{
				// disable Keep-Alive
				// @see https://www.cnblogs.com/cobbliu/p/4517598.html
				// @see https://nanxiao.me/en/a-brief-intro-of-tcp-keep-alive-in-gos-http-implementation/
				DisableKeepAlives: disableKeepalive,
			},
			Timeout: timeout,
		}
		for {
			connectAndShowResponse(client, "http://api-server-http-services:80/addr")
			time.Sleep(waitDuration)
		}

	})
	err := router.Run(":8080")
	if err != nil {
		log.Println("Router dead ☠️")
	}
}

func connectAndShowResponse(client *http.Client, url string) {
	resp, err := client.Get(url)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}
	defer resp.Body.Close()

	counter++
	body, _ := ioutil.ReadAll(resp.Body)
	fmt.Printf("#%d: %s\n", counter, string(body))
}
