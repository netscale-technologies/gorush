package gorush

import (
	"sync"
	"encoding/json"
	"net/http"
	"io/ioutil"
	"bytes"
	"github.com/gin-gonic/gin"
)

// InitWorkers for initialize all workers.
func InitWorkers(workerNum int64, queueNum int64) {
	LogAccess.Debug("worker number is ", workerNum, ", queue number is ", queueNum)
	QueueNotification = make(chan PushNotification, queueNum)
	for i := int64(0); i < workerNum; i++ {
		go startWorker()
	}
}

// SendNotification is send message to iOS or Android
func SendNotification(msg PushNotification) {
	switch msg.Platform {
	case PlatformIos:
		PushToIOS(msg)
	case PlatformAndroid:
		PushToAndroid(msg)
	case PlatformWeb:
		PushToWeb(msg)
	}
}

func startWorker() {
	for {
		notification := <-QueueNotification
		SendNotification(notification)
	}
}

func waitAndPerformCallback(callbackUrl string, count int, wg *sync.WaitGroup, log *[]LogPushEntry) {
	if (wg != nil) {
		wg.Wait()
		reqBody, err := json.Marshal(gin.H{
			"success": "ok",
			"counts":  count,
			"logs":    log,
		})
		if err != nil {
			var msg = "Error converting logs to JSON."
			LogAccess.Debug(msg)
			return
		}
		resp, err2 := http.Post(callbackUrl, "application/json", bytes.NewBuffer(reqBody))
		if err2 != nil {
			var msg = "Error posting logs to callback URL."
			LogAccess.Debug(msg)
			return
		}
		defer resp.Body.Close()
		ioutil.ReadAll(resp.Body)
	}
}

// queueNotification add notification to queue list.
func queueNotification(req RequestPush) (int, []LogPushEntry) {
	var count int
	var doSync = PushConf.Core.Sync
	if (req.Sync != nil) {
		doSync = *req.Sync
	}
	var callbackUrl = PushConf.Core.CallbackUrl
	if (req.CallbackUrl != nil) {
		callbackUrl = *req.CallbackUrl
	}
	wg := sync.WaitGroup{}
	newNotification := []PushNotification{}
	for _, notification := range req.Notifications {
		switch notification.Platform {
		case PlatformIos:
			if !PushConf.Ios.Enabled && !notification.Voip {
				continue
			}
			if !PushConf.Ios.VoipEnabled && notification.Voip {
				continue
			}
		case PlatformAndroid:
			if !PushConf.Android.Enabled {
				continue
			}
		case PlatformWeb:
			if !PushConf.Web.Enabled {
				continue
			}
		}
		notification.sync = doSync || callbackUrl != ""
		newNotification = append(newNotification, notification)
	}

	log := make([]LogPushEntry, 0, count)
	for _, notification := range newNotification {
		if doSync || callbackUrl != "" {
			notification.wg = &wg
			notification.log = &log
			notification.AddWaitCount()
		}
		QueueNotification <- notification
		switch notification.Platform {
		case PlatformWeb:
			count += len(notification.Subscriptions)
		default:
			count += len(notification.Tokens)
		}
		// Count topic message
		if notification.To != "" {
			count++
		}
	}

	if doSync {
		wg.Wait()
	} else if callbackUrl != "" {
		go waitAndPerformCallback(callbackUrl, count, &wg, &log)
	}

	StatStorage.AddTotalCount(int64(count))

	return count, log
}
