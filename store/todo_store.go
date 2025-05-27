package store

import (
	"Basic_CLI_Application/consts"
	"Basic_CLI_Application/utils"
	"bufio"
	"encoding/csv"
	"errors"
	"fmt"
	"log"
	"os"
	"sort"
	"strconv"
	"sync"
)

type Todo struct {
	UserId     int
	TodoNumber int
	TodoItem   string
	Status     string
}

type TodoStore struct {
	store map[string]Todo
	mutex sync.RWMutex
	Open  bool
}

var (
	// todoStore The Todo store itself
	todoStore TodoStore

	fileMutex sync.RWMutex

	ErrStoreNotOpen = errors.New("store not open")

	ErrStoreAlreadyOpen = errors.New("store already open")
)

// Open opens the store with the supplied CSV file and loads the contents. If csvFile is nil it will open empty.
func Open(csvFile *os.File) error {
	if todoStore.Open {
		return ErrStoreAlreadyOpen
	} else {
		err := todoStore.openStore(csvFile)
		if err != nil {
			utils.Logger.Fatalln(err)
			return err
		}
		fmt.Println("Store opened")
		return nil
	}
}

// Close - closes the store and removes all records
func Close() {
	if !todoStore.Open {
		utils.Logger.Println("Store already closed")
		return
	} else {
		todoStore.mutex.Lock()
		defer todoStore.mutex.Unlock()
		todoStore.store = nil
		todoStore.Open = false
		fmt.Println("Store closed")
	}
}

// Count is the count of TODO items in the store
func Count() int {
	todoStore.mutex.Lock()
	defer todoStore.mutex.Unlock()
	return len(todoStore.store)
}

func GetRecord(userId string, todoNumber string) Todo {
	todoStore.mutex.Lock()
	defer todoStore.mutex.Unlock()
	return todoStore.store[getKey(userId, todoNumber)]
}

func PutRecord(userId string, todoNumber string, description string, newStatus string) error {
	if !todoStore.Open {
		return ErrStoreNotOpen
	}

	if userId == "" || !isNumeric(userId) {
		// TODO: Log error
		return errors.New("userId is invalid")
	}

	if todoNumber == "" || !isNumeric(todoNumber) {
		// TODO: Log error
		return errors.New("todoNumber is invalid")
	}

	key := getKey(userId, todoNumber)

	user, _ := strconv.Atoi(userId)
	number, _ := strconv.Atoi(todoNumber)
	todoStore.mutex.Lock()
	defer todoStore.mutex.Unlock()
	_, exists := todoStore.store[key]
	if !exists {
		todoStore.store[key] = Todo{user, number, description, newStatus}
	} else {
		// Return an error indicating item already exists
		return errors.New(userId + " " + todoNumber + " already exists")
	}
	return nil
}

func UpdateRecord(userId string, todoNumber string, description string, newStatus string) error {
	if !todoStore.Open {
		return ErrStoreNotOpen
	}

	if userId == "" || !isNumeric(userId) {
		// TODO: Log error
		return errors.New("userId is invalid")
	}

	if todoNumber == "" || !isNumeric(todoNumber) {
		// TODO: Log error
		return errors.New("todoNumber is invalid")
	}

	key := getKey(userId, todoNumber)

	user, _ := strconv.Atoi(userId)
	number, _ := strconv.Atoi(todoNumber)
	todoStore.mutex.Lock()
	defer todoStore.mutex.Unlock()
	_, exists := todoStore.store[key]
	if exists {
		todoStore.store[key] = Todo{user, number, description, newStatus}
	} else {
		// Return an error indicating item does not exist
		return errors.New(userId + " " + todoNumber + " does not exist")
	}
	return nil
}

func RemoveRecord(userId string, todoNumber string) error {
	if !todoStore.Open {
		return ErrStoreNotOpen
	}

	if userId == "" || !isNumeric(userId) {
		// TODO: Log error
		return errors.New("userId is invalid")
	}

	if todoNumber == "" || !isNumeric(todoNumber) {
		// TODO: Log error
		return errors.New("todoNumber is invalid")
	}

	key := getKey(userId, todoNumber)

	todoStore.mutex.Lock()
	defer todoStore.mutex.Unlock()
	_, exists := todoStore.store[key]
	if exists {
		delete(todoStore.store, key)
	} else {
		// Return an error indicating does not exist
		return errors.New(todoNumber + " does not exist")
	}
	return nil
}

func SortedTodos() []string {
	if !todoStore.Open {
		return []string{}
	}

	data := make([]string, Count())
	i := 0
	for _, record := range todoStore.store {
		data[i] = fmt.Sprintf("[UserId]: %s [TodoNumber]: %s [Todo Item]: %s [Todo Status]: %s \n", strconv.Itoa(record.UserId), strconv.Itoa(record.TodoNumber), record.TodoItem, record.Status)
		i++
	}
	sort.Strings(data)
	return data
}

func WriteTodosToFile() error {
	fileMutex.Lock()
	defer fileMutex.Unlock()

	// TODO: Not very efficient - change to improve

	writer, file, err := utils.CreateCSVWriter(consts.FileName)
	if err != nil {
		utils.Logger.Fatalln("Error creating CSV writer:", err)
		return err
	}

	// Question: Should the range operator be wrapped in a mutex?
	for _, record := range todoStore.store {
		_ = writer.Write([]string{strconv.Itoa(record.UserId), strconv.Itoa(record.TodoNumber), record.TodoItem, record.Status})
	}

	writer.Flush()
	_ = file.Close()
	return nil
}

func putRecordsInMap(records [][]string) error {
	for _, record := range records {
		err := PutRecord(record[0], record[1], record[2], record[3])
		if err != nil {
			return err
		}
	}
	return nil
}

func loadTodoItems(csvFile *os.File) error {
	reader := csv.NewReader(bufio.NewReader(csvFile))
	records, err := reader.ReadAll()
	if err != nil {
		utils.Logger.Fatalln(err)
	}
	e := csvFile.Close()
	if e != nil {
		utils.Logger.Fatalln(err)
	}

	// Read records and put into the store
	fmt.Println("Opening the list of todo items: ")
	err = putRecordsInMap(records)
	if err != nil {
		utils.Logger.Fatalln(err)
	}
	return nil
}

func getKey(userId string, todoNumber string) string {
	return userId + "_" + todoNumber
}

func isNumeric(s string) bool {
	_, err := strconv.ParseFloat(s, 64)
	return err == nil
}

func (store *TodoStore) openStore(csvFile *os.File) error {
	store.store = make(map[string]Todo)
	store.Open = true

	if csvFile != nil {
		err := loadTodoItems(csvFile)
		if err != nil {
			log.Fatal(err)
			return err
		}
	}
	return nil
}
