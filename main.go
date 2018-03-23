package main

import (
	"log"
	"os"
	"time"
	"strings"
	"io"
	"io/ioutil"
	"net/http"
	"encoding/json"
	"github.com/rjeczalik/notify"
	"github.com/mmcdole/gofeed"
)

const root = "/home/eisbaer/torrent-files/"

type Tracker struct{
	Name string
	Path string
	Rss string
	Interval int64
	Active bool
}

func create_tracker(path string) {
	//Load json file
	b, err := ioutil.ReadFile(path)
	if err != nil {
		log.Fatal(err)
	}
	//Unmarshal json content into the tracker struct
	var tracker Tracker
	err = json.Unmarshal(b, &tracker)
	if err != nil {
		log.Fatal(err)
	}
	go tracking(tracker)
}

func tracking(tracker Tracker) {
	//Create Directory in the torrent-files directory
	os.MkdirAll(root+tracker.Path, 0700)
	prefix := "\033[36m"+tracker.Name+":\033[0m"
	for {
		log.Println(prefix, "Checking for new torrents")
		fullpath := root + tracker.Path
		//get the rss feed ascosiated with the tracker
		rssparser := gofeed.NewParser()
		feed, err := rssparser.ParseURL(tracker.Rss)
		if err != nil {
			log.Fatal(err)
		}
		//get all existing torrent files
		files, err :=  ioutil.ReadDir(fullpath)
		if err != nil {
			log.Fatal(err)
		}
		//check if all items are already in directoy ascosiated with the tracker
		for _, i := range feed.Items {
			url_slice := strings.Split(i.Link, "/")
			filename := url_slice[len(url_slice)-1]
			if stringInSlice(filename, files) {

			} else {
				downloadTorrent(i.Link, fullpath, prefix)
			}
		}
		//Sleep for the minutes specified in the json file
		time.Sleep(time.Duration(tracker.Interval) * time.Minute)
	}
}

func stringInSlice(a string, list []os.FileInfo) bool {
	for _, b := range list {
		if b.Name() == a {
			return true
		}
	}
	return false
}
//Function downloading the torrent to the 
func downloadTorrent (url string, path string, prefix string) {
	//Split url to get Filename
	tokens := strings.Split(url, "/")
	filename := tokens[len(tokens)-1]

	//Create Download file
	output, err := os.Create(path+filename)
	if err != nil {
		log.Fatal(err)
	}
	defer output.Close()

	//Get response body
	response, err := http.Get(url)
	if err != nil {
		log.Fatal(err)
	}
	defer response.Body.Close()

	//Copy the contents from the download to the file
	_ ,err = io.Copy(output, response.Body)
	if err != nil {
		log.Fatal(err)
	}
	log.Println(prefix, "downloaded", filename)
}

func main() {
	//Creating Channel for notify checking the directory for new trackerfiles
	c := make(chan notify.EventInfo, 1)
	if err := notify.Watch("/home/eisbaer/nyaa-tracker", c, notify.InCreate, notify.Remove); err != nil {
		log.Fatal(err)
	}
	defer notify.Stop(c)

	// Block until an event is received
	for {
		switch event := <-c; event.Event() {
			//Creating a new tracker if for the new file inside the tracker folder
			case notify.InCreate:
				create_tracker(event.Path())
		}
	}
}
