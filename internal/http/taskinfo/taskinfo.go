package taskinfo

import (
	"encoding/json"
	"log"
	"net/http"
	"sync"

	"github.com/flosch/pongo2/v4"
	"github.com/go-chi/chi/v5"

	"github.com/the-maldridge/yurt-tools/internal/consul"
	"github.com/the-maldridge/yurt-tools/internal/nomad"
)

type TaskInfo struct {
	cc *consul.Consul

	// We cache the task data in memory so we don't have to hit
	// consul every time we load up information on a given task.
	data map[string]map[string]map[string]map[string]consul.TaskData

	dMutex sync.RWMutex

	tmpls *pongo2.TemplateSet
}

func New() (*TaskInfo, error) {
	c, err := consul.New()
	if err != nil {
		return nil, err
	}

	sbl, err := pongo2.NewSandboxedFilesystemLoader("theme/p2")
	if err != nil {
		return nil, err
	}

	x := &TaskInfo{
		cc:    c,
		data:  make(map[string]map[string]map[string]map[string]consul.TaskData),
		tmpls: pongo2.NewSet("html", sbl),
	}
	x.tmpls.Debug = true
	x.update()

	return x, nil
}

func (ti *TaskInfo) HTTPEntry() chi.Router {
	r := chi.NewRouter()

	r.Get("/data/all", ti.dumpAll)
	r.Get("/data/detail/{namespace}/{job}/{group}/{task}", ti.dumpTask)

	r.Get("/view/", ti.viewAll)

	return r
}

func (ti *TaskInfo) update() {
	known, err := ti.cc.KnownTasks()
	if err != nil {
		log.Printf("Error while updating cache: %s", err)
		return
	}

	for _, t := range known {
		d, err := ti.cc.LoadAllForTask(t)
		if err != nil {
			log.Printf("Error refreshing information for task: %s", err)
			continue
		}
		ti.fillPath(t, d)
	}
}

func (ti *TaskInfo) fillPath(task nomad.Task, data consul.TaskData) {
	ti.dMutex.Lock()
	defer ti.dMutex.Unlock()

	if ti.data[task.Namespace] == nil {
		ti.data[task.Namespace] = make(map[string]map[string]map[string]consul.TaskData)
	}
	if ti.data[task.Namespace][task.Job] == nil {
		ti.data[task.Namespace][task.Job] = make(map[string]map[string]consul.TaskData)
	}
	if ti.data[task.Namespace][task.Job][task.Group] == nil {
		ti.data[task.Namespace][task.Job][task.Group] = make(map[string]consul.TaskData)
	}

	ti.data[task.Namespace][task.Job][task.Group][task.Name] = data
}

func (ti *TaskInfo) dumpAll(w http.ResponseWriter, r *http.Request) {
	enc := json.NewEncoder(w)
	ti.dMutex.RLock()
	defer ti.dMutex.RUnlock()
	enc.Encode(ti.data)
}

func (ti *TaskInfo) dumpTask(w http.ResponseWriter, r *http.Request) {
	ti.dMutex.RLock()
	defer ti.dMutex.RUnlock()

	enc := json.NewEncoder(w)
	d, ok := ti.data[chi.URLParam(r, "namespace")][chi.URLParam(r, "job")][chi.URLParam(r, "group")][chi.URLParam(r, "task")]
	if !ok {
		w.WriteHeader(http.StatusNotFound)
		enc.Encode(map[string]string{"error": "task not found"})
		return
	}

	enc.Encode(d)
}

func (ti *TaskInfo) viewAll(w http.ResponseWriter, r *http.Request) {
	t, err := ti.tmpls.FromCache("taskinfo.p2")
	if err != nil {
		ti.templateErrorHandler(w, err)
		return
	}
	ti.dMutex.RLock()
	defer ti.dMutex.RUnlock()
	ctx := make(map[string]interface{})
	ctx["data"] = ti.data
	if err := t.ExecuteWriter(ctx, w); err != nil {
		ti.templateErrorHandler(w, err)
	}
}

func (ti *TaskInfo) templateErrorHandler(w http.ResponseWriter, err error) {
	enc := json.NewEncoder(w)
	w.WriteHeader(http.StatusInternalServerError)
	enc.Encode(map[string]string{"error": err.Error()})
}
