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

func HandleUpdate(writer http.ResponseWriter, request *http.Request) {
	utils.Logger.SetPrefix(request.Context().Value("TraceID ").(string))
	utils.Logger.Println("Calling the Update handler")

	if request.Method != http.MethodPut {
		utils.Logger.Println("Method is not " + http.MethodPut)
		writers.WriteResponse(writer, http.StatusMethodNotAllowed)
		return
	}

	requestUrl, err := url.ParseRequestURI(request.RequestURI)
	if err != nil {
		utils.Logger.Println(err.Error())
		writers.WriteResponseWithMessage(writer, http.StatusBadRequest, err.Error())
		return
	}

	if requestUrl == nil {
		utils.Logger.Println("URL is nil")
		writers.WriteResponse(writer, http.StatusBadRequest)
		return
	}

	utils.Logger.Println("Path = " + requestUrl.Path)
	items := strings.Split(requestUrl.Path[1:], consts.PathSeparator)
	if len(items) < 5 || len(items[1]) == 0 {
		utils.Logger.Println("Error: null-length todo item")
		writers.WriteResponseWithMessage(writer, http.StatusBadRequest, "Error: null-length todo item")
		return
	}

	if !utils.IsStatusValid(items[4]) {
		utils.Logger.Println("Error: supplied status is not valid")
		writers.WriteResponseWithMessage(writer, http.StatusBadRequest, "Error: supplied status is not valid, it must be either  \"not started\", \"started\", or \"completed\"")
		return
	}

	err = store.UpdateRecord(items[1], items[2], items[3], strings.ToLower(items[4]))
	if err != nil {
		utils.Logger.Println(err.Error())
		writers.WriteResponseWithMessage(writer, http.StatusBadRequest, err.Error())
		return
	} else {
		n, err := writer.Write([]byte("Item updated: " + items[1]))
		if err != nil {
			utils.Logger.Println("Error calling writer.Write: " + err.Error())
			writers.WriteResponse(writer, http.StatusInternalServerError)
			return
		}

		err = store.WriteTodosToFile()
		if err != nil {
			utils.Logger.Print(err)
			writers.WriteResponse(writer, http.StatusInternalServerError)
			return
		}

		utils.Logger.Println("Item added" + strconv.Itoa(n))
		writers.WriteResponse(writer, http.StatusOK)
	}
}
