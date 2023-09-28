package main

import (
	"github.com/gin-gonic/gin"
	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/NeurologiaLogic/gochat/graph"
	"github.com/NeurologiaLogic/gochat/websocket"
)

// Defining the Graphql handler
func graphqlHandler() gin.HandlerFunc {
	// NewExecutableSchema and Config are in the generated.go file
	// Resolver is in the resolver.go file
	h := handler.NewDefaultServer(
		graph.NewExecutableSchema(graph.ConfigureResolver()))
	return func(c *gin.Context) {
		h.ServeHTTP(c.Writer, c.Request)
	}
}

// Defining the Playground handler
func playgroundHandler() gin.HandlerFunc {
	h := playground.Handler("GraphQL", "/query")

	return func(c *gin.Context) {
		h.ServeHTTP(c.Writer, c.Request)
	}
}


func main() {
	// Setting up Gin
	// gin.SetMode(gin.ReleaseMode)
	r := gin.Default()
	ws := websocket.NewWebsocketManager()
	go r.POST("/query", graphqlHandler())
	go r.GET("/", playgroundHandler())
	go r.GET("/ws",ws.Handler)
	// r.RunTLS("localhost:8080","server.crt","server.key")
	r.Run("localhost:8080")
}