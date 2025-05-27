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

	if requestUrl == nil {
		utils.Logger.Println("URL is nil")
		writers.WriteResponse(writer, http.StatusBadRequest)
		return
	}

	utils.Logger.Println("Path = " + requestUrl.Path + " ")
	urlArgs := strings.Split(requestUrl.Path[1:], consts.PathSeparator)
	utils.Logger.Println("Get entered with: " + strconv.Itoa(len(urlArgs)) + " items")

	if len(urlArgs) == 3 {
		// Get a specific todo
		userId, _ := strconv.Atoi(urlArgs[consts.UrlTodoNumber])
		todo := store.SortedTodos()[userId]
		writers.WriteResponseWithMessage(writer, http.StatusOK, todo)

	} else {
		t, err := template.New("ToDoItems").ParseFiles("templates/layout.html")
		if err != nil {
			utils.Logger.Println(err.Error())
			writers.WriteResponse(writer, http.StatusBadRequest)
			return
		}
		err = t.ExecuteTemplate(writer, "T", store.SortedTodos())
		if err != nil {
			utils.Logger.Println(err.Error())
			writers.WriteResponse(writer, http.StatusBadRequest)
			return
		}
	}
}
