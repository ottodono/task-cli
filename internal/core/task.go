package core

import (
	"fmt"
	"github.com/ottodono/task-cli/pkg/utils"
	"time"
)

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
	fmt.Printf("%s %s %s %t\n", task.id, task.content, utils.FormatTimeToString(task.createdDate), task.complete)
}
