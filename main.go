package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/go-redis/redis"
	"github.com/gorilla/mux"
)

// Todo going to a function to add a task
// todo going to add an endpoint to execute the task
// TODO going to
type Task struct {
	Command  string
	Argument string
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
	router.HandleFunc("/executeTask", executeTask).Methods("GET")
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
	client.RPush("tasks", string(json))
}

func executeTask() {

	pubsub := client.PSubscribe("tasks")
	defer pubsub.Close()
	ch := pubsub.Channel()

	for msg := range ch {
		fmt.Println(msg.Channel)
	}
	/*pubsub.Channel()
	taskStr := client.LPop("tasks").Val()
	var task Task
	json.Unmarshal([]byte(taskStr), &task)
	out, err := exec.Command(task.Command, task.Argument).Output()
	if err != nil {
		fmt.Println(err)
	}

	w.Write(out)
	*/
}

func main() {
	router := Router()
	http.ListenAndServe(":5050", router)
}
