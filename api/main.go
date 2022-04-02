package main

import (
	//"crypto/rand"
	"context"
	"fmt"
	"log"
	randMath "math/rand"
	"net/http"
	// "sync"
	"time"
	"uptime/api/ent"
	Monitor "uptime/api/ent/monitor"
	"uptime/api/logentry"

	"github.com/gin-gonic/gin"
	// "entgo.io/ent"
	// "golang.org/x/crypto/bcrypt"


	_ "github.com/mattn/go-sqlite3"
)

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

	list, _ := client.Monitor.Query().All(ctx)
	list[0].Retries = 0
	c.JSON(http.StatusOK, list)
	// c.SecureJSON(http.StatusOK, list)
}

func cookieTest(c *gin.Context) {
	cookie, err := c.Cookie("auth")

	if err != nil {
		cookie = "NotSet"
		c.SetSameSite(http.SameSiteStrictMode)
		c.SetCookie("auth", "test", 3600, "/", "localhost", true, true)
	}

	fmt.Printf("Cookie value: %s \n", cookie)

	c.String(http.StatusOK, fmt.Sprint())
}

// func redir(c *gin.Context) {
// 	c.Redirect(http.StatusPermanentRedirect, "/tea")
// }

// func getClient(c *gin.Context) {
// 	c.File("./client.html")
// }

var ctx = context.Background()
var client *ent.Client = nil

func main() {

	randMath.Seed(time.Now().UnixNano())

	fmt.Println("Starting Server")

	client2, err := ent.Open("sqlite3", "./ent.sqlite3?mode=memory&cache=shared&_fk=1")
	client = client2

    if err != nil {
        log.Fatalf("failed opening connection to sqlite: %v", err)
    }
    defer client.Close()

    // Run the auto migration tool.
    if err := client.Schema.Create(context.Background()); err != nil {
        log.Fatalf("failed creating schema resources: %v", err)
    }

	// ctx := context.Background()

	// create test entry
	u, err := CreateService("google", ctx, client)
	log.Println(err, u)

	// status = append(status, Status{Name: "Google", Service: "http://google.com", Status: false})
	// status = append(status, Status{Name: "Lel", Service: "http://lel.yellowtech.ch", Status: false})

	go periodic(ctx, client)

	gin.SetMode(gin.DebugMode)
	// gin.SetMode(gin.ReleaseMode)
	router := gin.Default()

	router.GET("/cookie", cookieTest)
	router.GET("/api/status", getStatus)

	// use this for outside of Docker
	router.Run("localhost:8000")
	//router.Run("web-input:8080")
}

func CreateService(name string, ctx context.Context, client *ent.Client) (*ent.Monitor, error) {
	// create service if not existing
	u, err := client.Monitor.Query().Where(Monitor.Name(name)).First(ctx)

    if err != nil {
		u, err = client.Monitor.Create().
			SetName(name).
			SetInterval(10).
			SetNextCheck(0).
			SetMode("http").
			SetStatus(true).
			SetInverted(false).
			SetLogs(make([]logentry.LogEntry, 0)).
			SetNrLogs(4).
			SetStatusMessage("Never checked").
			SetURL("http://google.com").
			SetRetries(5).
			Save(ctx)
		log.Println("new monitor was created: ", u)
    }

    return u, err
}


func periodic(ctx context.Context, client *ent.Client) {
	timeout := time.Duration(5 * time.Second)
	httpClient := http.Client{
		Timeout: timeout,
	}

	for {
		monitors := client.Monitor.Query().AllX(ctx)
		
		for i := range monitors {
			item := monitors[i]
			
			// if needs update
			if item.NextCheck < time.Now().Unix() {
				fmt.Println("Updating ", item.Name)
				itemUpdate := item.Update()
				resp, err := httpClient.Get(item.URL)

				newLogs := item.Logs

				// item.mu.Lock()
				if err != nil || resp.StatusCode >= 300 {
					fmt.Println(item.Name, err.Error())
					newLogs = append(newLogs, 
						logentry.LogEntry{Failed: true,Message: "ERROR",Time: time.Now().Unix()},
					)
					itemUpdate.SetStatus(false)
				} else {
					fmt.Println(item.Name, resp.StatusCode)
					newLogs = append(newLogs, 
						logentry.LogEntry{Failed: false,Message: "OK",Time: time.Now().Unix()},
					)
					itemUpdate.SetStatus(true)
				}
				
				if (len(newLogs) > item.NrLogs) {
					newLogs = newLogs[len(newLogs) - item.NrLogs :]
				}
				
				itemUpdate.SetStatusMessage(fmt.Sprint(resp.StatusCode))
				itemUpdate.SetLogs(newLogs)
				itemUpdate.SetNextCheck(time.Now().Unix() + item.Interval)
				itemUpdate.SaveX(ctx)
				// item.mu.Unlock()
			}
		}

		time.Sleep(time.Second * 5)
	}
}
