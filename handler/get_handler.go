package handler

import (
	"Basic_CLI_Application/consts"
	"Basic_CLI_Application/store"
	"Basic_CLI_Application/utils"
	"Basic_CLI_Application/writers"
	"html/template"
	"net/http"
	"net/url"
	"strconv"
	"strings"
)

func HandleGet(writer http.ResponseWriter, request *http.Request) {
	utils.Logger.SetPrefix(request.Context().Value("TraceID ").(string))
	utils.Logger.Println("Calling the get handler")

	if request.Method != http.MethodGet {
		utils.Logger.Println("Method is not " + http.MethodGet)
		writers.WriteResponse(writer, http.StatusMethodNotAllowed)
		return
	}

	requestUrl, err := url.ParseRequestURI(request.RequestURI)
	if err != nil {
		utils.Logger.Println(err.Error())
		writers.WriteResponse(writer, http.StatusBadRequest)
		return
	}

	utils.Logger.Println("Path = " + requestUrl.Path + " ")
	urlArgs := strings.Split(requestUrl.Path[1:], consts.PathSeparator)
	if len(urlArgs) < 2 || len(urlArgs[consts.UrlHttpVerb]) == 0 {
		utils.Logger.Println("Error: null-length todo item")
		writers.WriteResponseWithMessage(writer, http.StatusBadRequest, "Error: null-length todo item")
		return
	}

	utils.Logger.Println("Get entered with: " + strconv.Itoa(len(urlArgs)) + " items")

	if len(urlArgs) == 3 {
		// Get a specific todo
		userId, _ := strconv.Atoi(urlArgs[consts.UrlTodoUserId])
		todoNumber, _ := strconv.Atoi(urlArgs[consts.UrlTodoNumber])
		data := store.SortedTodos(userId, todoNumber)
		writers.WriteResponseWithMessage(writer, http.StatusOK, data[0])

	} else {
		t, err := template.New("ToDoItems").ParseFiles("templates/layout.html")
		if err != nil {
			utils.Logger.Println(err.Error())
			writers.WriteResponse(writer, http.StatusBadRequest)
			return
		}

		userId, _ := strconv.Atoi(urlArgs[consts.UrlTodoUserId])
		data := store.SortedTodos(userId, -1)

		if len(data) == 0 {
			writers.WriteResponseWithMessage(writer, http.StatusOK, "No records found!")
		} else {
			err = t.ExecuteTemplate(writer, "T", data)
			if err != nil {
				utils.Logger.Println(err.Error())
				writers.WriteResponse(writer, http.StatusBadRequest)
				return
			}
		}
	}
}
