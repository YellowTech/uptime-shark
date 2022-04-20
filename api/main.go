package main

import (
	//"crypto/rand"
	"context"
	"errors"
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
	"github.com/google/uuid"
	// "entgo.io/ent"
	// "golang.org/x/crypto/bcrypt"
	"net/url"

	"github.com/gin-contrib/cors"

	_ "github.com/mattn/go-sqlite3"
)

// type State struct {
// 	mu sync.Mutex
// 	v  map[(string, string)]bool
// }

type MonitorEdit struct {
	Id *string `json:"id" xml:"user"  binding:"required"`
	Name string `json:"name" binding:"required"`
	Interval int64 `json:"interval" binding:"required"`
	Status *bool `json:"status" binding:"required"`
	Inverted *bool `json:"inverted" binding:"required"`
	Mode string `json:"mode" binding:"required"`
	Url string `json:"url" binding:"required"`
}

const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

func RandStringBytes(n int) string {
	b := make([]byte, n)
	for i := range b {
		b[i] = letterBytes[randMath.Intn(len(letterBytes))]
	}
	return string(b)
}

func postEdit(c *gin.Context) {
	var editMonitor MonitorEdit
	if err := c.ShouldBindJSON(&editMonitor); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// test the url
	var err error
	_, err = url.ParseRequestURI(editMonitor.Url)
	if (err == nil) {
		httpClient := http.Client{
			Timeout: time.Duration(5 * time.Second),
		}
		var resp *http.Response
		resp, err = httpClient.Get(editMonitor.Url)
		if (err == nil) {
			if (resp.StatusCode >= 300) {
				err = errors.New("URL: status code not OK")
			}
		}
	}

	// check if valid uuid
	var checkedId uuid.UUID
	if (err == nil && len(*editMonitor.Id) != 0) {
		checkedId, err = uuid.Parse(*editMonitor.Id)
		if (err == nil) {
			_, err = client.Monitor.Query().Where(Monitor.ID(checkedId)).First(ctx)
		}
	}
	var updated *ent.Monitor
	if (err == nil) {
		// url is valid and id is either empty or valid
		if (len(*editMonitor.Id)==0) {
			updated, err = client.Monitor.Create().
				SetName(editMonitor.Name).
				SetInterval(editMonitor.Interval).
				SetNextCheck(0).
				SetMode(editMonitor.Mode).
				SetStatus(*editMonitor.Status).
				SetInverted(*editMonitor.Inverted).
				SetLogs(make([]logentry.LogEntry, 0)).
				SetNrLogs(20).
				SetStatusMessage("Never checked").
				SetURL(editMonitor.Url).
				SetRetries(5).
				Save(ctx)
		} else {
			updated, err = client.Monitor.Query().Where(Monitor.ID(checkedId)).First(ctx)
			if(err == nil){
				updated, err = updated.Update().
					SetName(editMonitor.Name).
					SetInterval(editMonitor.Interval).
					SetNextCheck(0).
					SetMode(editMonitor.Mode).
					SetStatus(*editMonitor.Status).
					SetInverted(*editMonitor.Inverted).
					SetURL(editMonitor.Url).
					Save(ctx)
			}
		}
	}


	if (err != nil) {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	log.Println("new edit received: ", editMonitor)
	log.Println(updated)

	c.JSON(http.StatusOK, editMonitor)//gin.H{"status": "ok"})
}

func postRemove(c *gin.Context) {
	var editMonitor MonitorEdit
	if err := c.ShouldBindJSON(&editMonitor); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var err error = nil
	var checkedId uuid.UUID

	// check if valid uuid
	checkedId, err = uuid.Parse(*editMonitor.Id)

	// delete if exists
	if (err == nil) {
		err = client.Monitor.DeleteOneID(checkedId).Exec(ctx)
	}

	if (err != nil) {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	} else {
		log.Println("new remove received: ", editMonitor.Name)
		c.JSON(http.StatusOK, editMonitor.Id)
		// Todo secure
	}
}

func getStatus(c *gin.Context) {
	list, _ := client.Monitor.Query().All(ctx)
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
	u, err := CreateService("Google", ctx, client)
	log.Println(err, u)

	// status = append(status, Status{Name: "Google", Service: "http://google.com", Status: false})
	// status = append(status, Status{Name: "Lel", Service: "http://lel.yellowtech.ch", Status: false})

	go periodic(ctx, client)

	gin.SetMode(gin.DebugMode)
	// gin.SetMode(gin.ReleaseMode)
	router := gin.Default()

	router.Use(cors.New(cors.Config{
        AllowOrigins:     []string{"http://localhost:8080"},
        AllowMethods:     []string{"GET", "PATCH"},
        AllowHeaders:     []string{"Origin"},
        ExposeHeaders:    []string{"Content-Length"},
        AllowCredentials: true,
        AllowOriginFunc: func(origin string) bool {
            return origin == "https://github.com"
        },
        MaxAge: 12 * time.Hour,
    }))

	router.GET("/cookie", cookieTest)
	router.GET("/api/status", getStatus)
	router.POST("/api/edit", postEdit);
	router.POST("/api/remove", postRemove);

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
			SetNrLogs(20).
			SetStatusMessage("Never checked").
			SetURL("http://google.com").
			SetRetries(5).
			Save(ctx)

		u, err = client.Monitor.Create().
			SetName("Yellowtech").
			SetInterval(10).
			SetNextCheck(0).
			SetMode("http").
			SetStatus(true).
			SetInverted(false).
			SetLogs(make([]logentry.LogEntry, 0)).
			SetNrLogs(20).
			SetStatusMessage("Never checked").
			SetURL("https://yellowtech.ch").
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
				
				if resp != nil {
					itemUpdate.SetStatusMessage(fmt.Sprint(resp.StatusCode))
				} else {
					itemUpdate.SetStatusMessage("Error")
				}
				
				itemUpdate.SetLogs(newLogs)
				itemUpdate.SetNextCheck(time.Now().Unix() + item.Interval)
				itemUpdate.SaveX(ctx)
				// item.mu.Unlock()
			}
		}

		time.Sleep(time.Second * 5)
	}
}
