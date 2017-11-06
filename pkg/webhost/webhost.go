// Package webhost handles the logic of the webhosting panel
package webhost

import (
	"io"
	"net/http"
	"os"

	"github.com/Ennovar/gPanel/pkg/api"
	"github.com/Ennovar/gPanel/pkg/logging"
	"github.com/Ennovar/gPanel/pkg/routing"
)

type PrivateHost struct {
	Directory string
}

// NewPrivateHost returns a new PrivateHost type.
func NewPrivateHost() PrivateHost {
	return PrivateHost{
		Directory: "document_roots/webhost/",
	}
}

// ServeHTTP function routes all requests for the private webhost server. It is used in the main
// function inside of the http.ListenAndServe() function for the private webhost host.
func (priv *PrivateHost) ServeHTTP(res http.ResponseWriter, req *http.Request) {
	path := req.URL.Path[1:]
	if len(path) == 0 {
		path = (priv.Directory + "index.html")
	} else {
		path = (priv.Directory + path)
	}

	if reqAuth(path) {
		if !checkAuth(res, req) {
			http.Error(res, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
			return
		}
	}

	isApi, _ := api.HandleAPI(path, res, req)

	if isApi {
		// API methods handle HTTP logic from here
		return
	}

	f, err := os.Open(path)

	if err != nil {
		routing.HttpThrowStatus(http.StatusNotFound, res)
		logging.Console(logging.PRIVATE_PREFIX, logging.NORMAL_LOG, "Path \""+path+"\" rendered a 404 error.")
		return
	}

	contentType, err := routing.GetContentType(path)

	if err != nil {
		routing.HttpThrowStatus(http.StatusUnsupportedMediaType, res)
		logging.Console(logging.PUBLIC_PREFIX, logging.NORMAL_LOG, "Path \""+path+"\" content type could not be determined, 404 error.")
		return
	}

	res.Header().Add("Content-Type", contentType)
	_, err = io.Copy(res, f)

	if err != nil {
		routing.HttpThrowStatus(http.StatusInternalServerError, res)
		logging.Console(logging.PUBLIC_PREFIX, logging.NORMAL_LOG, "Path \""+path+"\" rendered a 500 error.")
		return
	}
}
