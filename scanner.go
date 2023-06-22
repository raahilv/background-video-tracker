package main

import (
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"
)

var SLEEPTIME = time.Duration(60) * time.Second
var COMPLETION_CUTOFF = 0.85

func runPlayerctl(arguments []string) string {
	pc := exec.Command("playerctl", arguments...)
	output, _ := pc.Output()
	return string(output)
}
func moveToParser(filename string) {
	directory, _ := strings.CutSuffix(os.Args[0], "scanner")
	parser := exec.Command("python3", directory+filename)
	parser.Stdout = os.Stdout
	parser.Run()
}
func saveToFile(file_input *[][]string, already_written *[]string) {
	string_to_save := ""
	for _, details := range *file_input {
		current_time, _ := strconv.Atoi(details[0][1:])
		max_time, _ := strconv.Atoi(details[1])
		if float64(current_time)/float64(max_time) >= COMPLETION_CUTOFF {
			in_written := false
			for _, name := range *already_written {
				if name == details[2] {
					in_written = true
					break
				}
			}
			if in_written {
				continue
			}

			string_to_save += details[2] + "\n"
			*already_written = append(*already_written, details[2])
		}
	}
	os.WriteFile("nearing_completion.txt", []byte(string_to_save), 0640)
}

func main() {
	already_written := make([]string, 4)
	split_streams := [][]string{}
	for {
		if runPlayerctl([]string{"-l"}) == "" {
			time.Sleep(SLEEPTIME)
			continue
		}
		removalList := strings.Split(runPlayerctl([]string{"-a", "metadata", "-f", "{{album}}"}), "\n")
		list_streams := strings.Split(runPlayerctl([]string{"-a", "metadata", "-f", "'{{position}};{{mpris:length}};{{title}};{{mpris:length}}'"}), "\n")
		for i := 0; i < len(removalList); i++ {
			if removalList[i] != "" {
				list_streams = append(list_streams[0:i], list_streams[i+1:]...)
			}
		}

		for _, stream := range list_streams[:len(list_streams)-1] {
			split_streams = append(split_streams, strings.Split(stream, ";"))
		}
		i := 0
		for i < len(list_streams)-1 {
			if (split_streams[i][1] + "'") != split_streams[i][len(split_streams[i])-1] {
				split_streams = append(split_streams[0:i], split_streams[i+1:]...)
			} else {
				split_streams[i] = split_streams[i][:len(split_streams[i])-1]
				i++
			}
		}
		saveToFile(&split_streams, &already_written)
		moveToParser("parser.py")
	}

}
