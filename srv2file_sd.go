package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"net"
	"strconv"
	"encoding/json"
	"io/ioutil"
	"time"
)

func main() {
	hostname := flag.String("srv", "", "Hostname that points to a srv record")
	outfile := flag.String("out", "", "Path to JSON file to write")
	loop := flag.Bool("loop", false, "Loop forever")
	looptime := flag.Int("time", 300, "Time to wait between hostname resolution refresh cycles")
	flag.Parse()

	if *hostname == "" || *outfile == "" {
		flag.PrintDefaults()
		os.Exit(1)
	}
	if *loop {
		for {
			generate(*hostname, *outfile)
			time.Sleep(time.Duration(*looptime) * time.Second)
		}
	} else {
		generate(*hostname, *outfile)
	}

}

func generate(hostname string, outfile string) {
	targets := []string{}
	log.Println(fmt.Sprintf("Looking up hostname %s", hostname))

	_, addrs, err := net.LookupSRV("", "", hostname)
	if err != nil {
		log.Println(fmt.Sprintf("Error: %s", err))
		return
	} else {
		for i := 0; i < len(addrs); i++ {
			srvaddrs, err := net.LookupHost(addrs[i].Target)
			log.Println(fmt.Sprintf("%d Target: %s", i, addrs[i].Target))
			log.Println(fmt.Sprintf("%d Port: %d", i, addrs[i].Port))
			if err == nil {
				for n := 0; n < len(srvaddrs); n++ {
					log.Println(fmt.Sprintf("%d IP: %s", i, srvaddrs[n]))
					targets = append(targets, srvaddrs[n]+":"+strconv.Itoa(int(addrs[i].Port)))
				}
			}
		}
	}
	file_sd := [1]map[string][]string{}
	file_sd[0] = map[string][]string{"targets": targets}
	file_sd_str, _ := json.Marshal(file_sd)

	file, err := ioutil.TempFile("", "file-sd")
	if err != nil {
		log.Println(err)
		return
	}
	if _, err = file.Write([]byte(file_sd_str)); err != nil {
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
