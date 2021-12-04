/*
libLoFloccusDav - LoFloccus WebDav Server

Complied with:
- go version go1.14.2 windows/amd64
- go version go1.14.2 darwin/amd64

How to install/compile this as a static library:
1. go get "golang.org/x/net/webdav"
2. Build ir for your architecture:
2.1. Windows: go build -buildmode c-archive -o libLoFloccusDavWin64.a libLoFloccusDav.go
2.2. macOS: go build -buildmode c-shared -o libLoFloccusDavDarwin.a libLoFloccusDav.go
 */

package main 

import (
	"C"
	"context"
	"golang.org/x/net/webdav"
	"log"
	"net/http"
	"strings"
)

type ServerConfig struct {
	address string
	port    string
	dir     string
	user    string
	passwd  string
}

var (
	serverHandle *http.Server
)

//export serverStart
func serverStart(configAddress *C.char, configPort *C.char, configDir *C.char, configUser *C.char, configPassword *C.char) {

	serverConfig := ServerConfig{
		address: C.GoString(configAddress),
		port:    C.GoString(configPort),
		dir:     C.GoString(configDir),
		user:    C.GoString(configUser),
		passwd:  C.GoString(configPassword),
	}


	log.Printf("WEBDAV: Starting server at " + serverConfig.address + ":" + serverConfig.port + "...")

	mux := http.NewServeMux()
	srvWebdav := &webdav.Handler{
		FileSystem: webdav.Dir(serverConfig.dir),
		LockSystem: webdav.NewMemLS(),
		Logger: func(request *http.Request, err error) {
			if err != nil {
				log.Printf("WEBDAV [%s]: %s, ERROR: %s\n", request.Method, request.URL, err)
			} else {
				log.Printf("WEBDAV [%s]: %s \n", request.Method, request.URL)
			}
		},
	}

	serverHandle = &http.Server{
		Addr:    serverConfig.address + ":" + serverConfig.port,
		Handler: mux,
	}

	mux.HandleFunc("/", func(response http.ResponseWriter, request *http.Request) {

		// Authentication
		username, password, ok := request.BasicAuth()
		if !ok {
			response.Header().Set("WWW-Authenticate", `Basic realm="Restricted"`)
			response.WriteHeader(http.StatusUnauthorized)
			return
		}
		if username != serverConfig.user || password != serverConfig.passwd {
			http.Error(response, "WebDAV: need authorized!", http.StatusUnauthorized)
			return
		}

		// Restrict WebDav to the current folder & read/writes to .xbel files
		if (!strings.HasSuffix(request.RequestURI, ".xbel") && !strings.HasSuffix(request.RequestURI, ".xbel.lock") && request.RequestURI != "/") || (request.RequestURI == "/" && (request.Method != "HEAD" && request.Method != "PROPFIND")) {
			errorFsAccessMsg := "LoFloccus: unauthorized filesystem access detected. LoFloccus WebDAV server is restricted to '*.xbel' and '*.xbel.lock' files."
			log.Printf(request.RequestURI)
			log.Printf(request.Method)
			log.Printf(errorFsAccessMsg)
			http.Error(response, errorFsAccessMsg, http.StatusUnauthorized)
			return
		}

		srvWebdav.ServeHTTP(response, request)
	})

	go func() {
		if err := serverHandle.ListenAndServe(); err != http.ErrServerClosed {
			log.Fatalf("Error with WebDAV server: %s", err)
		}
	}()

}

//export serverStop
func serverStop() {
	log.Printf("WEBDAV: Shutting down...")
	if err := serverHandle.Shutdown(context.TODO()); err != nil {
		log.Fatalf("WEBDAV: error shutting down - %s", err)
	}
	serverHandle = nil
	log.Printf("WEBDAV: Server is down.")
}

func main() {
	// Need a main function to make CGO compile package as C shared library
}
