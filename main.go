package main

import (
	"fmt"
	"log"

	"github.com/ireina7/mockgpt/server"
	"github.com/ireina7/mockgpt/utils"
	"github.com/joho/godotenv"
)

func main() {
	// Load environment args
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	// Init logger
	err = utils.InitLogger()
	if err != nil {
		log.Fatal(err)
	}

	// Initializing app
	grpcServer := new(server.RpcServer)
	stubServer := new(server.StubServer)
	driver := Driver{
		Name:       "MockGPT",
		Version:    "v0.1",
		OutputPath: "./output",
		GrpcServer: *grpcServer,
		StubServer: *stubServer,
	}
	err = driver.Run()
	if err != nil {
		log.Fatal(err)
	}
	// utils.Use(driver)
	fmt.Println("Hello, MockGPT")
}
