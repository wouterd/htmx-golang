package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"
	"wouterd/playground/htmx/task"

	"github.com/gorilla/mux"
)

type TaskInsert struct {
    Task string `json:"task"`
}

type TaskPatch struct {
    Name *string `json:"name"`
    Completed *bool `json:"completed"`
}

func (patch TaskPatch) toStr() string {
    s := "{"
    if patch.Name != nil {
        s += " Name=" + *patch.Name
    }
    if patch.Completed != nil {
        s += " Completed=" + strconv.FormatBool(*patch.Completed)
    }
    return s + " }"
}

func taskCreateHandler(tasks *task.Tasks) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        var taskJson TaskInsert
        if r.Body == nil {
            w.WriteHeader(http.StatusBadRequest)
            return
        }
        
        err := json.NewDecoder(r.Body).Decode(&taskJson)
        if err != nil {
            fmt.Println(err.Error())
            w.WriteHeader(http.StatusBadRequest)
            return
        } 

        task := task.Task {
            Name: taskJson.Task,
            Created: time.Now().UTC(),
            Completed: false,
        }
        addedTask := tasks.Add(task)
        w.Header().Add("Content-Type", "text/html")
        w.WriteHeader(http.StatusOK)
        template := task_list_item(&addedTask)
        template.Render(context.Background(), w)
    }
}

func taskUpdateHandler(tasks *task.Tasks) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        fmt.Println("INCOMING PATCH!")
        vars := mux.Vars(r)
        id, err := strconv.Atoi(vars["id"])
        if err != nil {
            fmt.Printf("ID is not an int! (ID=%d)\n", id)
            w.WriteHeader(http.StatusInternalServerError)
            return
        }

        task, err := tasks.Get(id)
        if err != nil {
            w.WriteHeader(http.StatusNotFound)
            return
        }
        
        var taskJson TaskPatch
        if r.Body == nil {
            w.WriteHeader(http.StatusBadRequest)
            fmt.Fprintln(w, "No body found, can't patch")
            return
        }

        decodeErr := json.NewDecoder(r.Body).Decode(&taskJson)
        if decodeErr != nil {
            w.WriteHeader(http.StatusBadRequest)
            fmt.Fprintln(w, "No valid JSON body")
            return
        }

        fmt.Println("incoming task: ", taskJson.toStr())
        fmt.Println("current task: ", task)
        if taskJson.Name != nil {
            task.Name = *taskJson.Name
        }
        if taskJson.Completed != nil {
            task.Completed = *taskJson.Completed
        }
        fmt.Println("Merged task: ", task)
        tasks.Update(id, task)
    }
}

func taskUpdateFormHandler(tasks *task.Tasks) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        fmt.Println("INCOMING PATCH VIA FORM!")
        vars := mux.Vars(r)
        id, err := strconv.Atoi(vars["id"])
        if err != nil {
            fmt.Printf("ID is not an int! (ID=%d)\n", id)
            w.WriteHeader(http.StatusInternalServerError)
            return
        }

        task, err := tasks.Get(id)
        if err != nil {
            w.WriteHeader(http.StatusNotFound)
            return
        }

        task.Completed = r.PostFormValue("completed") != ""
        tasks.Update(id, task)

        template := tasks_list(tasks.Tasks())
        err = template.Render(context.Background(), w)
        if err != nil {
            fmt.Println(err.Error())
            w.WriteHeader(http.StatusInternalServerError)
        }
    }
}

func main() {
    fmt.Println("Big F")
    r := mux.NewRouter()

    tasks := task.Tasks{}

    r.Path("/tasks").Methods(http.MethodPost).HandlerFunc(taskCreateHandler(&tasks))
 
    r.Path("/tasks/{id:[0-9]+}").Methods(http.MethodPatch).Headers("Content-Type", "application/json").HandlerFunc(taskUpdateHandler(&tasks))

    r.Path("/tasks/{id:[0-9]+}").Methods(http.MethodPatch).Headers("Content-Type", "application/x-www-form-urlencoded").HandlerFunc(taskUpdateFormHandler(&tasks))

    r.Path("/").Methods(http.MethodGet).HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        template := index(tasks.Tasks())
        err := template.Render(context.Background(), w)
        if err != nil {
            fmt.Println(err.Error())
        }
    })

    fmt.Println("Starting HTTP server on port 8080")

    err := http.ListenAndServe(":8080", r)

    if err != nil {
        fmt.Println("Error starting server..:")
        fmt.Println(err.Error())
        return
    }

    quit := make(chan bool)

    <- quit
}
