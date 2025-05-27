package handler

import (
	"Basic_CLI_Application/consts"
	"Basic_CLI_Application/store"
	"Basic_CLI_Application/utils"
	"Basic_CLI_Application/writers"
	"net/http"
	"net/url"
	"strings"
)

func HandleDelete(writer http.ResponseWriter, request *http.Request) {
	utils.Logger.SetPrefix(request.Context().Value("TraceID ").(string))
	utils.Logger.Println("Calling the Delete handler")

	if request.Method != http.MethodDelete {
		utils.Logger.Println("Method is not " + http.MethodDelete)
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
	items := strings.Split(requestUrl.Path[1:], consts.PathSeparator)
	if len(items) < 3 || len(items[1]) == 0 {
		utils.Logger.Println("Error: null-length todo item")
		writers.WriteResponseWithMessage(writer, http.StatusBadRequest, "Error: null-length todo item")
		return
	}

	err = store.RemoveRecord(items[1], items[2])
	if err != nil {
		utils.Logger.Print(err)
		writers.WriteResponseWithMessage(writer, http.StatusInternalServerError, err.Error())
		return
	} else {
		err = store.WriteTodosToFile()
		if err != nil {
			utils.Logger.Print(err)
			writers.WriteResponseWithMessage(writer, http.StatusInternalServerError, err.Error())
			return
		}

		utils.Logger.Println("Item deleted" + items[1])
		_, err := writer.Write([]byte("Item deleted: " + items[1]))
		if err != nil {
			utils.Logger.Println("Error calling writer.Write: " + err.Error())
			writers.WriteResponseWithMessage(writer, http.StatusInternalServerError, err.Error())
			return
		}

		writers.WriteResponse(writer, http.StatusOK)
		return
	}
}
