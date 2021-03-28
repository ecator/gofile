package server

import (
	"log"
	"net/http"
	"path/filepath"
	"strconv"
	"sync"
	"time"

	"github.com/julienschmidt/httprouter"
)

// Server infomation
type ServerInfo struct {
	ListenAddr string
	ListenPort int
	WebDir     string
	DataDir    string
	ExpireSec  int
}

var serverInfo ServerInfo

// StartServer starts server
func StartServer(si ServerInfo) {
	serverInfo = si
	// set routers
	router := httprouter.New()

	// not found
	router.NotFound = http.HandlerFunc(handleNotFound)

	// web
	router.GET("/", handleIndex)
	router.GET("/index.ico", handleIndex)
	router.GET("/404.ico", handleIndex)
	router.ServeFiles("/assets/*filepath", http.Dir(filepath.Join(serverInfo.WebDir, "assets")))

	// file
	router.GET("/file/:token", handleFileDown)
	router.POST("/file", handleFileUp)
	router.GET("/file", handleFileInfos)

	// start http server
	log.Println("Starting service...")
	errCh := make(chan error)
	var wg sync.WaitGroup
	listenAddr := si.ListenAddr + ":" + strconv.Itoa(si.ListenPort)
	go func() {
		wg.Add(1)
		errCh <- http.ListenAndServe(listenAddr, router)
		wg.Done()
	}()

	select {
	case err := <-errCh:
		log.Fatalln(err)
	case <-time.After(time.Second * 3):
		log.Println("Listening at " + listenAddr)
		wg.Wait()
	}
}
