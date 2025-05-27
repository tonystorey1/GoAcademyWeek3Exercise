package main

import (
	"Basic_CLI_Application/consts"
	"Basic_CLI_Application/store"
	"fmt"
	"github.com/stretchr/testify/assert"
	"log"
	"strconv"
	"sync"
	"testing"
)

func Test_CanPutRecord(t *testing.T) {
	teardownTest := setupTest(t)
	defer teardownTest(t)

	err := store.PutRecord("1", "1", "A Todo", "notStarted")
	if err != nil {
		t.Errorf("Test_CanPutRecordInMap error = %q", err.Error())
	}

	if store.Count() != 1 {
		t.Error("New Todo not added")
	}
}

func Test_CanUpdateRecord(t *testing.T) {
	teardownTest := setupTest(t)
	defer teardownTest(t)

	// Arrange
	const userId = "1"
	const todoNumber = "1"
	err := store.PutRecord(userId, todoNumber, "A Todo", consts.TodoStatusNotStarted)
	if err != nil {
		t.Errorf("Test_CanUpdateRecord error adding test todo to map = %q", err.Error())
	}

	// Act
	err = store.UpdateRecord(userId, todoNumber, "A different todo", consts.TodoStatusStarted)
	if err != nil {
		t.Errorf("Test_CanUpdateRecord error updating test todo to map = %q", err.Error())
	}

	// Assert
	todo := store.GetRecord(userId, todoNumber)
	assert.Equal(t, todo.TodoItem, "A different todo")
	assert.Equal(t, todo.Status, consts.TodoStatusStarted)
}

func Test_CannotUpdateRecord(t *testing.T) {
	teardownTest := setupTest(t)
	defer teardownTest(t)

	// Assert
	const userId = "1"
	const todoNumber = "1"
	todo := store.GetRecord(userId, todoNumber)
	assert.Equal(t, todo.TodoItem, "")
	assert.Equal(t, todo.Status, "")
}

func Test_CanRemoveRecord(t *testing.T) {
	teardownTest := setupTest(t)
	defer teardownTest(t)

	// Arrange
	const userId = "1"
	const todoNumber = "1"
	err := store.PutRecord(userId, todoNumber, "A Todo", consts.TodoStatusNotStarted)
	if err != nil {
		t.Errorf("Test_CanRemoveRecord error adding test todo to map = %q", err.Error())
	}

	// Act
	err = store.RemoveRecord(userId, todoNumber)
	if err != nil {
		t.Errorf("Test_CanRemoveRecord error updating test todo to map = %q", err.Error())
	}

	// Assert
	assert.Equal(t, store.Count(), 0)
}

func Test_CountReturnsExpectedCount(t *testing.T) {
	teardownTest := setupTest(t)
	defer teardownTest(t)

	// Arrange
	const userId = "1"
	const count = 10
	for i := 0; i < count; i++ {
		err := store.PutRecord(userId, strconv.Itoa(i), "Todo:"+strconv.Itoa(i), consts.TodoStatusNotStarted)
		if err != nil {
			t.Errorf("Test_CountReturnsExpectedCount error adding test todo to map = %q", err.Error())
		}
	}

	// Assert
	assert.Equal(t, store.Count(), count)
}

func Test_CanGetIndividualRecord(t *testing.T) {
	teardownTest := setupTest(t)
	defer teardownTest(t)

	// Arrange
	const userId = "1"
	const count = 10
	for i := 0; i < count; i++ {
		err := store.PutRecord(userId, strconv.Itoa(i), "Todo:"+strconv.Itoa(i), consts.TodoStatusNotStarted)
		if err != nil {
			t.Errorf("Test_CountReturnsExpectedCount error adding test todo to map = %q", err.Error())
		}
	}

	// Assert
	todo := store.GetRecord(userId, "7")
	assert.Equal(t, todo.TodoNumber, 7)
	assert.Equal(t, todo.TodoItem, "Todo:"+strconv.Itoa(7))
}

func Test_SortedTodosReturnsExpected(t *testing.T) {
	teardownTest := setupTest(t)
	defer teardownTest(t)

	// Arrange
	const userId = "1"
	const count = 10
	for i := 0; i < count; i++ {
		err := store.PutRecord(userId, strconv.Itoa(i), "Todo:"+strconv.Itoa(i), consts.TodoStatusNotStarted)
		if err != nil {
			t.Errorf("Test_CountReturnsExpectedCount error adding test todo to map = %q", err.Error())
		}
	}

	// Assert
	todo := store.SortedTodos()
	assert.Equal(t, len(todo), count)
	assert.Equal(t, todo[4], "[UserId]: 1 [TodoNumber]: 4 [Todo Item]: Todo:4 [Todo Status]: not started \n")
}

func Test_CanConcurrently_PutRecords(t *testing.T) {
	teardownTest := setupTest(t)
	defer teardownTest(t)

	// Arrange
	const users = 500
	const todosToCreate = 1000
	var wg sync.WaitGroup
	for i := 0; i < users; i++ {
		wg.Add(1)
		go createTodos(i, todosToCreate, &wg, t)
	}
	wg.Wait()

	fmt.Println("Waiting done")

	count := store.Count()
	assert.Equal(t, count, users*todosToCreate)
	todo := store.GetRecord(strconv.Itoa(users-1), strconv.Itoa(todosToCreate-1))
	assert.Equal(t, todo.TodoItem, "Todo:"+strconv.Itoa(todosToCreate-1))
}

func Test_CanConcurrently_RemoveRecords(t *testing.T) {
	teardownTest := setupTest(t)
	defer teardownTest(t)

	// Arrange
	const users = 500
	const todosToCreate = 1000
	var addWg sync.WaitGroup
	for i := 0; i < users; i++ {
		addWg.Add(1)
		go createTodos(i, todosToCreate, &addWg, t)
	}
	addWg.Wait()

	fmt.Println("Todo's added, count is: " + strconv.Itoa(store.Count()))

	// Act
	var deleteWg sync.WaitGroup
	for i := 0; i < users; i++ {
		deleteWg.Add(1)
		go deleteTodos(i, todosToCreate, &deleteWg, t)
	}
	deleteWg.Wait()

	fmt.Println("Todo's removed, count is: " + strconv.Itoa(store.Count()))

	// Assert
	count := store.Count()
	assert.Equal(t, count, 0)
}

func createTodos(userId int, count int, wg *sync.WaitGroup, t *testing.T) {
	defer wg.Done()

	user := strconv.Itoa(userId)
	for i := 0; i < count; i++ {
		err := store.PutRecord(user, strconv.Itoa(i), "Todo:"+strconv.Itoa(i), consts.TodoStatusNotStarted)
		if err != nil {
			t.Errorf("createTodos error adding test todo to map = %q", err.Error())
		}
	}
}

func deleteTodos(userId int, count int, wg *sync.WaitGroup, t *testing.T) {
	defer wg.Done()

	user := strconv.Itoa(userId)
	for i := 0; i < count; i++ {
		err := store.RemoveRecord(user, strconv.Itoa(i))
		if err != nil {
			t.Errorf("deleteTodos error removing test todo to map = %q", err.Error())
		}
	}
}

func setupTest(tb testing.TB) func(tb testing.TB) {
	log.Println("setup test")
	err := store.Open(nil)
	if err != nil {
		tb.Fatal(err)
	}

	return func(tb testing.TB) {
		log.Println("teardown test")
		store.Close()
	}
}
