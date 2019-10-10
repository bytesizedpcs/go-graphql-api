package gql

import (
  "github.com/bytesizedpcs/go-graphql-api/postgres"
  "github.com/graphql-go/graphql"
)

type Root struct {
  Query *graphql.Object
}

func NewRoot(db *postgres.Db) *Root {
  // Create resolver holding our database
  resolver := Resolver{db: db}

  // create a new root that describes our base query set up
  // user query that takes one argument called name
  root := Root{
    Query: graphql.NewObject(
      graphql.ObjectConfig{
        Name: "Query"
        Fields: graphql.Fields{
          "users": &graphql.Field{
            Type: graphql.NewList(User),
            Args: graphql.FieldConfigArgument{
              "name": &graphql.ArgumentConfig{
                Type: graphql.String,
              },
            },
            Resolve: resolver.UserResolver,
          },
        },
      },
    ),
  }
  return &root
}
