package main

import (
	"strings"
	"time"
	"log"
	conf "github.com/eisbaer-north/nyaa-tracker/config"
	trac "github.com/eisbaer-north/nyaa-tracker/tracker"
	gc "github.com/rthornton128/goncurses"
)

func printTracker(stdscr *gc.Window, t trac.Tracker, row int) {
	for {
		msg := <-t.Out
		stdscr.ColorOn(1)
		stdscr.MovePrint(4+row, 1, t.Name)
		stdscr.ColorOff(1)
		switch msg {
		case "keepalive":
			stdscr.MovePrint(4+row, 40, time.Now().Format("2006-01-02 15:04:05"))
		default:
			stdscr.MovePrint(4+row, 40, time.Now().Format("2006-04-15 15:04:05") + "\t" + msg)
		}
		stdscr.Refresh()
	}
}

func printTitle(stdscr *gc.Window, title string) {
	_, cols := stdscr.MaxYX()
	spacing := strings.Repeat(" ", ((cols - len(title)) / 2) - 1 )
	stdscr.Print(spacing)
	stdscr.Print(title)
}
func printRowOfChar(stdscr *gc.Window, char string, row int) {
	_, cols := stdscr.MaxYX()
	line := strings.Repeat(char, cols)
	stdscr.MovePrint(row, 0, line)
}
func printColHeadings(stdscr *gc.Window) {
	stdscr.MovePrint(2, 0, "Name")
	stdscr.MovePrint(2, 40, "last Update")
	stdscr.MovePrint(2, 64, "last episode")
}

func main () {
	//setup ncurses
	stdscr, err := gc.Init()
	if err != nil {
		log.Fatal("init", err)
	}
	defer gc.End()

	//enable multiple options for golang
	gc.Raw(true)
	gc.Echo(true)
	gc.Cursor(0)
	gc.UseDefaultColors()
	stdscr.Keypad(true)

	if err := gc.StartColor(); err != nil {
		log.Fatal(err)
	}

	gc.InitPair(1, gc.C_CYAN, gc.C_BLACK)


	//print the Title
	printTitle(stdscr, "NYAA-Tracker")
	printRowOfChar(stdscr, "=", 1)
	printColHeadings(stdscr)
	printRowOfChar(stdscr, "-", 3)
	stdscr.Refresh()

	//Set the default configuration path
	config_path_file := "./config.json"
	//Load the configuration 
	config := conf.LoadConfig(config_path_file)
	if config.Autostart {
		trackers := trac.LoadTracker(config.Path)
		for row ,tracker := range trackers {
			go tracker.StartTracking()
			go printTracker(stdscr, tracker, row)
			time.Sleep(time.Second)
		}
	}
	stdscr.GetChar()
	/*for {
		time.Sleep(10 * time.Second)
	}*/
}
