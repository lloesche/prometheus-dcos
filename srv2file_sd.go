package main

import (
	"encoding/json"
	"flag"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"
)

type tasks struct {
	Tasks []task `json:"tasks"`
}

type task struct {
	Host  string `json:"host"`
	Ports []int  `json:"ports"`
}

func extractHosts(blob []byte) (*[]string, error) {

	var t tasks
	err := json.Unmarshal(blob, &t)

	if err != nil {
		return nil, err
	}

	result := make([]string, len(t.Tasks))

	for index, task := range t.Tasks {
		result[index] = task.Host + ":" + strconv.Itoa(task.Ports[0])
	}

	return &result, nil
}

func getTasks(url string) ([]byte, error) {
	resp, err := http.Get(url)

	if err != nil {
		return nil, err
	}

	bodyBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return bodyBytes, nil
}

func taskToString(tasks *[]string) *[]byte {
	fileSd := [1]map[string][]string{}
	fileSd[0] = map[string][]string{"targets": *tasks}
	fileSdStr, _ := json.Marshal(fileSd)

	return &fileSdStr
}

func writeToFile(blob *[]byte, outfile string) {
	file, err := ioutil.TempFile("", "file-sd")
	if err != nil {
		log.Println(err)
		return
	}
	if _, err = file.Write([]byte(*blob)); err != nil {
		log.Println(err)
		return
	}

	if err = file.Close(); err != nil {
		log.Println(err)
		return
	}

	if err = os.Chmod(file.Name(), 0644); err != nil {
		log.Println(err)
		return
	}

	if err = os.Rename(file.Name(), outfile); err != nil {
		log.Println(err)
		return
	}
}

func createURL(leader string, application string) string {
	url := leader + "/service/marathon/v2/apps" + application + "/tasks"
	return url
}

func generate(url string, outfile string) {
	rest, _ := getTasks(url)
	tasks, _ := extractHosts(rest)
	blob := taskToString(tasks)
	writeToFile(blob, outfile)
}

func main() {
	srv := flag.String("srv", "", "service name e.g /nginx")
	outfile := flag.String("out", "", "Path to JSON file to write")
	loop := flag.Bool("loop", false, "Loop forever")
	looptime := flag.Int("time", 300, "Time to wait between hostname resolution refresh cycles in seconds")
	marathon := flag.String("marathon", "http://leader.mesos", "marathon host location")
	flag.Parse()

	if *srv == "" || *outfile == "" {
		flag.PrintDefaults()
		os.Exit(1)
	}
	url := createURL(*marathon, *srv)

	if *loop {
		for {
			generate(url, *outfile)
			time.Sleep(time.Duration(*looptime) * time.Second)
		}
	} else {
		generate(url, *outfile)
	}

}
