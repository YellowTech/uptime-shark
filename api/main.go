package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"time"
	"uptime/api/authentication"
	"uptime/api/ent"
	Monitor "uptime/api/ent/monitor"
	"uptime/api/logentry"
	"uptime/api/notifications"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"net/url"

	_ "github.com/mattn/go-sqlite3"
)

type MonitorEdit struct {
	Id *string `json:"id" xml:"user"  binding:"required"`
	Name string `json:"name" binding:"required"`
	Interval int64 `json:"interval" binding:"required"`
	Enabled *bool `json:"enabled" binding:"required"`
	Inverted *bool `json:"inverted" binding:"required"`
	Mode string `json:"mode" binding:"required"`
	Url string `json:"url" binding:"required"`
}

// Monitor is the model entity for the Monitor schema.
type MonitorPrivate struct {
	ID uuid.UUID `json:"id"`
	Name string `json:"name"`
	Enabled bool `json:"enabled"`
	Status bool `json:"status"`
	StatusMessage string `json:"statusMessage"`
	Inverted bool `json:"inverted"`
	Logs []logentry.LogEntry `json:"logs"`
}

// receiving edits from user
func postEdit(c *gin.Context) {
	var editMonitor MonitorEdit
	if err := c.ShouldBindJSON(&editMonitor); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// test the url
	var err error

	// prepend https if there is no protocol mentioned
	if !strings.HasPrefix(editMonitor.Url, "https://") && !strings.HasPrefix(editMonitor.Url, "http://") {
		editMonitor.Url = "https://" + editMonitor.Url
	}
	_, err = url.ParseRequestURI(editMonitor.Url)

	// do not check reachability if inverted
	if (err == nil && !*editMonitor.Inverted) {
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
	} else if err != nil {
		// url not valid
		c.JSON(http.StatusBadRequest, gin.H{"error": "URI is not valid"})
		return
	}

	// check if valid uuid and if exists
	var checkedId uuid.UUID
	if (err == nil && len(*editMonitor.Id) != 0) {
		checkedId, err = uuid.Parse(*editMonitor.Id)
		if (err == nil) {
			_, err = client.Monitor.Query().Where(Monitor.ID(checkedId)).First(ctx)
		}
	}

	var updated *ent.Monitor
	if (err == nil) {
		// uuid was found in db
		// url is valid and id is either empty or valid
		if (len(*editMonitor.Id)==0) {
			updated, err = client.Monitor.Create().
				SetName(editMonitor.Name).
				SetInterval(editMonitor.Interval).
				SetNextCheck(0).
				SetMode(editMonitor.Mode).
				SetEnabled(*editMonitor.Enabled).
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
					SetEnabled(*editMonitor.Enabled).
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

	log.Println("New edit received: ", editMonitor)
	log.Println(updated)

	c.JSON(http.StatusOK, editMonitor)
}

// receiving deletion requests
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
		// TODO secure
	}
}

func getStatus(c *gin.Context) {
	list, _ := client.Monitor.Query().All(ctx)

	// if not authenticated, send only limited information
	if authentication.CheckAuthCookie(c) {
		c.JSON(http.StatusOK, list)
	} else {
		var listPrivate []MonitorPrivate
		for _, x := range list {
			listPrivate = append(listPrivate, MonitorPrivate{ID: x.ID, Name: x.Name, Enabled: x.Enabled, Status: x.Status,
				StatusMessage: x.StatusMessage, Inverted: x.Inverted, Logs: x.Logs})
		}
		
		c.JSON(http.StatusOK, listPrivate)
	}
}

var ctx = context.Background()
var client *ent.Client = nil
var notify []*notifications.NotificationService = nil

func main() {
	fmt.Println("Starting Server")

	client2, err := ent.Open("sqlite3", "./db/db.sqlite3?mode=memory&cache=shared&_fk=1")
	client = client2

    if err != nil {
        log.Fatalf("failed opening connection to sqlite: %v", err)
    }
    defer client.Close()

    // Run the auto migration tool.
    if err := client.Schema.Create(context.Background()); err != nil {
        log.Fatalf("failed creating schema resources: %v", err)
    }

	// create test entry
	CreateService("Google", ctx, client)

	// delete all notifications
	client.Notification.Delete().ExecX(ctx)

	// create notification service from environment variables
	{
		key := os.Getenv("TELEGRAMKEY")
		chatId := os.Getenv("TELEGRAMCHAT")
		if key != "" && chatId != "" {
			client.Notification.Create().
				SetName("Telegram").
				SetActive(true).
				SetSettings([]string{"telegram", key, chatId}).
				SaveX(ctx)
		}
	}

	// Initialize notification setup
	notify, err = notifications.SetupNotifications(client, ctx)

	if (err != nil) {
		log.Fatalf("Failed loading notification entries: %v", err)
    }

	log.Println("Notification services loaded: ", notify)

	// launch go routine for periodic checking
	go periodic(ctx, client)

	if os.Getenv("DEBUG") == "true" {
		gin.SetMode(gin.DebugMode)
	} else {
		gin.SetMode(gin.ReleaseMode)
	}

	router := gin.Default()

	// Recovery middleware recovers from any panics and writes a 500 if there was one.
	router.Use(gin.Recovery())

	router.GET("/api/status", getStatus)

	authentication.InitializeAuth(router)

	// Authorization group
	authorized := router.Group("/", authentication.AuthRequired)
	{
		authorized.POST("/api/edit", postEdit);
		authorized.POST("/api/remove", postRemove);
	}


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
			SetEnabled(true).
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
			SetEnabled(true).
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

// go routine that periodically checks all monitors
func periodic(ctx context.Context, client *ent.Client) {
	timeout := time.Duration(5 * time.Second)
	httpClient := http.Client{
		Timeout: timeout,
	}

	// do indefinitely
	for {
		monitors, err := client.Monitor.Query().Where(Monitor.Enabled(true)).All(ctx)
		if err != nil {
			log.Println("Error loading monitors: ", err)
			notifications.SendMessageToAll(notify, "Error: Loading Monitors failed: " + fmt.Sprint(err))
			time.Sleep(time.Second * 10)
			continue
		}
		
		for _, item := range monitors {
			
			// if needs update
			if item.NextCheck < time.Now().Unix() {
				log.Println("Updating ", item.Name)
				itemUpdate := item.Update()
				// make http request
				resp, err := httpClient.Get(item.URL)
				
				// the new log entry that is being created
				var newLogEntry logentry.LogEntry
				
				// evaluating the results of the request
				if err != nil || resp.StatusCode >= 300 {
					newLogEntry = logentry.LogEntry { Failed: true, Message: "ERROR", Time: time.Now().Unix()}
				} else {
					fmt.Println(item.Name, resp.StatusCode)
					newLogEntry = logentry.LogEntry{Failed: false,Message: "OK",Time: time.Now().Unix()}
				}

				// if inverted -> invert the failed status
				if item.Inverted {
					newLogEntry.Failed = !newLogEntry.Failed
				}

				// the log entry list of the monitor
				logList := item.Logs

				logList = append(logList, newLogEntry)

				// if the logList is too long, shorten to the correct amount
				if (len(logList) > item.NrLogs) {
					logList = logList[len(logList) - item.NrLogs :]
				}

				// check if the last entry of the monitor needs to be changed
				// if failed and positive status => alert
				// if not failed and negative status => alert
				if newLogEntry.Failed == item.Status {
					itemUpdate.SetStatus(!newLogEntry.Failed)
					log.Println("Status change for " + item.Name + " deteceted")

					// change messages depending on if inverted or not
					if item.Inverted {
						if newLogEntry.Failed {
							notifications.SendMessageToAll(notify, "🔴 " + item.Name + "(inverted) is reachable")
						} else {
							notifications.SendMessageToAll(notify, "🟢 " + item.Name + "(inverted) is unreachable")
						}
					} else {
						if newLogEntry.Failed {
							notifications.SendMessageToAll(notify, "🔴 " + item.Name + " went down")
						} else {
							notifications.SendMessageToAll(notify, "🟢 " + item.Name + " is up")
						}
					}
				}

				// set status message of monitor item
				if resp != nil {
					itemUpdate.SetStatusMessage(fmt.Sprint(resp.StatusCode))
				} else {
					itemUpdate.SetStatusMessage("Error")
				}
				
				itemUpdate.SetLogs(logList)
				itemUpdate.SetNextCheck(time.Now().Unix() + item.Interval)
				_, err = itemUpdate.Save(ctx)
				if err!= nil {
					log.Println("Error saving monitor update for " + item.Name + ":", err)
				}
			}
		}

		time.Sleep(time.Second * 5)
	}
}
