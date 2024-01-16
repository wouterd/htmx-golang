package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
	"wouterd/playground/htmx/task"

	"github.com/gorilla/mux"
)

type TaskInsert struct {
    Task string `json:"task"`
}

func main() {
    r := mux.NewRouter()

    tasks := task.Tasks{}

    r.Path("/tasks").Methods(http.MethodPost).HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
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
        tasks.Add(task)
        w.Header().Add("Content-Type", "text/html")
        w.WriteHeader(http.StatusOK)
        fmt.Fprint(w, "<li>", taskJson.Task, "</li>")
    })
 
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
