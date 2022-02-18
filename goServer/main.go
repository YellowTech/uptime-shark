package main

import (
	//"crypto/rand"
	"fmt"
	randMath "math/rand"
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
)

type Status struct {
	// ID        int64  `json:"id"`
	// Color     string `json:"color"`
	Name    string `json:"name"`
	Service string `json:"service"`
	Status  bool   `json:"status"`
	mu      sync.Mutex
}

// type State struct {
// 	mu sync.Mutex
// 	v  map[(string, string)]bool
// }

const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

func RandStringBytes(n int) string {
	b := make([]byte, n)
	for i := range b {
		b[i] = letterBytes[randMath.Intn(len(letterBytes))]
	}
	return string(b)
}

func getStatus(c *gin.Context) {
	// var stat = make([]status, len(state.v))

	// var i = 0
	// for key, value := range state.v {
	// 	stat[i].Service = key
	// 	stat[i].Status = value
	// 	i++
	// }

	c.SecureJSON(http.StatusOK, status)
}

func cookieTest(c *gin.Context) {
	cookie, err := c.Cookie("auth")

	if err != nil {
		cookie = "NotSet"
		c.SetSameSite(http.SameSiteStrictMode)
		c.SetCookie("auth", "test", 3600, "/", "localhost", true, true)
	}

	fmt.Printf("Cookie value: %s \n", cookie)

	c.String(http.StatusOK, fmt.Sprint(status))
}

func redir(c *gin.Context) {
	c.Redirect(http.StatusPermanentRedirect, "/tea")
}

func getClient(c *gin.Context) {
	c.File("./client.html")
}

var status = make([]Status, 0)

func main() {
	randMath.Seed(time.Now().UnixNano())

	fmt.Println("Starting Server")

	status = append(status, Status{Name: "Google", Service: "http://google.com", Status: false})
	status = append(status, Status{Name: "Lel", Service: "http://lel.yellowtech.ch", Status: false})

	go periodic()

	gin.SetMode(gin.DebugMode)
	// gin.SetMode(gin.ReleaseMode)
	router := gin.Default()

	router.GET("/cookie", cookieTest)
	router.GET("/api/status", getStatus)
	// router.GET("/input", getInput)
	// router.GET("/secret", getSecret)
	// router.GET("/assign", getAssign)
	// router.GET("/tea", getTea)
	// router.GET("/", redir)
	// router.GET("/client", getClient)

	// use this for outside of Docker
	router.Run("localhost:8000")
	//router.Run("web-input:8080")
}

func periodic() {
	timeout := time.Duration(5 * time.Second)
	client := http.Client{
		Timeout: timeout,
	}

	for {
		for i := range status {
			item := &status[i]
			resp, err := client.Get(item.Service)
			item.mu.Lock()
			if err != nil {
				fmt.Println(item.Name, err.Error())
				item.Status = false
			} else {
				fmt.Println(item.Name, resp.StatusCode)
				item.Status = true
			}
			item.mu.Unlock()
		}

		time.Sleep(time.Second * 180)
	}
}
