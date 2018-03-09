package hw

import "net/http"

var Version = "unset"

// OK simply return 200
func OK(w http.ResponseWriter, req *http.Request) {}

// V houses version information
type V struct {
	Hostname string `json:"hostname"`
	V        string `json:"version"`
}
