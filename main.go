package main

import (
  "fmt"
  "log"
  "net/http"

  "github.com/bytesizedpcs/go-graphql-api/gql"
  "github.com/bytesizedpcs/go-graphql-api/postgres"
  "github.com/bytesizedpcs/go-graphql-api/server"
  "github.com/go-chi/chi"
  "github.com/go-chi/chi/middleware"
  "github.com/go-chi/render"
  "github.com/graphql-go/graphql"
)

func main() {
  router, db := initializeAPI()
  defer db.Close()

  log.Fatal(httpListenAndServer(":4000", router))
}

func initializeAPI() (*chi.Mux, *postgres.Db) {
  //create a new router

  router := chi.NewRouter()

  // postgres username and database
  db, err := postgres.New(
    postgres.ConnString("localhost", 5432, "", ""),
  )
  if err != nil { log.Fatal(err) }

  // Create our root query for graphql
  rootQuery := gql.NewRoot(db)

  // Create a new graphql schema, passing in the root query
  sc, err := graphql.NewSchema(
    graphql.SchemaConfig{ Query: rootQuery.Query },
  )
  if err != nil {
    fmt.Println("Error creating schema: ", err)
  }

  // Create a server struct that holds a pointer to our databse as well
  // as the address of our graphql schema
  s := server.Server {
    GqlSchema: &sc,
  }

  // Add some middleware to our router
  router.Use(
    render.SerContentType(render.ContentTypeJSON), // set content-type headers as app/json
    middleware.Logger, // log api request calls 
    middleware.DefaultCompress, // compress results, mostly gzipping assets and json
    middleware.StripSlashes, // match paths with a trailing slash, strip, and route tru mux
    middleware.Recoverer, //recover from panics without crashing server
  )

  // Create the graphql route with a Server mehtod to handle it
  router.Post("/graphql", s.Graphql())

  return router, db

}
