package dbclient

import (
	"log"

	"github.com/adityaadpandey/memory-x/go-api/prisma/db"
)

var PrismaClient *db.PrismaClient

func InitClient() error {
	PrismaClient = db.NewClient()

	if err := PrismaClient.Prisma.Connect(); err != nil {
		return err
	}
	return nil
}

// Disconnect the Prisma client
func Disconnect() {
	if err := PrismaClient.Prisma.Disconnect(); err != nil {
		log.Fatal("Error disconnecting Prisma client:", err)
	}
}
