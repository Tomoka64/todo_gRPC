package main

import (
	"golang.org/x/net/context"

	"flag"
	"fmt"
	"os"
	"strings"

	todo "github.com/Tomoka64/todoWithGRPC/todo"
	"google.golang.org/grpc"
)

var T = flag.String("t", "clean up my room", "put your todo-list")
var D = flag.String("d", "3000-00-00", "set up a deadline for your todo (format-3000-00-00)")

func main() {
	flag.Parse()

	conn, err := grpc.Dial(":8888", grpc.WithInsecure())
	if err != nil {
		fmt.Fprintf(os.Stderr, "could not connect to backend: %v\n", err)
		os.Exit(1)
	}
	client := todo.NewTasksClient(conn)

	switch cmd := flag.Arg(0); cmd {
	case "list":
		err = list(context.Background(), client)
	case "add":
		err = add(context.Background(), client, strings.Join(flag.Args()[1:], " "))
	default:
		err = fmt.Errorf("unknown subcommand %s", cmd)
	}
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func add(ctx context.Context, client todo.TasksClient, text string) error {
	_, err := client.Add(ctx, &todo.Text{Text: text})
	if err != nil {
		return fmt.Errorf("could not add: %v", err)
	}

	fmt.Println("task added successfully")
	return nil
}

func list(ctx context.Context, client todo.TasksClient) error {
	l, err := client.List(ctx, &todo.Empty{})
	if err != nil {
		return fmt.Errorf("%v", err)
	}
	for _, t := range l.Tasks {
		fmt.Printf("✔️ %s\n", t.Text)
	}
	return nil
}
