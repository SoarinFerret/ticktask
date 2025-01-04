package todotxt

import (
	"time"

	"github.com/KEINOS/go-todotxt/todo"
	"github.com/soarinferret/ticktask/internal/config"
	"github.com/soarinferret/ticktask/internal/profile"
)

// AddTask adds a task to the todo.txt file
func AddTask(task todo.Task) (*todo.Task, error) {
	// get active profile
	projects, contexts := profile.GetActiveProfileFilter()

	// add projects and contexts to task
	task.Projects = append(task.Projects, projects...)
	task.Contexts = append(task.Contexts, contexts...)

	// add created date
	task.CreatedDate = time.Now()

	// load tasks
	list, err := todo.LoadFromPath(config.GetTodoTxtPath())
	if err != nil {
		return nil, err
	}

	// add task
	list.AddTask(&task)

	// save tasks
	err = list.WriteToPath(config.GetTodoTxtPath())
	if err != nil {
		return nil, err
	}

	return &task, nil
}

// GetTasks returns all tasks from the todo.txt file
func GetTasks() (todo.TaskList, error) {
	// get active profile
	projects, contexts := profile.GetActiveProfileFilter()

	// load tasks only from active projects and contexts
	list, err := todo.LoadFromPath(config.GetTodoTxtPath())
	if err != nil {
		return nil, err
	}

	// filter tasks to only those that match the active profile
	for _, project := range projects {
		list = list.Filter(todo.FilterByProject(project))
	}
	for _, context := range contexts {
		list = list.Filter(todo.FilterByContext(context))
	}

	return list, nil
}

func GetTask(id int) (*todo.Task, error) {
	// load tasks
	list, err := GetTasks()
	if err != nil {
		return nil, err
	}

	// get task by id
	task, err := list.GetTask(id)
	if err != nil {
		return nil, err
	}

	return task, nil
}

func RemoveTask(id int) error {
	// check if task exists in context
	_, err := GetTask(id)
	if err != nil {
		return err
	}

	// load tasks
	list, err := todo.LoadFromPath(config.GetTodoTxtPath())
	if err != nil {
		return err
	}

	// remove task
	list.RemoveTaskByID(id)

	// save tasks
	err = list.WriteToPath(config.GetTodoTxtPath())
	if err != nil {
		return err
	}

	return nil
}

// Definitely a better way to do this...
func CompleteTask(id int) (*todo.Task, error) {
	_, err := GetTask(id)
	if err != nil {
		return nil, err
	}

	// task exists in context, so now we need to get a full list of tasks
	// so when we save, it doesn't overwrite the list with just the filtered task list
	list, err := todo.LoadFromPath(config.GetTodoTxtPath())
	if err != nil {
		return nil, err
	}

	// get task by id
	task, err := list.GetTask(id)
	if err != nil {
		return nil, err
	}

	// complete task
	task.Complete()

	// save tasks
	err = list.WriteToPath(config.GetTodoTxtPath())
	if err != nil {
		return nil, err
	}

	return task, nil
}

func ReopenTask(id int) (*todo.Task, error) {
	_, err := GetTask(id)
	if err != nil {
		return nil, err
	}

	// task exists in context, so now we need to get a full list of tasks
	// so when we save, it doesn't overwrite the list with just the filtered task list
	list, err := todo.LoadFromPath(config.GetTodoTxtPath())
	if err != nil {
		return nil, err
	}

	// get task by id
	task, err := list.GetTask(id)
	if err != nil {
		return nil, err
	}

	// incomplete task
	task.Reopen()

	// save tasks
	err = list.WriteToPath(config.GetTodoTxtPath())
	if err != nil {
		return nil, err
	}

	return task, nil
}

func AddTimeToTask(id int, duration time.Duration, override bool) (*todo.Task, error) {
	_, err := GetTask(id)
	if err != nil {
		return nil, err
	}

	// task exists in context, so now we need to get a full list of tasks
	// so when we save, it doesn't overwrite the list with just the filtered task list
	list, err := todo.LoadFromPath(config.GetTodoTxtPath())
	if err != nil {
		return nil, err
	}

	// get task by id
	task, err := list.GetTask(id)
	if err != nil {
		return nil, err
	}

	// add time to task
	if _, exists := task.AdditionalTags["time"]; !exists {
		if task.AdditionalTags == nil {
			task.AdditionalTags = make(map[string]string)
		}
		task.AdditionalTags["time"] = duration.String()
	} else {
		// add time to existing time
		existingDuration, err := time.ParseDuration(task.AdditionalTags["time"])
		if err != nil {
			return nil, err
		}
		if override {
			existingDuration = 0
		}
		task.AdditionalTags["time"] = (existingDuration + duration).String()
	}

	// save tasks
	err = list.WriteToPath(config.GetTodoTxtPath())
	if err != nil {
		return nil, err
	}

	return task, nil
}

func SetPriority(id int, priority string) (*todo.Task, error) {
	_, err := GetTask(id)
	if err != nil {
		return nil, err
	}

	// task exists in context, so now we need to get a full list of tasks
	// so when we save, it doesn't overwrite the list with just the filtered task list
	list, err := todo.LoadFromPath(config.GetTodoTxtPath())
	if err != nil {
		return nil, err
	}

	// get task by id
	task, err := list.GetTask(id)
	if err != nil {
		return nil, err
	}

	// add priority
	task.Priority = priority

	// save tasks
	err = list.WriteToPath(config.GetTodoTxtPath())
	if err != nil {
		return nil, err
	}

	return task, nil
}
