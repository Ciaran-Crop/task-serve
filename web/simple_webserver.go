package web

import (
	"fmt"
	"io"
	"log"
	"net/http"
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
	fmt.Print("a")
}

func taskCreateHandle(w http.ResponseWriter, request *http.Request) {
	fmt.Print("a")
}

func taskCancelHandle(w http.ResponseWriter, request *http.Request) {
	fmt.Print("a")
}

func taskTaskListHandle(w http.ResponseWriter, request *http.Request) {
	fmt.Print("a")
}

func taskTaskViewHandle(w http.ResponseWriter, request *http.Request) {
	fmt.Print("a")
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
