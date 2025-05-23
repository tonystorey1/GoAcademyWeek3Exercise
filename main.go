package main

import (
	"Basic_CLI_Application/handler"
	"Basic_CLI_Application/middleware"
	"Basic_CLI_Application/store"
	"Basic_CLI_Application/utils"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
)

// TODO: Introduce actor pattern to encapsulate a user Id for each todo
// TODO: Move handlers into go routines
// TODO: Complete the test.http file
// TODO: Convert to JSON file format
// TODO: Sort superfluous HTTP responses

func main() {
	utils.SetupLogger()

	csvFile := utils.OpenOrCreateFile()
	err := store.Open(csvFile)
	if err != nil {
		log.Fatal(err)
	}

	done := make(chan bool)
	exitWatch(done)

	mux := http.NewServeMux()

	mux.Handle("GET /static/", http.StripPrefix("/static", http.FileServer(http.Dir("./static/"))))

	fmt.Println("Adding HTTP handlers")
	stack := []middleware.Middleware{middleware.ContextMiddleware, middleware.LogMiddleware}
	mux.HandleFunc("/get/", middleware.CompileMiddleware(handler.HandleGet, stack))
	mux.HandleFunc("/add/", middleware.CompileMiddleware(handler.HandleAdd, stack))
	mux.HandleFunc("/update/", middleware.CompileMiddleware(handler.HandleUpdate, stack))
	mux.HandleFunc("/delete/", middleware.CompileMiddleware(handler.HandleDelete, stack))

	srv := &http.Server{Handler: mux}

	fmt.Println("Listening... ")
	ln, err := net.Listen("tcp", "localhost:3000")
	if err != nil {
		utils.Logger.Fatalln(err)
	}

	fmt.Println("Serving... ")
	err = srv.Serve(ln)
	if err != nil {
		utils.Logger.Fatalln(err)
	}

	<-done
	store.Close()
	fmt.Println("Store closed")
}

func exitWatch(done chan bool) {
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		sig := <-sigs
		utils.Logger.Println("\rReceived signal:", sig)
		fmt.Println("\rReceived", sig)
		done <- true
	}()
}
