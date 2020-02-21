package main

import (
	"log"
	"net/http"
	"html/template"

	"github.com/the-maldridge/yurt-tools/internal/consul"
)

var (
	taskdata map[string]map[string]map[string]consul.TaskData

	cs *consul.Consul
)

func update() {
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
		taskdata[task.Job][task.Group][task.Name] = d
	}
}

func homePageHandler(w http.ResponseWriter, r *http.Request) {
	t := template.Must(template.ParseFiles("tpl/base.tpl", "tpl/home.tpl"))
	t.Execute(w, taskdata)
}

func main() {
	taskdata = make(map[string]map[string]map[string]consul.TaskData)

	var err error
	cs, err = consul.New()
	if err != nil {
		log.Fatal(err)
	}

	update()

	http.HandleFunc("/", homePageHandler)
	log.Fatal(http.ListenAndServe(":8080", nil))
}
