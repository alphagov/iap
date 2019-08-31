package internal

import (
	"encoding/json"
	"net/http"

	"github.com/sirupsen/logrus"
)

// JSONResponse should take our data and other attributes and simply present end user with JSON.
func JSONResponse(ctx Context, w http.ResponseWriter, code int, data interface{}) {
	blob, err := json.Marshal(data)
	if err != nil {
		ctx.Logger.WithFields(logrus.Fields{
			"error": err,
		}).Error("failed to parse json blob")

		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(code)
	w.Write(blob)
	return
}
