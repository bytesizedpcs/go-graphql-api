package postgres

import (
  "database/sql"
  "fmt"

  // postgres driver
  _ "github.com/lib/pq"
)

// Db is our database struct used for interacting with the database
type Db struct {
  *sql.DB
}

func New(connString string) (*Db, error) {
  db, err := sql.Open("postgres", connString)
  if err != nil {
    return nil, err
  }

  err = db.Ping()
  if err != nil {
    return nil, err
  }

  return &Db{db}, nil
}

func ConnString(host string, port int, user string, dbName string) string {
  return fmt.Sprintf(
    "host=%s port=%d user=%s dbname=%s sslmode=disable",
    host, port, user, dbName,
  )
}

// User shape
type User struct {
  ID int
  Name string
  Age int
  Profession string
  Friendly bool
}

// is called within our user query for graphql
func (d *Db) GetUsersByName(name string) []User {
  stmt, err := d.Prepare("Select * FROM users WHERE name=$1")
  if err != nil {
    fmt.Println("GetUsersByName Preparation Err: ", err)
  }

  rows, err := stmt.Query(name)
  if err != nil {
    fmt.Println("GetUsersByName Query Err:", err)
  }

  // Create user struct for holding each row's data
  var r User

  users := []User{}

  for rows.Next() {
    err = rows.Scan(
      &r.ID,
      &r.Name,
      &r.Age,
      &r.Profession,
      &r.Friendly,
    )
    if err != nil {
      fmt.Println("Error scanning rows: ", err)
    }
    users = append(users, r)
  }

  return users
}
