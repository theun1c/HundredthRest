package todo

import "time"

type Task struct {
	Title       string
	Text        string
	IsDone      bool
	CreatedAt   time.Time
	CompletedAt *time.Time // if task is not done - nil, else - completed time
	// cause its not required field
}

func NewTask(title, text string) Task {

	return Task{
		Title:       title,
		Text:        text,
		IsDone:      false,
		CreatedAt:   time.Now(),
		CompletedAt: nil,
	}
}

func (t *Task) Complete() {
	now := time.Now()
	t.IsDone = true
	t.CompletedAt = &now
}

func (t *Task) Uncomplete() {
	t.IsDone = false
	t.CompletedAt = nil
}
