package server

import (
	"bytes"
	"fmt"
	"log"
	"os"
	"os/exec"
	"os/signal"
	"path"
	"strings"
	"syscall"
)

type RpcServer struct {
	Addr      string
	Port      string
	AdminPort string
}

func (self *RpcServer) Run(
	outputPath string,
	importPath string,
	protoPaths []string,
) {
	fmt.Println("Initializing RPC server...")
	if outputPath == "" {
		outputPath = os.Getenv("output_path") + "/src/grpc"
	}

	// for safety
	outputPath += "/"
	if _, err := os.Stat(outputPath); os.IsNotExist(err) {
		os.Mkdir(outputPath, os.ModePerm)
	}

	// run admin stub server
	// stub.RunStubServer(stub.Options{
	// 	StubPath: *stubPath,
	// 	Port:     *adminport,
	// 	BindAddr: *adminBindAddr,
	// })

	if len(protoPaths) == 0 {
		log.Fatal("RPC server: Need at least one proto file")
	}

	importDirs := strings.Split(importPath, ",")

	// generate pb.go and grpc server based on proto
	generateProtoc(protocParam{
		protoPath:   protoPaths,
		adminPort:   self.AdminPort,
		grpcAddress: self.Addr,
		grpcPort:    self.Port,
		output:      outputPath,
		imports:     importDirs,
	})

	// build the server
	//buildServer(output)

	// and run
	run, runerr := runGrpcServer(outputPath)

	var term = make(chan os.Signal)
	signal.Notify(term, syscall.SIGTERM, syscall.SIGKILL, syscall.SIGINT)
	select {
	case err := <-runerr:
		log.Fatal(err)
	case <-term:
		fmt.Println("Stopping gRPC Server...")
		run.Process.Kill()
	}
}

type protocParam struct {
	protoPath   []string
	adminPort   string
	grpcAddress string
	grpcPort    string
	output      string
	imports     []string
}

func getProtodirs(protoPath string, imports []string) []string {
	// deduced protodir from protoPath
	splitpath := strings.Split(protoPath, "/")
	protodir := ""
	if len(splitpath) > 0 {
		protodir = path.Join(splitpath[:len(splitpath)-1]...)
	}

	// search protodir prefix
	protodirIdx := -1
	for i := range imports {
		dir := path.Join("protogen", imports[i])
		if strings.HasPrefix(protodir, dir) {
			protodir = dir
			protodirIdx = i
			break
		}
	}

	protodirs := make([]string, 0, len(imports)+1)
	protodirs = append(protodirs, protodir)
	// include all dir in imports, skip if it has been added before
	for i, dir := range imports {
		if i == protodirIdx {
			continue
		}
		protodirs = append(protodirs, dir)
	}
	return protodirs
}

func generateProtoc(param protocParam) {
	param.protoPath = fixGoPackage(param.protoPath)
	protodirs := getProtodirs(param.protoPath[0], param.imports)

	// estimate args length to prevent expand
	args := make([]string, 0, len(protodirs)+len(param.protoPath)+2)
	for _, dir := range protodirs {
		args = append(args, "-I", dir)
	}

	// the latest go-grpc plugin will generate subfolders under $GOPATH/src based on go_package option
	pbOutput := os.Getenv("output_path") + "/src"

	args = append(args, param.protoPath...)
	args = append(args, "--go_out=plugins=grpc:"+pbOutput)
	args = append(args, fmt.Sprintf("--mockgpt_out=admin-port=%s,grpc-address=%s,grpc-port=%s:%s",
		param.adminPort, param.grpcAddress, param.grpcPort, param.output))
	protoc := exec.Command("protoc", args...)
	protoc.Stdout = os.Stdout
	protoc.Stderr = os.Stderr
	err := protoc.Run()
	if err != nil {
		log.Fatal("Fail on protoc ", err)
	}

}

// append gopackage in proto files if doesn't have any
func fixGoPackage(protoPaths []string) []string {
	fixgopackage := exec.Command("fix_gopackage.sh", protoPaths...)
	buf := &bytes.Buffer{}
	fixgopackage.Stdout = buf
	fixgopackage.Stderr = os.Stderr
	err := fixgopackage.Run()
	if err != nil {
		log.Println("error on fixGoPackage", err)
		return protoPaths
	}

	return strings.Split(strings.TrimRight(buf.String(), "\n"), "\n")
}

func runGrpcServer(output string) (*exec.Cmd, <-chan error) {
	run := exec.Command("go", "run", output+"server.go")
	run.Stdout = os.Stdout
	run.Stderr = os.Stderr
	err := run.Start()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("grpc server pid: %d\n", run.Process.Pid)
	runerr := make(chan error)
	go func() {
		runerr <- run.Wait()
	}()
	return run, runerr
}
