package task

import (
	"time"
)

type Task struct {
    Name string
    Created time.Time
    Completed bool
}

type Tasks struct {
    tasks []Task
}

func (list *Tasks) Add(task Task) {
    list.tasks = append(list.tasks, task)
}

func (list Tasks) Tasks() TaskIterator {
    return TaskIterator{
        curr: len(list.tasks),
        tasks: list.tasks[0:len(list.tasks)],
    }
}

type TaskIterator struct {
    curr int
    tasks []Task
}

func (it *TaskIterator) Next() *Task {
    it.curr--
    if it.curr < 0 {
        return nil
    }
    return &(it.tasks[it.curr])
}

func (it *TaskIterator) HasNext() bool {
    return it.curr > 0
}
