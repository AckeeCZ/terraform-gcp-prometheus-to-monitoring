package main

import (
	"io"
	"net/http"

	dto "github.com/prometheus/client_model/go"
	"github.com/prometheus/common/expfmt"
	zlog "github.com/rs/zerolog/log"
)

func parseMF(data io.Reader) (map[string]*dto.MetricFamily, error) {
	var parser expfmt.TextParser
	mf, err := parser.TextToMetricFamilies(data)
	if err != nil {
		return nil, err
	}
	return mf, nil
}

func main() {
	resp, err := http.Get("https://api.notificator.cz/metrics")
	if err != nil {
		zlog.Print(err)
	}

	mf, err := parseMF(io.Reader(resp.Body))
	if err != nil {
		zlog.Print(err)
	}

	config := make(map[string]*string)
	load_vars(config)

	sendToGCPMonitoring(mf, config)
}
