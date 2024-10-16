package main

import (
	"encoding/csv"
	"errors"
	"fmt"
	"os"
	"strings"
	"syscall"
	"time"

	uuid "github.com/google/uuid"
)

const (
	FILE_NAME string = "todo.csv"
)

// Structures
type Task struct {
	id          string
	content     string
	createdDate time.Time
	complete    bool
}

func NewTask(id string, content string, createdDate time.Time, complete bool) *Task {
	return &Task{
		id:          id,
		content:     content,
		createdDate: createdDate,
		complete:    complete,
	}
}

func (task *Task) GetId() string {
	return task.id
}

func (task *Task) GetContent() string {
	return task.content
}

func (task *Task) GetCreatedDate() time.Time {
	return task.createdDate
}

func (task *Task) GetComplete() bool {
	return task.complete
}

func (task *Task) Complete() {
	task.complete = true
}

func (task *Task) Afficher() {
	fmt.Printf("%s %s %s %t\n", task.id, task.content, formatTimeToString(task.createdDate), task.complete)
}

// Date
func formatTimeToString(date time.Time) string {
	str := date.Format(time.RFC3339)
	return str
}

func formatStringToTime(str string) time.Time {
	date, err := time.Parse(time.RFC3339, str)
	if err != nil {
		fmt.Println("Could not parse time:", err)
	}
	return date
}

// Main
func main() {
	if len(os.Args) < 2 {
		fmt.Println("Not enough arguments")
		return
	}

	command := os.Args[1]
	if command == "add" {
		if len(os.Args) < 3 {
			fmt.Println("Not enough arguments for the command add")
			return
		} else {
			id := strings.Split(uuid.NewString(), "-")[0]
			content := os.Args[2]
			task := NewTask(id, content, time.Now(), false)
			save(task)
		}
	} else if command == "list" {
		tasks, err := findAll()
		if err != nil {
			fmt.Println(err.Error())
		}
		for _, task := range tasks {
			task.Afficher()
		}
	} else if command == "delete" {
		if len(os.Args) < 3 {
			fmt.Println("Not enough arguments for the command add")
			return
		} else {
			id := os.Args[2]
			deleteTask(id)
		}
	} else if command == "complete" {
		if len(os.Args) < 3 {
			fmt.Println("Not enough arguments for the command add")
			return
		} else {
			id := os.Args[2]
			completeTask(id)
		}
	} else {
		fmt.Println("Unknown command.")
		return
	}
}

func save(task *Task) error {
	tasks, err := findAll()
	if err != nil {
		return err
	}
	tasks = append(tasks, *task)
	err = saveAll(tasks)
	if err != nil {
		return err
	}
	return nil
}

func saveAll(tasks []Task) error {
	records := tasksToRecords(tasks)

	file, err := loadFile(FILE_NAME)
	if err != nil {
		return err
	}
	defer closeFile(file)

	if err := file.Truncate(0); err != nil {
		return err
	}

	if _, err := file.Seek(0, 0); err != nil {
		return err
	}

	writer := csv.NewWriter(file)

	for _, record := range records {
		if err := writer.Write(record); err != nil {
			return err
		}
	}

	writer.Flush()
	if err := writer.Error(); err != nil {
		return err
	}
	return nil
}

func findAll() ([]Task, error) {
	file, err := loadFile(FILE_NAME)
	if err != nil {
		return nil, err
	}
	defer closeFile(file)

	reader := csv.NewReader(file)
	records, err := reader.ReadAll()
	if err != nil {
		return nil, err
	}

	tasks := recordsToTasks(records)

	return tasks, nil
}

func deleteTask(id string) error {
	tasks, err := findAll()
	if err != nil {
		return err
	}

	index := -1
	for i, task := range tasks {
		if task.GetId() == id {
			index = i
		}
	}
	if index == -1 {
		return errors.New("Task with the id not found")
	}

	tasks = remove(tasks, index)

	err = saveAll(tasks)
	if err != nil {
		return err
	}
	return nil
}

func completeTask(id string) error {
	tasks, err := findAll()
	if err != nil {
		return err
	}

	index := -1
	for i, task := range tasks {
		if task.GetId() == id {
			index = i
		}
	}
	if index == -1 {
		return errors.New("Task with the id not found")
	}

	tasks[index].Complete()

	saveAll(tasks)

	return nil
}

func remove(tasks []Task, index int) []Task {
	fmt.Println(tasks)
	tasks = append(tasks[:index], tasks[index+1:]...)
	fmt.Println(tasks)
	return tasks
}

func recordsToTasks(records [][]string) []Task {
	tasks := make([]Task, 0)
	for _, record := range records {
		task := recordToTask(record)
		tasks = append(tasks, *task)
	}
	return tasks
}

func recordToTask(record []string) *Task {
	return NewTask(
		record[0],
		record[1],
		formatStringToTime(record[2]),
		stringToBool(record[3]),
	)
}

func tasksToRecords(tasks []Task) [][]string {
	records := make([][]string, 0)
	for _, task := range tasks {
		record := taskToRecord(&task)
		records = append(records, record)
	}
	return records
}

func taskToRecord(task *Task) []string {
	date := formatTimeToString(task.GetCreatedDate())
	return []string{
		task.GetId(),
		task.GetContent(),
		date,
		boolToString(task.GetComplete()),
	}
}

func boolToString(value bool) string {
	if value {
		return "true"
	}
	return "false"
}

func stringToBool(value string) bool {
	return value == "true"
}

func loadFile(filepath string) (*os.File, error) {
	f, err := os.OpenFile(filepath, os.O_RDWR|os.O_CREATE, os.ModePerm)
	if err != nil {
		return nil, fmt.Errorf("failed to open file for reading")
	}

	// Exclusive lock obtained on the file descriptor
	if err := syscall.Flock(int(f.Fd()), syscall.LOCK_EX); err != nil {
		_ = f.Close()
		return nil, err
	}

	return f, nil
}

func closeFile(f *os.File) error {
	syscall.Flock(int(f.Fd()), syscall.LOCK_UN)
	return f.Close()
}
