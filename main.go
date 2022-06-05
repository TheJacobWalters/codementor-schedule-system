package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os/exec"
	"strconv"

	"github.com/go-redis/redis"
	"github.com/gorilla/mux"
	"github.com/robfig/cron"
)

// Todo going to a function to add a task
// todo going to add an endpoint to execute the task
// TODO going to
type Task struct {
	Command  string
	Argument string
	Time     string
}

var client redis.Client

func createRedisClient() redis.Client {
	client := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})
	return *client
}

func index(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("welcome to my program"))
}

func Router() *mux.Router {
	client = createRedisClient()
	router := mux.NewRouter()
	router.HandleFunc("/", index).Methods("GET")
	router.HandleFunc("/addTask", addTask).Methods("GET")
	return router
}

func addTask(w http.ResponseWriter, r *http.Request) {
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
	client.Publish("tasks", string(json))
}

func executor(command string, argument string, timer string) {
	c := cron.New()
	c.AddFunc(timer, func() { out, _ := exec.Command(command, argument).Output(); fmt.Println(string(out)) })
	c.Start()
}

func executeTask() {
	fmt.Println("starting execute task")
	pubsub := client.PSubscribe("tasks")
	defer pubsub.Close()
	ch := pubsub.Channel()
	var task Task
	for msg := range ch {
		taskStr := msg.Payload
		json.Unmarshal([]byte(taskStr), &task)
		fmt.Printf("Running Command %s %s at time %s \n", task.Command, task.Argument, task.Time)
		if task.Time != "" {
			go executor(task.Command, task.Argument, task.Time)
		} else {
			out, err := exec.Command(task.Command, task.Argument).Output()
			if err != nil {
				fmt.Println(err)
			}
			fmt.Println("Results of Task:")
			fmt.Println(string(out))
		}
	}
}

func main() {
	go executeTask()
	router := Router()
	http.ListenAndServe(":5050", router)
}
