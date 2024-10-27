package core

type TaskRepository interface {
	FindAll() ([]Task, error)
	SaveAll(tasks []Task) error
	Save(task Task) error
	DeleteById(id string) error
}
