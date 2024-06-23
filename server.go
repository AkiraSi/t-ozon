package main

import (
	"log"
	"net/http"
	"os"
	"t-ozon/graph"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/gofor-little/env"
)

const defaultPort = "8080"

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = defaultPort
	}
	if err := env.Load(".env"); err != nil {
		panic(err)
	}
	typeMemory := env.Get("TYPE_MEMORY", "False")
	resolver, err := graph.NewResolver(typeMemory)
	if err != nil {
		log.Fatalf("failed to create resolver: %v", err)
	}
	srv := handler.NewDefaultServer(graph.NewExecutableSchema(graph.Config{Resolvers: resolver}))
	http.Handle("/", playground.Handler("GraphQL playground", "/query"))
	http.Handle("/query", srv)
	log.Printf("connect to http://localhost:%s/ for GraphQL playground", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
