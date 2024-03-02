package main

import (
	"fmt"
	"justdoit/templates"
	"log"
	"net/http"
	"time"
)

type Task struct {
	Id        string
	Name      string
	Completed bool
}

// auto_increment for id
var auto_increment = 1

// not a real DB of course
var DB = map[string]Task{
	fmt.Sprint(auto_increment): {Id: fmt.Sprint(auto_increment), Name: "Placeholder task :D", Completed: false},
}

func PAGE(w http.ResponseWriter, r *http.Request) {
	unfinished_tasks := map[string]Task{}
	for k, v := range DB {
		if !v.Completed {
			unfinished_tasks[k] = v
		}
	}

	data := struct {
		Tasks map[string]Task
		Day   string
	}{unfinished_tasks, time.Now().Weekday().String()}

	err := templates.General.Execute(w, data)
	if err != nil {
		panic(err)
	}
}

func TASKS(w http.ResponseWriter, r *http.Request) {
	// create a task
	if r.Header.Get("HX-Request") != "true" {
		http.Error(w, "Not a HTMX request", http.StatusForbidden)
		return
	}

	if r.Method == "POST" {
		task := r.PostFormValue("task")

		auto_increment++
		id := fmt.Sprint(auto_increment)
		data := Task{
			// theres no "auto_increment++" expression in go :(
			Id:        fmt.Sprint(auto_increment),
			Name:      task,
			Completed: false,
		}
		DB[id] = data

		templates.General.ExecuteTemplate(w, "task-list", data)
	}
}

func EDITFORM(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	task, ok := DB[id]
	if !ok {
		http.Error(w, "Task not found", http.StatusNotFound)
		return
	} else if r.Header.Get("HX-Request") != "true" {
		http.Error(w, "Not a HTMX request", http.StatusForbidden)
		return
	}

	templates.FormTemplate.Execute(w, task)
}

func TASK(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	task, ok := DB[id]

	if !ok {
		http.Error(w, "Task not found", http.StatusNotFound)
		return
	} else if r.Header.Get("HX-Request") != "true" {
		http.Error(w, "Not a HTMX request", http.StatusForbidden)
		return
	}

	if r.Method == "GET" || r.Method == "" {
		templates.General.ExecuteTemplate(w, "task", task)

	} else if r.Method == "PUT" {
		el_name := r.Header.Get("HX-Trigger-Name")

		if el_name == "rename-task" {
			task.Name = r.PostFormValue("name")
			DB[id] = task

			templates.General.ExecuteTemplate(w, "task", task)
		} else {
			task.Completed = true
			DB[id] = task

			data := struct {
				Task Task
				Type string
			}{Task: task, Type: "success"}

			templates.General.ExecuteTemplate(w, "toast", data)
		}
	} else if r.Method == "DELETE" {
		delete(DB, id)

		data := struct {
			Task Task
			Type string
		}{Task: task, Type: "default"}

		templates.General.ExecuteTemplate(w, "toast", data)
	}
}

func main() {
	http.HandleFunc("/", PAGE)
	http.HandleFunc("/task", TASKS)
	http.HandleFunc("/task/{id}", TASK)
	http.HandleFunc("/task/{id}/edit", EDITFORM)

	http.Handle("/img/", http.StripPrefix("/img/", http.FileServer(http.Dir("./img"))))

	log.Fatal(http.ListenAndServe(":8080", nil))
}
