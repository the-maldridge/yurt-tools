package main

import (
	"encoding/json"
	"log"

	"github.com/the-maldridge/yurt-tools/internal/consul"
)

var (
	taskdata map[string]map[string]map[string]consul.TaskData
)

func main() {
	taskdata = make(map[string]map[string]map[string]consul.TaskData)

	cs, err := consul.New()
	if err != nil {
		log.Fatal(err)
	}
	tasks, err := cs.KnownTasks()
	if err != nil {
		log.Fatal(err)
	}

	for _, task := range tasks {
		d, err := cs.LoadAllForTask(task)
		if err != nil {
			log.Printf("Error loading data for %s: %v", task.Name, err)
			continue
		}
		if taskdata[task.Job] == nil {
			taskdata[task.Job] = make(map[string]map[string]consul.TaskData)
		}
		if taskdata[task.Job][task.Group] == nil {
			taskdata[task.Job][task.Group] = make(map[string]consul.TaskData)
		}
		log.Println(task.Job, task.Group, task.Name)
		taskdata[task.Job][task.Group][task.Name] = d
	}

	bytes, err := json.Marshal(taskdata)
	if err != nil {
		log.Println(err)
	}
	log.Println(string(bytes[:]))
}
