package todo

import "sync"

type List struct {
	tasks map[string]Task
	mtx   sync.RWMutex
}

func NewList() *List {
	return &List{
		tasks: make(map[string]Task),
	}
}

// с точки зрения принципа разделения ответственности
// лучше всего в методы листа передавать уже созданную из вне задачу
// вместо передачи параметров текста и заголовка внутри методов
// AddTask(task) <-- better when --> AddTask(title, text ...)

// ресивер должен быть по указателю
// иначе мы значение будем менять у копии
// а не у изначального объекта !
func (l *List) AddTast(task Task) error {
	// стоит выбор сохранения таск
	// 1 - слайс
	// 2 - мапа
	//
	// выбор: мапа (ключ - значение)
	// обоснование выбора:
	// удобнее всего, тк удаление/получение будет реализовано
	// через заголовок задачи (ключ)
	// а создание задачи будет через 2 параметра
	// название (ключ) - текст (значение)
	//
	// слайсом мы бы прогоняли каждый раз по циклу - неэффективно
	//
	// также в структуре листа мы бы удобно хранили значение таскс
	// мапа где ключ - строка (тайтл) а значение - объект таск

	// если такая задача существует "по названию"
	// то создаем новую с именет "тайтл + дата сейчас"
	// иначе просто создаем задачу, добавляя таск в мапу
	l.mtx.Lock()
	defer l.mtx.Unlock()

	if _, ok := l.tasks[task.Title]; ok {
		return ErrTaskAlreadyExist
	}

	l.tasks[task.Title] = task

	return nil
}

func (l *List) DeleteTask(title string) error {
	l.mtx.Lock()
	defer l.mtx.Unlock()

	if _, ok := l.tasks[title]; ok {
		return ErrTaskNotFound
	}

	delete(l.tasks, title)

	return nil
}

func (l *List) CompleteTask(title string) error {

	l.mtx.Lock()
	defer l.mtx.Unlock()

	task, ok := l.tasks[title]
	if !ok {
		return ErrTaskNotFound
	}

	// можео сделать и так, но следует придерживаться
	// принципа единой ответственности и выделить логику отметки
	// в отдельную функцию структуры таск
	//
	// task.IsDone = true
	// now := time.Now()
	// task.CompletedAt = &now

	task.Complete()

	l.tasks[title] = task

	return nil
}

// нужно обратить внимание!
// возвращать так мапу чревато какими-либо последствиями
//
// нужно нагуглить проблему
func (l *List) GetTasks() map[string]Task {
	l.mtx.RLock()
	defer l.mtx.RUnlock()

	tmp := make(map[string]Task, len(l.tasks))

	for k, v := range l.tasks {
		tmp[k] = v
	}

	return tmp
}

func (l *List) ListNotCompletedTasks() map[string]Task {
	l.mtx.RLock()
	defer l.mtx.RUnlock()

	notCompletedTasks := make(map[string]Task)

	for k, v := range l.tasks {
		if !v.IsDone {
			notCompletedTasks[k] = v
		}
	}

	return notCompletedTasks
}

func (l *List) GetTask(title string) (Task, error) {
	l.mtx.RLock()
	defer l.mtx.RUnlock()

	task, ok := l.tasks[title]
	if !ok {
		return Task{}, ErrTaskNotFound
	}

	return task, nil

}

func (l *List) UncompleteTask(title string) error {
	l.mtx.Lock()
	defer l.mtx.Unlock()

	task, ok := l.tasks[title]
	if !ok {
		return ErrTaskNotFound
	}

	task.Uncomplete()

	l.tasks[title] = task

	return nil
}
