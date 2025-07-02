package server

import (
	"net/http"

	"github.com/lf-edge/ekuiper/internal/pip"
)

func handlePipListQuery(w http.ResponseWriter, r *http.Request) {
	got, err := pip.GetPipInstallList()
	if err != nil {
		handleError(w, err, "", logger)
		return
	}
	jsonResponse(got, w, logger)
}
