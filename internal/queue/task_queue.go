package queue

type TaskQueue struct {
	ch chan string
}

func NewTaskQueue(bufferSize int) *TaskQueue {
	return &TaskQueue{ch: make(chan string, bufferSize)}
}

func (q *TaskQueue) Enqueue(taskID string) {
	q.ch <- taskID
}

func (q *TaskQueue) Consume() <-chan string {
	return q.ch
}
