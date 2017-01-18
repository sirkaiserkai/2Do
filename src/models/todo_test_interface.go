package models

import (
	"fmt"
	"log"
)

var todoMap = make(map[string]Todo)

// the UserStorage interface
type TestTodoStorage struct {
	todos *map[string]Todo
}

func newTestTodoStorage() *TestTodoStorage {
	t := TestTodoStorage{}
	t.todos = &todoMap
	return &t
}
func (tus *TestTodoStorage) Close() {
	log.Println("Closing TestUserStorage")
}

func (tus *TestTodoStorage) GetAllTodos() ([]Todo, error) {
	return nil, nil
}

func (tus *TestTodoStorage) GetTodoById(id string) (*Todo, error) {
	return nil, nil
}

func (tus *TestTodoStorage) GetTodosForUserId(id string) ([]Todo, error) {
	todos := *tus.todos
	if todos == nil {
		log.Fatal("todos map was not assigned in TestTodoStorage!")
	}

	ts := make([]Todo, 0)
	for _, t := range todos {
		if t.Ownerid == id {
			ts = append(ts, t)
		}
	}

	return ts, nil
}

func (tus *TestTodoStorage) InsertTodo(t Todo) error {
	(*tus.todos)[t.Id.Hex()] = t
	return nil
}

func (tus *TestTodoStorage) ModifyTodo(todoId, userId string, changes map[string]interface{}) error {
	return nil
}

func (tus *TestTodoStorage) DeleteTodo(id, userId string) error {
	todos := (*tus.todos)
	t, ok := todos[id]
	if !ok {
		return TodoNotFoundError
	}

	if t.Ownerid != userId {
		return fmt.Errorf("User: %s does not own: %s", userId, id)
	}

	delete(todos, id)
	return nil
}
