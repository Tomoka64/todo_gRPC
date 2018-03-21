package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"os"

	todo "github.com/Tomoka64/todoWithGRPC/todo"

	"golang.org/x/net/context"
	grpc "google.golang.org/grpc"
)

type taskServer struct{}

const (
	dbPath = "../config/data.json"
)

func main() {
	fmt.Println("listening to port 8888")
	srv := grpc.NewServer()
	var tasks taskServer
	todo.RegisterTasksServer(srv, tasks)
	l, _ := net.Listen("tcp", ":8888")
	log.Fatal(srv.Serve(l))
}

func (taskServer) Add(ctx context.Context, text *todo.Text) (*todo.Task, error) {
	task := &todo.Task{
		Text: text.Text,
	}

	b, err := json.Marshal(task)
	if err != nil {
		return &todo.Task{}, fmt.Errorf("%s", err)
	}

	f, err := os.OpenFile(dbPath, os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		return &todo.Task{}, fmt.Errorf("%s", err)
	}
	_, err = f.Write(b)
	if err != nil {
		return &todo.Task{}, fmt.Errorf("%s", err)
	}

	defer f.Close()
	return task, nil
}

func (taskServer) List(ctx context.Context, empty *todo.Empty) (*todo.TaskList, error) {
	contents, err := ioutil.ReadFile(dbPath)
	if err != nil {
		return &todo.TaskList{}, fmt.Errorf("could not read %s: %v", dbPath, err)
	}
	var tasks *todo.TaskList
	err = json.Unmarshal(contents, tasks)
	fmt.Println(tasks)
	return tasks, nil
}
