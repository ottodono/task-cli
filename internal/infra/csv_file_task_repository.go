package infra

import (
	"encoding/csv"
	"errors"
	"github.com/ottodono/task-cli/internal/core"
	"github.com/ottodono/task-cli/pkg/utils"
	"os"
)

type CsvFileTaskRepository struct {
	filePath string
}

func NewCsvFileTaskRepository(filePath string) *CsvFileTaskRepository {
	return &CsvFileTaskRepository{filePath: filePath}
}

func (repository *CsvFileTaskRepository) FindAll() (tasks []core.Task, error error) {
	file, err := utils.LoadFile(repository.filePath)
	if err != nil {
		return nil, err
	}
	defer func(f *os.File) {
		_ = utils.CloseFile(f)
	}(file)

	reader := csv.NewReader(file)
	records, err := reader.ReadAll()
	if err != nil {
		return nil, err
	}

	return recordsToTasks(records), nil
}

func (repository *CsvFileTaskRepository) SaveAll(tasks []core.Task) (error error) {
	records := tasksToRecords(tasks)

	file, err := utils.LoadFile(repository.filePath)
	if err != nil {
		return err
	}
	defer func(f *os.File) {
		_ = utils.CloseFile(f)
	}(file)

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

func (repository *CsvFileTaskRepository) Save(task core.Task) (error error) {
	tasks, err := repository.FindAll()
	if err != nil {
		return err
	}
	tasks = append(tasks, task)
	err = repository.SaveAll(tasks)
	if err != nil {
		return err
	}
	return nil
}

func (repository *CsvFileTaskRepository) DeleteById(id string) (error error) {
	tasks, err := repository.FindAll()
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

	tasks = remove(tasks, index)

	err = repository.SaveAll(tasks)
	if err != nil {
		return err
	}
	return nil
}

func remove(tasks []core.Task, index int) []core.Task {
	tasks = append(tasks[:index], tasks[index+1:]...)
	return tasks
}

func recordsToTasks(records [][]string) []core.Task {
	tasks := make([]core.Task, 0)
	for _, record := range records {
		task := recordToTask(record)
		tasks = append(tasks, *task)
	}
	return tasks
}

func recordToTask(record []string) *core.Task {
	return core.NewTask(
		record[0],
		record[1],
		utils.FormatStringToTime(record[2]),
		stringToBool(record[3]),
	)
}

func tasksToRecords(tasks []core.Task) [][]string {
	records := make([][]string, 0)
	for _, task := range tasks {
		record := taskToRecord(&task)
		records = append(records, record)
	}
	return records
}

func taskToRecord(task *core.Task) []string {
	date := utils.FormatTimeToString(task.GetCreatedDate())
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
