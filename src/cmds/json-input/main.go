package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"worker"
)

var cmdsfile = flag.String("cmdsfile", "cmds.json", "json file with commands to execute")

func ParseCommands(filename string) []*worker.JobPayload {
	jsonBytes, err := ioutil.ReadFile(filename)
	if err != nil {
		panic(err)
	}
	cmds := make([]*worker.JobPayload, 0)
	err = json.Unmarshal(jsonBytes, &cmds)
	if err != nil {
		panic(err)
	}
	return cmds
}

func main() {
	flag.Parse()
	cmds := ParseCommands(*cmdsfile)

	scheduler := worker.NewScheduler(200)

	go func() {
		for res := range scheduler.ResultsChan {
			fmt.Printf("Results: \n%#v\n", res)
			if res.FailedJob != nil {
				out, _ := json.Marshal(res.FailedJob)
				fmt.Printf("Failed job:\n%v\n", string(out))
			}
		}
	}()

	go func() {
		scheduler.Start()
	}()

	for _, cmd := range cmds {
		scheduler.Push(cmd)
	}
	sleep := make(chan bool)
	sleep <- true
}
