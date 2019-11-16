package main

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin/binding"
	"log"
	"net/http"
	"os"
	"os/signal"
	"spm/model"
	"spm/routes"
	"spm/service"
	"spm/util"
	"syscall"
	"time"
)

var help = `Parameters:
	migrate
	reset
	seed
	run
`

func runserver() {
	service.Start()

	// setuo router
	r := routes.APIRouter()

	srv := &http.Server{
		Addr:    "127.0.0.1:3000",
		Handler: r,
	}

	go func() {
		// service connections
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
	}()

	// Wait for interrupt signal to gracefully shutdown the server with
	// a timeout of 5 seconds.
	quit := make(chan os.Signal, 1)
	// kill (no param) default send syscanll.SIGTERM
	// kill -2 is syscall.SIGINT
	// kill -9 is syscall. SIGKILL but can"t be catch, so don't need add it
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutdown Server ...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("Server Shutdown: ", err)
	}

	log.Println("Server exiting")
}

func main() {
	model.Connect()
	defer model.Close()
	// set validator to v9
	binding.Validator = new(util.DefaultValidator)

	args := os.Args[1:]
	if len(args) == 0 {
		runserver()
	} else {
		switch args[0] {
		case "migrate":
			model.Migrate()
		case "reset":
			model.Reset()
		case "seed":
			model.Seed()
		case "run":
			runserver()
		default:
			fmt.Print(help)
		}
	}
}
