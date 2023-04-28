package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"regexp"
	"strings"
	"time"

	"github.com/radovskyb/watcher"
	"github.com/robfig/cron/v3"
)

var (
	watchFolder = initEnvValue("WATCH_FOLDER", "tmp/orders")
	apiServer   = initEnvValue("API_SERVER", "http://localhost:8080")
	_watcher    = watcher.New()
)

func main() {
	_watcher.SetMaxEvents(1)
	_watcher.FilterOps(watcher.Write, watcher.Create)

	reg := regexp.MustCompile("\\.csv$")
	_watcher.AddFilterHook(watcher.RegexFilterHook(reg, false))

	go watching()

	if err := _watcher.AddRecursive(watchFolder); err != nil {
		log.Panicln(err)
	}

	printWatchingList()
	startCronJobToStopWatchingOldFiles()

	if err := _watcher.Start(time.Microsecond * 100); err != nil {
		log.Panicln(err)
	}
}

func watching() {
	for {
		select {
		case event := <-_watcher.Event:
			log.Println(event.String())
			log.Println(postFileChangeEvent(&event))
		case err := <-_watcher.Error:
			log.Panicln(err)
		case <-_watcher.Closed:
			return
		}
	}
}

func postFileChangeEvent(e *watcher.Event) (bool, error) {
	reqBody := strings.NewReader(fmt.Sprintf("{\"csvFile\":%q}", e.Path))
	resp, err := http.Post(apiServer+"/orders", "application/json", reqBody)
	if err != nil {
		log.Println(err)
		return false, err
	}

	defer resp.Body.Close()
	respBody, respErr := io.ReadAll(resp.Body)
	if respErr != nil {
		log.Println(respErr)
		return false, respErr
	}
	log.Printf("response body: %s\n", respBody)

	return true, nil
}

func printWatchingList() {
	files := _watcher.WatchedFiles()
	fmt.Printf("Currently %d file(s) under watching:\n", len(files))
	for path := range files {
		fmt.Printf("DEBUG: %s\n", path)
	}
	fmt.Println("====================================")
}

func startCronJobToStopWatchingOldFiles() {
	// trigger it for a first time
	go stopWatchingOldFiles()

	c := cron.New()
	c.AddFunc("@daily", stopWatchingOldFiles)
	c.Start()
	log.Printf("INFO: Cronjob started with entries: %v\n", c.Entries())
}

func stopWatchingOldFiles() {
	log.Println("INFO: Start trying to stop watching old files...")
	files, err := os.ReadDir(watchFolder)
	if err != nil {
		log.Printf("Error: %v\n", err)
		return
	}

	countRemovedFiles := 0
	now := time.Now()
	days3Ago := now.AddDate(0, 0, -3)
	log.Printf("DEBUG: 3 days ago: %s\n", days3Ago.Format(time.DateTime))

	for _, file := range files {
		fileInfo, infoErr := file.Info()
		if infoErr != nil {
			log.Printf("Warning: Ignoring the error: %v\n", infoErr)
			continue
		}

		if fileInfo.ModTime().Compare(days3Ago) < 0 {
			log.Printf("INFO: Stop watching old file: %q\n", file.Name())

			if file.IsDir() {
				_watcher.RemoveRecursive(file.Name())
			} else {
				_watcher.Remove(file.Name())
			}
			countRemovedFiles++
		}
	}
	log.Printf("INFO: %v file(s) stopped watching.\n", countRemovedFiles)
}

func initEnvValue(envName string, defaultValue string) (val string) {
	val, ok := os.LookupEnv(envName)
	if !ok {
		val = defaultValue
	}
	log.Printf("INFO: Using env[%s] value: %s\n", envName, val)
	return
}
