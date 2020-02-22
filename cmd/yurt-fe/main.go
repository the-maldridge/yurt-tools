package main

import (
	"log"
	"net/http"
	"time"

	"github.com/labstack/echo"

	"github.com/the-maldridge/yurt-tools/internal/consul"
	"github.com/the-maldridge/yurt-tools/internal/web"
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

func homePageView(c echo.Context) error {
	return c.Render(http.StatusOK, "home", taskdata)
}

func taskDetailView(c echo.Context) error {
	data := taskdata[c.Param("job")][c.Param("group")][c.Param("task")]
	return c.Render(http.StatusOK, "task-detail", data)
}

func updateTrigger(c echo.Context) error {
	update()
	return c.String(http.StatusOK, "Updated")
}

func main() {
	taskdata = make(map[string]map[string]map[string]consul.TaskData)

	var err error
	cs, err = consul.New()
	if err != nil {
		log.Fatal(err)
	}

	go func() {
		update()
		for range time.Tick(time.Hour) {
			update()
		}
	}()

	h, err := web.New("static", "tmpl")
	if err != nil {
		log.Fatal(err)
	}

	h.GET("/", homePageView)
	h.GET("/detail/:job/:group/:task", taskDetailView)
	h.GET("/update-now", updateTrigger)

	if err := h.Serve(); err != nil {
		log.Fatal(err)
	}
}
