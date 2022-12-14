package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
)

var port string
var host string
var remoteMetricsEndpoint string
var allowArray []string
var blockArray []string

func getEnv(key, def string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return def
}

func stringInSlice(s string, strList []string) bool {
	for _, substr := range strList {
		if strings.Contains(s, substr) {
			return true
		}
	}
	return false
}

func seperatorSplit(r rune) bool {
	return r == ',' || r == '|' || r == ';'
}

func filterMetrics(text string) string {
	var ret []string
	var hasAllowed bool = len(allowArray) > 0
	var hasBlocked bool = len(blockArray) > 0
	lines := strings.Split(text, "\n")

	for _, line := range lines {
		var isAllowed bool = (hasAllowed && stringInSlice(line, allowArray))
		var isBlocked bool = (hasBlocked && stringInSlice(line, blockArray))

		if hasAllowed {
			if isAllowed && !isBlocked {
				ret = append(ret, line)
			}

			continue
		}

		if isBlocked {
			continue
		}

		ret = append(ret, line)
	}

	return strings.Join(ret, "\n")
}

func statusHandler(w http.ResponseWriter, r *http.Request) {
	_, _ = w.Write([]byte("OK"))
}

func requestHandler(w http.ResponseWriter, r *http.Request) {
	resp, err := http.Get(remoteMetricsEndpoint)

	if err != nil {
		w.Write([]byte(fmt.Sprintf("requestHandler error %s\n", err.Error())))
		return
	}

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		w.Write([]byte(fmt.Sprintf("requestHandler error %s\n", err.Error())))
		return
	}

	_, _ = w.Write([]byte(filterMetrics(string(body))))
}

func main() {
	port = getEnv("PORT", "9200")
	host = getEnv("HOST", "0.0.0.0")
	remoteMetricsEndpoint = getEnv("REMOTE_METRICS_ENDPOINT", "")
	allowList := getEnv("ALLOW_LIST", "")
	blockList := getEnv("BLOCK_LIST", "")

	allowArray = strings.FieldsFunc(allowList, seperatorSplit)
	blockArray = strings.FieldsFunc(blockList, seperatorSplit)

	http.HandleFunc("/healthz", statusHandler)
	http.HandleFunc("/", requestHandler)

	log.Printf("Authentication service started on %s:%s...\n", host, port)
	http.ListenAndServe(host+":"+port, nil)
}
