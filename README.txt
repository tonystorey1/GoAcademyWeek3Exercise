Golang Academy Week 3 Concurrency/Actor Pattern
===============================================
Week 3 code to add concurrency and actor pattern to the code produced in weeks 1 & 2.

Solution developed with Goland 2025.1
-------------------------------------
Code itself has been run and debugged via Goland using go version 1.24 (select the main entry point)

Unit tests are supplied in the main_test.go file (again, run in Goland).

The code format of the todo is:

type Todo struct {
	UserId     int
	TodoNumber int
	TodoItem   string
	Status     string
}

Persisting the TODO's
----------------------
The todos are currently persisted to a CSV file (there is an outstanding TODO to convert to JSON file format)

The format of the file is:

<userId>,<Todo Number>,<Todo Item>,<status>

1,1,do something,not started
1,2,make a cuppa,started
2,3,wash the car,not started
3,4,do some shopping,done

No header is required.

Todo Status
-----------
Can be either "not started", "started", or "completed"

Commands (with Curl examples)
=============================

1. Add: adds a new TODO to the store
curl -X POST http://localhost:3000/add/<userId>/<todo item to add>/<status>

2. Delete: deletes an existing todo from the store
curl -X DELETE http://localhost:3000/delete/<userId>/<todo number to delete>

3. Update: updates an existing todo in the store
curl -X PUT http://localhost:3000/update/<userId>/<todo number to add>/<todo to update>/<todo status>

4. Get
This command can either get all todos or the todo's for a single user
curl -X GET http://localhost:3000/get/<userId>/<optional: either a todo number or blank for all todos for the user>

5. About: Returns simple about information
curl -X GET http://localhost:3000/static/about.html




Useful Links (copy from learnamp)
---------------------------------
https://bjss.learnamp.com/en/learnlists/golang-academy?overview=true course overview
https://quii.gitbook.io/learn-go-with-tests good place to start
https://gobyexample.com/ good place to start
https://www.youtube.com/watch?v=YzLrWHZa-Kc reference video, good place to start
https://www.youtube.com/watch?v=PyDMqgOkiR8 more on philosophy of go, method receivers
https://learning.oreilly.com/library/view/learning-go-2nd/9781098139285/ reference bookExtras
https://100go.co/ 100 mistakes
https://golangweekly.com/issues/551 good newsletter to subscribe to
https://go.dev/doc/effective_go go classic doc
https://pkg.go.dev/std standart library reference https://goplay.tools/ playground to run simple scripts




