package http

import (
	"encoding/json"
	"errors"
	"time"
)

type TaskDTO struct {
	Title string
	Text  string
}

func (t *TaskDTO) ValidateOnCreate() error {
	if t.Title == "" {
		return errors.New("title is empty")
	}

	if t.Text == "" {
		return errors.New("text is empty")
	}

	return nil
}

type ErrorDTO struct {
	Message string
	Time    time.Time
}

func (e *ErrorDTO) ToString() string {
	b, err := json.MarshalIndent(e, "", "    ")

	if err != nil {
		panic(err)
	}

	return string(b)
}
