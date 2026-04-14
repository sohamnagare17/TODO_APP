package main

import (
	
	"go-sqlite/database"
	"go-sqlite/handlers"
	"go-sqlite/repository"
	"go-sqlite/routes"
	"go-sqlite/services"
	"go-sqlite/observability"
	"log"
	"net/http"
	"os"
)


// func InitTracer() func(context.Context) error {
// 	ctx := context.Background()
// 	exporter, err := otlptracegrpc.New(ctx,
// 		otlptracegrpc.WithEndpoint("otel-collector:4317"),
// 		otlptracegrpc.WithInsecure(),
// 	)
// 	if err != nil {
// 		log.Fatal(err)
// 	}
// 	tp := sdktrace.NewTracerProvider(
// 		sdktrace.WithBatcher(exporter),
// 		sdktrace.WithResource(resource.NewWithAttributes(
// 			semconv.SchemaURL,
// 			semconv.ServiceName("todo-app"),
// 		)),
// 	)
// 	otel.SetTracerProvider(tp)
// 	return tp.Shutdown
// }

func main() {
	file, err := os.OpenFile("app.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)

	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	log.SetOutput(file)

	dbconn := db.Dbinit()

	defer dbconn.Close()

	repo := repository.NewTaskRepository(dbconn)
	service := services.NewTaskServices(repo)
	taskhandler := handlers.NewTaskHandler(service)

	repouser := repository.NewUserRepository(dbconn)
	userservices := services.NewUserServices(repouser)
	userhandler := handlers.NewUserHandler(userservices)

	routes.SetupRoutes(taskhandler, userhandler, dbconn)

	shutdown := observability.InitTracer()
	defer shutdown()
	
	log.Println("server running on port 8080")
	http.ListenAndServe(":8080", nil)
}
