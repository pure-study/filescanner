package main

import (
    "fmt"
    "log"
    "regexp"
    "time"

    "github.com/radovskyb/watcher"
)

var watchFolder = "tmp/results"
var _watcher = watcher.New()

func printWatchingList() {
    fmt.Println("Current files under watching:")
    for path := range _watcher.WatchedFiles() {
        fmt.Printf("%s\n", path)
    }
    fmt.Println("====================================")
}

func main() {
    _watcher.SetMaxEvents(1)
    _watcher.FilterOps(watcher.Write, watcher.Create)

    reg := regexp.MustCompile("\\.csv$")
    _watcher.AddFilterHook(watcher.RegexFilterHook(reg, false))

    go func()  {
        for {
            select {
            case event := <-_watcher.Event:
                log.Println(event.String())
            case err := <-_watcher.Error:
                log.Fatalln(err)
            case <-_watcher.Closed:
                return
            }
        }
    }()

    if err := _watcher.AddRecursive(watchFolder); err != nil {
        log.Fatalln(err)
    }

    printWatchingList()

    if err := _watcher.Start(time.Microsecond * 100); err != nil {
        log.Fatalln(err)
    }
}