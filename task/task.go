package task

import (
	"fmt"
	"time"
)

type Task struct {
    Name string
    Created time.Time
    Completed bool
    Id int
}

type Tasks struct {
    tasks []Task
}

func (list *Tasks) Add(task Task) Task {
    task.Id = len(list.tasks)
    list.tasks = append(list.tasks, task)
    return task
}

func (list *Tasks) Tasks() TaskIterator {
    length := len(list.tasks)
    return TaskIterator{
        currUncompleted: length,
        currCompleted: length,
        tasks: list.tasks[0:length],
        length: length,
    }
}

func (list *Tasks) Get(id int) (Task, error) {
    if id >= len(list.tasks) {
        return Task{}, fmt.Errorf("Can't find task with ID %d", id)
    }
    return list.tasks[id], nil
}

func (list *Tasks) Update(id int, task Task) error {
    if id >= len(list.tasks) {
        return fmt.Errorf("Can't find task with ID %d", id)
    }
    list.tasks[id] = task
    return nil
}

type TaskIterator struct {
    currUncompleted int
    currCompleted int
    tasks []Task
    length int
}

func (it *TaskIterator) nextCompleted() *Task {
    it.currCompleted--
    for it.currCompleted >= 0 && !it.tasks[it.currCompleted].Completed {
        it.currCompleted--
    }
    if it.currCompleted >= 0 {
        return &it.tasks[it.currCompleted]
    }
    return nil
}

func (it *TaskIterator) nextUnCompleted() *Task {
    it.currUncompleted--
    for it.currUncompleted >= 0 && it.tasks[it.currUncompleted].Completed {
        it.currUncompleted--
    }
    if it.currUncompleted >= 0 {
        return &it.tasks[it.currUncompleted]
    }
    return nil
}


func (it *TaskIterator) Next() *Task {
    next := it.nextUnCompleted()
    if next == nil {
        next = it.nextCompleted()
    }
    if next != nil {
        it.length--
    }
    return next
}

func (it *TaskIterator) HasNext() bool {
    return it.length > 0
}
