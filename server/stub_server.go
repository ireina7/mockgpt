package server

type StubServer struct {
	addr string
	port string
}

func (self *StubServer) Run() {
	// if self.port == "" {
	// 	self.port = DEFAULT_PORT
	// }
	// addr := opt.BindAddr + ":" + opt.Port
	// r := chi.NewRouter()
	// r.Post("/add", addStub)
	// r.Get("/", listStub)
	// r.Post("/find", handleFindStub)
	// r.Get("/clear", handleClearStub)

	// if opt.StubPath != "" {
	// 	readStubFromFile(opt.StubPath)
	// }

	// fmt.Println("Serving stub admin on http://" + addr)
	// go func() {
	// 	err := http.ListenAndServe(addr, r)
	// 	log.Fatal(err)
	// }()
}
