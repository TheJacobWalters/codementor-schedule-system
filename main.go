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
	task := Task{Command: query.Get("Command"), Argument: query.Get("Argument"), Time: query.Get("Time")}
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

// TODO 9 - Separate infinite loop processing into own function (L79 - 94)
func executeTask(c *redis.Client) {
	fmt.Println("starting execute task")

	pubsub := c.PSubscribe("tasks")
	defer pubsub.Close()

	ch := pubsub.Channel()

	// TODO 1 - create array of task statuses

	var task Task
	for msg := range ch {
		taskStr := msg.Payload
		json.Unmarshal([]byte(taskStr), &task)
		fmt.Printf("Running Command %s %s at time %s \n", task.Command, task.Argument, task.Time)
		if task.Time != "" {
			go executor(task.Command, task.Argument, task.Time)

			// TODO 2 - Update array with future execution status
		} else {
			out, err := exec.Command(task.Command, task.Argument).Output()
			if err != nil {
				fmt.Println(err)
			}
			fmt.Println("Results of Task:")
			fmt.Println(string(out))

			// TODO 3 - Update array with success status for task
		}
	}

	// TODO 4 - Return array of task statuses
}

func main() {
	go executeTask(&client)
	router := Router()
	http.ListenAndServe(":5050", router)
}
