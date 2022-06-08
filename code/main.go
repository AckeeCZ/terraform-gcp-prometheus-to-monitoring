package p

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"

	zlog "github.com/rs/zerolog/log"
)

func Consume(w http.ResponseWriter, r *http.Request) {
	var d struct {
		Service  string `json:"service"`
		Endpoint string `json:"endpoint"`
	}

	if err := json.NewDecoder(r.Body).Decode(&d); err != nil {
		switch err {
		case io.EOF:
			zlog.Error().Msg(fmt.Sprintf("json io.EOF: %v", err))
			return
		default:
			zlog.Error().Msg(fmt.Sprintf("json.NewDecoder: %v", err))
			http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
			return
		}
	}
	config := make(map[string]*string)
	load_vars(config)
	config["SERVICE"] = &d.Service

	resp, err := http.Get(d.Endpoint)
	if err != nil {
		zlog.Error().Msg(fmt.Sprintf("http call fail: %v", err))
		os.Exit(1)
	}

	mf, err := parseMF(io.Reader(resp.Body))
	if err != nil {
		zlog.Error().Msg(fmt.Sprintf("http parsing fail: %v", err))
		os.Exit(1)
	}

	sendToGCPMonitoring(mf, config)
}
