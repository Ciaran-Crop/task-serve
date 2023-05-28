package web

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"task-serve/api"
)

type HandleFnc func(http.ResponseWriter, *http.Request)

func logPanics(function HandleFnc) HandleFnc {
	return func(writer http.ResponseWriter, request *http.Request) {
		defer func() {
			if x := recover(); x != nil {
				log.Printf("[%v] caught panic : %v", request.RemoteAddr, x)

				http.Error(writer, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			}
		}()
		function(writer, request)
	}
}

func simpleHandle(w http.ResponseWriter, request *http.Request) {
	io.WriteString(w, "<h1>hello, world</h1>")
}

func RunSimpleServer() {
	http.HandleFunc("/test1", logPanics(simpleHandle))
	fmt.Println("Server start: localhost:9001")
	if err := http.ListenAndServe(":9001", nil); err != nil {
		panic(err)
	}
}

func taskPageHandle(w http.ResponseWriter, request *http.Request) {
	simpleHandle(w, request)
}

func taskCreateHandle(w http.ResponseWriter, request *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	switch request.Method {
	case "POST":
		request.ParseForm()
		responseJson := map[string]string{"result": "success"}
		taskName, taskCommand := request.Form.Get("task_name"), request.Form.Get("task_command")
		if taskId, err := api.CreateTask(taskName, taskCommand); err != nil {
			responseJson["result"] = "failed"
			rep, _ := json.Marshal(&responseJson)
			w.Write(rep)
		} else {
			responseJson["task_id"] = taskId
			rep, _ := json.Marshal(&responseJson)
			w.Write(rep)
		}
	}
}

func taskCancelHandle(w http.ResponseWriter, request *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	switch request.Method {
	case "GET":
		query := request.URL.Query()
		taskId := query.Get("task_id")
		responseJson := map[string]string{"result": "success"}
		if err := api.CancelTask(taskId); err != nil {
			responseJson["result"] = "failed"
			rep, _ := json.Marshal(&responseJson)
			w.Write(rep)
		} else {
			rep, _ := json.Marshal(&responseJson)
			w.Write(rep)
		}
	}
}

func taskTaskListHandle(w http.ResponseWriter, request *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	switch request.Method {
	case "GET":
		responseJson := make(map[string]interface{})
		if taskList, err := api.GetTasks(); err != nil {
			responseJson["result"] = "failed"
			rep, _ := json.Marshal(&responseJson)
			w.Write(rep)
		} else {
			responseJson["result"] = "success"
			responseTaskList := make([]map[string]interface{}, 0)
			for _, task := range taskList {
				responseTaskList = append(responseTaskList, task.GetMap())
			}
			responseJson["task_list"] = responseTaskList
			rep, _ := json.Marshal(&responseJson)
			w.Write(rep)
		}
	}
}

func taskTaskViewHandle(w http.ResponseWriter, request *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	switch request.Method {
	case "GET":
		query := request.URL.Query()
		taskId := query.Get("task_id")
		responseJson := make(map[string]interface{})
		if task, err := api.GetOneTask(taskId); err != nil {
			responseJson["result"] = "failed"
			rep, _ := json.Marshal(&responseJson)
			w.Write(rep)
		} else {
			responseJson["result"] = "success"
			responseJson["task_list"] = task.GetMap()
			rep, _ := json.Marshal(&responseJson)
			w.Write(rep)
		}
	}
}

func RunTaskServer() {
	http.HandleFunc("/", logPanics(taskPageHandle))
	http.HandleFunc("/create_task", logPanics(taskCreateHandle))
	http.HandleFunc("/cancel_task", logPanics(taskCancelHandle))
	http.HandleFunc("/get_task_list", logPanics(taskTaskListHandle))
	http.HandleFunc("/view_task", logPanics(taskTaskViewHandle))
	fmt.Println("Server start: localhost:9001")
	if err := http.ListenAndServe(":9001", nil); err != nil {
		panic(err)
	}
}
