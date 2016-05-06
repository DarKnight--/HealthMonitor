// The package implements the web interface for the OWTF monitor module
package webui

import (
	"io"
	"net/http"
)

func home(w http.ResponseWriter, r *http.Request) {
	io.WriteString(w, "OWTF Monitor")
}

func RunServer(port string) {
	http.HandleFunc("/", home)
	http.ListenAndServe(":"+port, nil)
}
