package main

import (
	"github.com/BurntSushi/toml"
	"github.com/Pauloo27/go-mpris"
	"github.com/godbus/dbus/v5"
	"os"
	"os/exec"
	"strings"
	"time"
)

var SLEEPTIME = time.Duration(60) * time.Second
var COMPLETION_CUTOFF = 0.85
var conf Config

type streamer struct {
	position, length float64
	title            string
}
type Config struct {
	Whitelist bool
	Players   []string
	Parser    string
}

func moveToParser() {
	directory, _ := strings.CutSuffix(os.Args[0], "scanner")
	parser := exec.Command(directory + conf.Parser)
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
func removeByConf(names []string) []string {
	// I don't think using a pointer here really provides any benefit

	i := 0
	for i < len(names) {
		var name string
		nameArray := strings.SplitAfter(names[i], ".")
		if len(nameArray) > 4 && nameArray[len(nameArray)-1][0:8] == "instance" {
			name = strings.Join(nameArray[0:len(nameArray)-1], "")
		} else {
			name = names[i]
		}
		if conf.Whitelist {
			if !removeByConfIsIn(&name) {
				if i == 0 {
					names = names[1:]
				} else if i == len(names)-1 {
					names = names[0 : len(names)-1]
				} else {
					names = append(names[0:i], names[i+1:]...)
				}
			} else {
				i++
			}
		} else { // if conf is set to blacklist
			if removeByConfIsIn(&name) {
				if i == 0 {
					names = names[1:]
				} else if i == len(names)-1 {
					names = names[0 : len(names)-1]
				} else {
					names = append(names[0:i-1], names[i+1:]...)
				}
			} else {
				i++
			}
		}
	}
	return names
}
func removeByConfIsIn(InName *string) bool {
	// helper function for removeByConf
	// Does Go not have nested functions? I need to look at it again
	for _, name := range conf.Players {
		if *InName == "org.mpris.MediaPlayer2."+name {
			return true
		}
	}
	return false
}

func main() {
	_, err := toml.DecodeFile("./conf.conf", &conf)
	if err != nil {
		panic(err)
	}
	already_written := make([]string, 0)
	bus, err := dbus.ConnectSessionBus()
	if err != nil {
		panic(err)
	}

	for {
		names, _ := mpris.List(bus)
		names = removeByConf(names)
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
		moveToParser()
	}

}
