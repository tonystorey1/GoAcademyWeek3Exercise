package handler

import (
	"Basic_CLI_Application/consts"
	"Basic_CLI_Application/store"
	"Basic_CLI_Application/utils"
	"Basic_CLI_Application/writers"
	"net/http"
	"net/url"
	"strconv"
	"strings"
)

func HandleAdd(writer http.ResponseWriter, request *http.Request) {
	utils.Logger.SetPrefix(request.Context().Value("TraceID ").(string))
	utils.Logger.Println("Calling the Add handler")

	if request.Method != http.MethodPost {
		utils.Logger.Println("Method is not " + http.MethodPost)
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

	utils.Logger.Println("Path = " + requestUrl.Path)
	urlArgs := strings.Split(requestUrl.Path[1:], consts.PathSeparator)
	if len(urlArgs) < 5 || len(urlArgs[0]) == 0 {
		utils.Logger.Println("Error: null-length todo item")
		writers.WriteResponse(writer, http.StatusBadRequest)
		return
	}

	if !utils.IsStatusValid(urlArgs[consts.UrlTodoStatus]) {
		utils.Logger.Println("Error: supplied status is not valid")
		writers.WriteResponseWithMessage(writer, http.StatusBadRequest, "Error: supplied status is not valid, it must be either  \"not started\", \"started\", or \"completed\"")
		return
	}

	err = store.PutRecord(urlArgs[consts.UrlSegmentUserId], urlArgs[consts.UrlTodoNumber], urlArgs[consts.UrlTodoDescription], strings.ToLower(urlArgs[consts.UrlTodoStatus]))
	if err != nil {
		utils.Logger.Print(err)
		writers.WriteResponseWithMessage(writer, http.StatusBadRequest, err.Error())
		return
	} else {
		n, err := writer.Write([]byte("Item added: " + urlArgs[consts.UrlSegmentUserId]))
		if err != nil {
			utils.Logger.Println("Error calling writer.Write: " + err.Error())
			writers.WriteResponseWithMessage(writer, http.StatusBadRequest, err.Error())
			return
		}

		err = store.WriteTodosToFile()
		if err != nil {
			utils.Logger.Println(err.Error())
			writers.WriteResponseWithMessage(writer, http.StatusInternalServerError, err.Error())
			return
		}

		utils.Logger.Println("Item added" + strconv.Itoa(n))
		writers.WriteResponse(writer, http.StatusOK)
	}
}
