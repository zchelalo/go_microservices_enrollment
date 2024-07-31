package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/joho/godotenv"
	"github.com/zchelalo/go_microservices_enrollment/internal/enrollment"
	"github.com/zchelalo/go_microservices_enrollment/pkg/bootstrap"
	"github.com/zchelalo/go_microservices_enrollment/pkg/handler"
)

func main() {
	logger := bootstrap.InitLogger()

	err := godotenv.Load()
	if err != nil {
		logger.Fatal("error loading .env file")
	}

	db, err := bootstrap.DBConnection()
	if err != nil {
		logger.Fatal(err)
	}

	pageLimDef := os.Getenv("PAGINATOR_LIMIT_DEFAULT")
	if pageLimDef == "" {
		logger.Fatal("paginator limit default is required")
	}

	config := enrollment.Config{
		LimPageDef: pageLimDef,
	}

	ctx := context.Background()
	enrollmentRepository := enrollment.NewRepository(logger, db)
	enrollmentService := enrollment.NewService(logger, enrollmentRepository)
	hdler := handler.NewEnrollmentHTTPServer(ctx, enrollment.MakeEndpoints(enrollmentService, config))

	port := os.Getenv("PORT")
	address := fmt.Sprintf("127.0.0.1:%s", port)

	server := &http.Server{
		// Handler:      http.TimeoutHandler(router, 5*time.Second, "Timeout!"),
		Handler:      accessControl(hdler),
		Addr:         address,
		WriteTimeout: 5 * time.Second,
		ReadTimeout:  5 * time.Second,
	}

	errCh := make(chan error)
	go func() {
		logger.Println("listen in ", address)
		errCh <- server.ListenAndServe()
	}()

	err = <-errCh
	if err != nil {
		log.Fatal(err)
	}
}

func accessControl(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PATCH, DELETE, OPTIONS, HEAD")
		w.Header().Set("Access-Control-Allow-Headers", "Accept, Authorization, Cache-Control, Content-Type, DNT, If-Modified-Since, Keep-Alive, Origin, User-Agent, X-Requested-With")

		if req.Method == "OPTIONS" {
			return
		}

		h.ServeHTTP(w, req)
	})
}
