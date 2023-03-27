package main

import (
	"flag"
	"fmt"
	"log"

	"github.com/ireina7/mockgpt/utils"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	err = utils.InitLogger()
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Initializing flags...")

	outputPointer := flag.String(
		"o", "", "directory to output server.go. Default is $output_path/src/grpc/",
	)
	grpcPort := flag.String("grpc-port", "4770", "Port of gRPC tcp server")
	grpcBindAddr := flag.String(
		"grpc-listen", "", "Adress the gRPC server will bind to. Default to localhost, set to 0.0.0.0 to use from another machine",
	)
	adminport := flag.String("admin-port", "4771", "Port of stub admin server")
	// adminBindAddr := flag.String("admin-listen", "", "Adress the admin server will bind to. Default to localhost, set to 0.0.0.0 to use from another machine")
	// stubPath := flag.String("stub", "", "Path where the stub files are (Optional)")
	imports := flag.String(
		"imports", "/protobuf", "comma separated imports path. default path /protobuf is where mockgpt Dockerfile install WKT protos",
	)
	// for backwards compatibility
	// if os.Args[0] == "gripmock" {
	// 	os.Args = append(os.Args[:1], os.Args[2:]...)
	// }

	utils.Use(outputPointer, grpcBindAddr, grpcPort, adminport, imports)
	flag.Parse()
	fmt.Println("Hello, mockgpt")
}
