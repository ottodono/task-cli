package core

import "errors"

type TaskService struct {
	taskRepository TaskRepository
}

func NewTaskService(taskRepository TaskRepository) *TaskService {
	return &TaskService{taskRepository: taskRepository}
}

func (service *TaskService) FindAll() ([]Task, error) {
	return service.taskRepository.FindAll()
}

func (service *TaskService) Save(task Task) (Task, error) {
	err := service.taskRepository.Save(task)
	return task, err
}

func (service *TaskService) Complete(id string) error {
	tasks, err := service.taskRepository.FindAll()
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
		return errors.New("task with the id not found")
	}

	tasks[index].Complete()

	return service.taskRepository.SaveAll(tasks)
}

func (service *TaskService) DeleteById(id string) error {
	return service.taskRepository.DeleteById(id)
}
