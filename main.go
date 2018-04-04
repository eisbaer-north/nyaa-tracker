package main

import (
	"fmt"
	conf "github.com/eisbaer-north/nyaa-tracker/config"
	trac "github.com/eisbaer-north/nyaa-tracker/tracker"
	"time"

)

func printTracker(t trac.Tracker) {
	for {
		msg := <-t.Out
		fmt.Println(msg)
	}
}

func main () {
	//Set the default configuration path
	config_path_file := "./config.json"
	//Load the configuration 
	config := conf.LoadConfig(config_path_file)
	if config.Autostart {
		trackers := trac.LoadTracker(config.Path)
		for _,tracker := range trackers {
			go tracker.StartTracking()
			go printTracker(tracker)
		}
	}
	for {
		time.Sleep(10 * time.Second)
	}
}
