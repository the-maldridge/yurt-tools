package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"sort"
	"time"

	"github.com/hashicorp/go-version"
	"github.com/heroku/docker-registry-client/registry"

	"github.com/the-maldridge/yurt-tools/internal/nomad"
)

var (
	pageinfo pagedata
	nc       *nomad.Client
)

type pagedata struct {
	TaskList []task
	Updated  time.Time
}

type task struct {
	Name    string
	Image   string
	Url     string
	Version string
	Newer   []string
	NoData  bool
}

func getTagsForImage(hub *registry.Registry, repo string) ([]string, error) {
	tags, err := hub.Tags(repo)
	if err != nil {
		return nil, err
	}
	return tags, nil
}

func getNewerVersions(tl []nomad.Task) ([]task, error) {
	out := make([]task, len(tl))

	url := "https://registry-1.docker.io/"
	username := os.Getenv("UP2DATE_REGISTRY_USERNAME")
	password := os.Getenv("UP2DATE_REGISTRY_PASSWORD")
	hub, err := registry.New(url, username, password)
	if err != nil {
		return nil, err
	}

	for i, task := range tl {
		if task.Driver != "docker" {
			continue
		}

		repoStr := task.Docker.Owner + "/" + task.Docker.Image
		if task.Docker.Owner == "" {
			repoStr = "library/" + task.Docker.Image
		}

		out[i].Name = task.Job + "/" + task.Name
		out[i].Image = repoStr
		out[i].Version = task.Docker.Tag
		out[i].Url = task.URL

		have, err := version.NewVersion(task.Docker.Tag)
		if err != nil {
			log.Printf("Task %s has uncomparable version: %s", task.Name, err)
			out[i].NoData = true
			continue
		}

		tags, err := getTagsForImage(hub, repoStr)
		if err != nil {
			log.Println(err)
			out[i].NoData = true
			continue
		}

		versions := []*version.Version{}
		for i := range tags {
			v, err := version.NewVersion(tags[i])
			if err != nil {
				continue
			}
			versions = append(versions, v)
		}
		sort.Sort(sort.Reverse(version.Collection(versions)))

		for _, v := range versions {
			if err != nil {
				log.Println("Attempted to parse unparseable version", task.Name, err)
				continue
			}
			if have.LessThan(v) {
				out[i].Newer = append(out[i].Newer, v.Original())
			}
			if len(out[i].Newer) > 5 {
				break
			}
		}
	}
	return out, nil
}

func updateData() {
	tasklist, err := nc.ListTasks(nomad.QueryOpts{})
	if err != nil {
		fmt.Fprintf(os.Stderr, err.Error())
		os.Exit(1)
	}

	tl, err := getNewerVersions(tasklist)
	if err != nil {
		fmt.Fprintf(os.Stderr, err.Error())
		os.Exit(1)
	}

	pageinfo.TaskList = tl
	pageinfo.Updated = time.Now()
	fmt.Println("Update complete!")
}

func statusHandler(w http.ResponseWriter, r *http.Request) {
	t := template.Must(template.ParseFiles("status.tpl"))
	t.Execute(w, pageinfo)
}

func main() {
	var err error
	nc, err = nomad.New()
	if err != nil {
		log.Fatal(err)
	}

	go func() {
		updateData()
		for range time.Tick(time.Hour * 4) {
			updateData()
		}
	}()

	http.HandleFunc("/", statusHandler)
	log.Fatal(http.ListenAndServe(":8080", nil))
}
