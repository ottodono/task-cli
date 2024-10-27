package main

import (
	"fmt"
	"github.com/google/uuid"
	"github.com/mergestat/timediff"
	"github.com/ottodono/task-cli/internal/core"
	"github.com/ottodono/task-cli/internal/infra"
	"os"
	"strings"
	"text/tabwriter"
	"time"
)

const (
	FileName string = "todo.csv"
)

func main() {
	taskRepository := infra.NewCsvFileTaskRepository(FileName)
	taskService := core.NewTaskService(taskRepository)
	handleCommand(*taskService)
}

func handleCommand(taskService core.TaskService) {
	if len(os.Args) < 2 {
		fmt.Println("not enough arguments")
	}

	command := os.Args[1]
	switch command {
	case "list":
		handleListCommand(taskService)
	case "add":
		handleAddCommand(taskService)
	case "complete":
		handleCompleteCommand(taskService)
	case "delete":
		handleDeleteCommand(taskService)
	default:
		fmt.Println("unknown command")
	}
}

func handleListCommand(taskService core.TaskService) {
	tasks, err := taskService.FindAll()
	if err != nil {
		fmt.Println(err.Error())
	}
	displayTasks(tasks)
}

func handleAddCommand(taskService core.TaskService) {
	if len(os.Args) < 3 {
		fmt.Println("not enough arguments for the command add")
	} else {
		id := strings.Split(uuid.NewString(), "-")[0]
		arg := os.Args[2]
		task := core.NewTask(id, arg, time.Now(), false)
		_, err := taskService.Save(*task)
		if err != nil {
			fmt.Println(err.Error())
		}
	}
}

func handleCompleteCommand(taskService core.TaskService) {
	if len(os.Args) < 3 {
		fmt.Println("not enough arguments for the command complete")
	} else {
		id := os.Args[2]
		err := taskService.Complete(id)
		if err != nil {
			fmt.Println(err.Error())
		}
	}
}

func handleDeleteCommand(taskService core.TaskService) {
	if len(os.Args) < 3 {
		fmt.Println("not enough arguments for the command delete")
	} else {
		id := os.Args[2]
		err := taskService.DeleteById(id)
		if err != nil {
			fmt.Println(err.Error())
		}
	}
}

func displayTasks(tasks []core.Task) {
	writer := tabwriter.NewWriter(os.Stdout, 0, 0, 3, ' ', 0)
	_, err := fmt.Fprintln(writer, "ID\tTask\tCreated\tDone")
	if err != nil {
		return
	}
	for _, task := range tasks {
		timeToDisplay := timediff.TimeDiff(task.GetCreatedDate())
		_, err := fmt.Fprintf(writer, "%s\t%s\t%s\t%t\n", task.GetId(), task.GetContent(), timeToDisplay, task.GetComplete())
		if err != nil {
			return
		}
	}
	err = writer.Flush()
	if err != nil {
		return
	}
}
