package tracker

import (
	"os"
	"io"
	"log"
	"time"
	"strings"
	"net/http"
	"io/ioutil"
	"encoding/json"
	"github.com/mmcdole/gofeed"
	"github.com/rjeczalik/notify"
)

type Tracker struct {
	Name string
	Prefix string
	Path string
	Rss string
	Interval int8
	Active bool
	In chan string
	Out chan string
}

func (t Tracker) StartTracking() {
	os.MkdirAll(t.Path, 0700)
	t.Prefix = "\033[36m" + t.Name + ":\033[0m "
	for {
		if t.Active {
			t.Out <-t.Prefix + "Checking for new torrents"
			rssparser := gofeed.NewParser()
			feed, err := rssparser.ParseURL(t.Rss)
			if err != nil {
				log.Fatal(err)
			}
			files, err := ioutil.ReadDir(t.Path)
			if err != nil {
				log.Fatal(err)
			}
			for _, i := range feed.Items {
				url_slice := strings.Split(i.Link, "/")
				filename := url_slice[len(url_slice)-1]
				if stringInSlice(files, filename)== false {
					t.DownloadTorrent(i.Link)
				}
			}
		}
		time.Sleep(time.Duration(t.Interval) * time.Minute)
	}
}

func (t Tracker) DownloadTorrent(link string) {
	tokens := strings.Split(link, "/")
	filename := tokens[len(tokens)-1]

	output, err := os.Create(t.Path+filename)
	if err != nil {
		log.Fatal(err)
	}
	defer output.Close()

	//Get response body
	response, err := http.Get(link)
	if err != nil {
		log.Fatal(err)
	}
	defer response.Body.Close()
	_, err = io.Copy(output, response.Body)
	if err != nil {
		log.Fatal(err)
	}
	t.Out <- t.Prefix + " downloaded " + filename
}

func (t Tracker) fileWatcher() {
	c := make(chan notify.EventInfo, 1)
	if err  := notify.Watch(t.Path, c, notify.InModify); err != nil {
		log.Fatal(err)
	}
	defer notify.Stop(c)
	for {
		switch event := <-c; event.Event() {
			case notify.InModify:
				t.In <- "tracker-changed"
		}
	}
}

func CreateTracker(path string) Tracker {
	file, err := ioutil.ReadFile(path)
	if err != nil {
		log.Fatal(err)
	}
	var tracker Tracker
	err = json.Unmarshal(file, &tracker)
	if err != nil {
		log.Fatal(err)
	}
	tracker.Out = make(chan string)
	tracker.In = make(chan string)
	return tracker
}

func LoadTracker(path string) []Tracker {
	var trackers []Tracker
	files, err := ioutil.ReadDir(path)
	if err != nil {
		log.Fatal(err)
	}
	for _, file := range files {
		filepath := path + "/" + file.Name()
		trackers = append(trackers, CreateTracker(filepath))
	}
	return trackers
}

func stringInSlice(list []os.FileInfo, a string) bool {
	for _, b := range list {
		if b.Name() == a {
			return true
		}
	}
	return false
}
