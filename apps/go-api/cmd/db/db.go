package main

import (
  "context"
  "encoding/json"
  "fmt"

  // adapt "demo" to your module name if it differs
	"github.com/adityaadpandey/memory-x/go-api/prisma/db"
)

func main() {
  if err := run(); err != nil {
    panic(err)
  }
}

func run() error {
  client := db.NewClient()
  if err := client.Prisma.Connect(); err != nil {
    return err
  }

  defer func() {
    if err := client.Prisma.Disconnect(); err != nil {
      panic(err)
    }
  }()

  ctx := context.Background()

  // create a post
  createdPost, err := client.Post.CreateOne(
    db.Post.Title.Set("Hi from Prisma!"),
    db.Post.Published.Set(true),
    db.Post.Desc.Set("Prisma is a database toolkit and makes databases easy."),
  ).Exec(ctx)
  if err != nil {
    return err
  }

  result, _ := json.MarshalIndent(createdPost, "", "  ")
  fmt.Printf("created post: %s\n", result)

  // find a single post
  post, err := client.Post.FindMany().Exec(ctx)
  if err != nil {
    return err
  }

  result, _ = json.MarshalIndent(post, "", "  ")
  if result == nil {
    fmt.Println("No posts found")
    return nil
  }
  fmt.Printf("post: %s\n", result)
  return nil
}
