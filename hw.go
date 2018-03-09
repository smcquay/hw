package hw

import "net/http"

// Version is version.
var Version = "unset"

// Git exists so we can override it
var Git = "unset"

// OK simply return 200
func OK(w http.ResponseWriter, req *http.Request) {}

// V houses version information
type V struct {
	Hostname string `json:"hostname"`
	V        string `json:"version"`
	G        string `json:"git"`
}
