package main

import (
	"github.com/Pauloo27/go-mpris"
	"github.com/godbus/dbus/v5"
	"os"
	"os/exec"
	"strings"
	"time"
)

var SLEEPTIME = time.Duration(60) * time.Second
var COMPLETION_CUTOFF = 0.85

type streamer struct {
	position, length float64
	title            string
}

func moveToParser(filename string) {
	directory, _ := strings.CutSuffix(os.Args[0], "scanner")
	parser := exec.Command("python3", directory+filename)
	parser.Stdout = os.Stdout
	err := parser.Run()
	if err != nil {
		panic(err)
	}
}
func saveToFile(in_streams *[]streamer, already_written *[]string) {
	string_to_save := ""
	for _, stream := range *in_streams {
		if stream.position/stream.length >= COMPLETION_CUTOFF {
			in_written := false
			for _, name := range *already_written {
				if name == stream.title {
					in_written = true
					break
				}
			}
			if in_written {
				continue
			}

			string_to_save += stream.title + "\n"
			*already_written = append(*already_written, stream.title)
		}
	}
	err := os.WriteFile("nearing_completion.txt", []byte(string_to_save), 0640)
	if err != nil {
		panic(err)
	}
}

func main() {
	already_written := make([]string, 4)
	bus, err := dbus.ConnectSessionBus()
	if err != nil {
		panic(err)
	}

	for {
		names, _ := mpris.List(bus)
		streams := make([]streamer, 0)
		for _, name := range names {
			player := mpris.New(bus, name)
			metadata, _ := player.GetMetadata()
			_, in := metadata["xesam:album"]
			if !in {
				pos, err1 := player.GetPosition()
				length, err2 := player.GetLength()
				title, err3 := metadata["xesam:title"].Value().(string)
				if err1 != nil || err2 != nil || !err3 {
					continue
				}
				streams = append(streams, streamer{pos, length, title})
			}

		}
		if len(streams) == 0 {
			time.Sleep(SLEEPTIME)
			continue
		}
		saveToFile(&streams, &already_written)
		moveToParser("parser.py")
	}

}
