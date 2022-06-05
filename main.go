package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os/exec"
	"strconv"

	"github.com/go-redis/redis"
)

// Todo going to a function to add a task
// todo going to add an endpoint to execute the task
// TODO going to
type Task struct {
	Command  string
	Argument string
}

var Tasks []Task

func main() {

	client := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})

	pong, err := client.Ping().Result()
	fmt.Println(pong, err)

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("welcome to my program"))
	})

	http.HandleFunc("/addTask", func(w http.ResponseWriter, r *http.Request) {
		query := r.URL.Query()
		if !query.Has("Command") || !query.Has("Argument") {
			w.Write([]byte(strconv.FormatBool(query.Has("Command"))))
			w.Write([]byte(strconv.FormatBool(query.Has("Argument"))))
			return
		}
		task := Task{Command: query.Get("Command"), Argument: query.Get("Argument")}
		json, err := json.Marshal(task)
		if err != nil {
			fmt.Println(err)
			return
		}
		client.RPush("tasks", string(json))
	})

	http.HandleFunc("/executeTask", func(w http.ResponseWriter, r *http.Request) {
		if client.LLen("tasks").Val() == 0 {
			w.Write([]byte("There are no Tasks"))
			return
		}

		taskStr := client.LPop("tasks").Val()
		var task Task
		json.Unmarshal([]byte(taskStr), &task)
		out, err := exec.Command(task.Command, task.Argument).Output()
		if err != nil {
			fmt.Println(err)
		}

		w.Write(out)
	})

	// listen to port
	http.ListenAndServe(":5050", nil)
}
