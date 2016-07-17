package ws

func init() {
	// InitSearcher run it
	SearcherVar = &searcher{
		directory:  make(map[string]*actor),
		search:     make(chan *searchActor, 256),
		register:   make(chan *registerActor, 256),
		unregister: make(chan string, 256),
		Clean:      make(chan bool),
	}
	go SearcherVar.Run()
}
