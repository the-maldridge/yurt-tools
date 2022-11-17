package taskinfo

import (
	"encoding/json"
	"log"
	"net/http"
	"sync"

	"github.com/flosch/pongo2/v4"
	"github.com/go-chi/chi/v5"
	"github.com/the-maldridge/yurt-tools/internal/kv"
	"github.com/the-maldridge/yurt-tools/internal/nomad"
)

type TaskInfo struct {
	cc *kv.KV

	// We cache the task data in memory so we don't have to hit
	// consul every time we load up information on a given task.
	data map[string]map[string]map[string]map[string]kv.TaskData

	dMutex sync.RWMutex

	tmpls *pongo2.TemplateSet
}

func New() (*TaskInfo, error) {
	c, err := kv.NewKVBackend()
	if err != nil {
		return nil, err
	}

	sbl, err := pongo2.NewSandboxedFilesystemLoader("theme/p2")
	if err != nil {
		return nil, err
	}

	x := &TaskInfo{
		cc:    c,
		data:  make(map[string]map[string]map[string]map[string]kv.TaskData),
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

	r.Get("/view", ti.viewAll)
	r.Get("/view/{namespace}/{job}/{group}/{task}", ti.viewAll)
	r.Get("/view/{namespace}/{job}/{group}/{task}/details", ti.viewDetails)

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

func (ti *TaskInfo) fillPath(task nomad.Task, data kv.TaskData) {
	ti.dMutex.Lock()
	defer ti.dMutex.Unlock()

	if ti.data[task.Namespace] == nil {
		ti.data[task.Namespace] = make(map[string]map[string]map[string]kv.TaskData)
	}
	if ti.data[task.Namespace][task.Job] == nil {
		ti.data[task.Namespace][task.Job] = make(map[string]map[string]kv.TaskData)
	}
	if ti.data[task.Namespace][task.Job][task.Group] == nil {
		ti.data[task.Namespace][task.Job][task.Group] = make(map[string]kv.TaskData)
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

func getTaskData(ti *TaskInfo, r *http.Request) (kv.TaskData, map[string]string, bool) {
	namespace := chi.URLParam(r, "namespace")
	job := chi.URLParam(r, "job")
	group := chi.URLParam(r, "group")
	task := chi.URLParam(r, "task")
	activeTask := make(map[string]string)
	d, ok := ti.data[namespace][job][group][task]
	if ok {
		activeTask["namespace"] = namespace
		activeTask["job"] = job
		activeTask["group"] = group
		activeTask["task"] = task
	}
	return d, activeTask, ok
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
	d, activeTask, hasActiveTask := getTaskData(ti, r)
	ctx["hasActiveTask"] = hasActiveTask
	ctx["activeTask"] = activeTask
	ctx["activeTaskData"] = d
	if err := t.ExecuteWriter(ctx, w); err != nil {
		ti.templateErrorHandler(w, err)
	}
}

func (ti *TaskInfo) viewDetails(w http.ResponseWriter, r *http.Request) {
	t, err := ti.tmpls.FromCache("taskinfo_details.p2")
	if err != nil {
		ti.templateErrorHandler(w, err)
		return
	}
	ti.dMutex.RLock()
	defer ti.dMutex.RUnlock()
	d, activeTask, _ := getTaskData(ti, r)
	jsonD, _ := json.MarshalIndent(d, "", "  ")
	log.Print(string(jsonD))
	ctx := make(map[string]interface{})
	ctx["activeTask"] = activeTask
	ctx["activeTaskData"] = d
	if err := t.ExecuteWriter(ctx, w); err != nil {
		ti.templateErrorHandler(w, err)
	}
}

func (ti *TaskInfo) templateErrorHandler(w http.ResponseWriter, err error) {
	enc := json.NewEncoder(w)
	w.WriteHeader(http.StatusInternalServerError)
	enc.Encode(map[string]string{"error": err.Error()})
}
