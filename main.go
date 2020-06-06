package main

import (
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	log "github.com/sirupsen/logrus"
	"os"
	"runtime"
	"strconv"
	"strings"
)

var logger = log.New()

func main() {
	file, err := os.OpenFile(os.Getenv("NTM_LOG_FILE"), os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0644)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()
	setupLogger(file)
	sanityCheck()
	a := App{}
	a.Initialize()
	a.RunExternal()
}

func setupLogger(file *os.File) {
	logger.SetOutput(file)
	logger.SetLevel(log.InfoLevel)
	//useful only if we are sending logs to some log consolidator
	//logger.SetFormatter(&log.JSONFormatter{})
	//logger.SetFormatter(&log.TextFormatter{})
	logger.SetReportCaller(true)

	// https://github.com/sirupsen/logrus/issues/63
	logger.SetFormatter(&log.TextFormatter{
		CallerPrettyfier: func(f *runtime.Frame) (string, string) {
			var filename string
			var function string
			if strings.Contains(f.File, "github.com") { // removing the file and function name from the logs for dependencies
				filename = ""
				function = ""
			} else {
				repopath := fmt.Sprintf("%s/src/ntm-backend/", os.Getenv("GOPATH"))
				filename = strings.Replace(f.File, repopath, "", -1) + ":" + strconv.Itoa(f.Line)
				function = strings.Replace(f.Function, "ntm-backend/", "", -1)
			}
			return function, filename
		},
	})
}

func sanityCheck() {
	envNames := []string{
		"NTM_SERVER_IP",
		"NTM_SERVER_PORT",
		"NTM_LOG_FILE",
		"ELASTICSEARCH_URL",
		"SNMP_ELASTICSEARCH_INDEX",
		"NTM_API_KEY",
	}
	notFound := []string{}
	for _, name := range envNames {
		if os.Getenv(name) == "" {
			notFound = append(notFound, name)
		}
	}
	if len(notFound) > 0 {
		logger.Printf("[ %v ] Environment variable is not defined", strings.Join(notFound, ","))
		fmt.Printf("[ %v ] Environment variable is not defined", strings.Join(notFound, ","))
		os.Exit(1)
	}
}
