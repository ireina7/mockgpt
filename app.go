package main

import (
	"flag"

	"github.com/ireina7/mockgpt/server"
	"github.com/ireina7/mockgpt/utils"
)

type Driver struct {
	Name       string
	Version    string
	OutputPath string
	GrpcServer server.RpcServer
	StubServer server.StubServer
}
type CmdOptions struct{}

func (driver *Driver) Run() error {
	err := driver.parse()
	if err != nil {
		return err
	}
	return nil
}

func (self *Driver) parse() error {
	self.OutputPath = *flag.String(
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
	flag.Parse()
	utils.Use(grpcBindAddr, grpcPort, adminport, imports)
	return nil
}
