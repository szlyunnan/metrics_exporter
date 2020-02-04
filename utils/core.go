package utils

import (
	"encoding/json"
	"flag"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/spf13/viper"
	"log"
	"metrics_exporter/config"
	"metrics_exporter/model"
	"metrics_exporter/prom"
	"net/http"
	"os"
	"strconv"
	"time"
)

var (
	conf   config.Configure
	listen string
	host   string
	port   int
	path   string
	help   bool
)

func Core() {

	dbfullpath := fmt.Sprintf("%s/%s", conf.Metrics.DbPath, conf.Metrics.DbName)

	// get data from sqlite
	mdb := model.NewMetricsDB("sqlite3", dbfullpath)
	t := time.Now()
	sqlStr := fmt.Sprintf("select si_pkg, si_code, si_msg, si_build from `aax_sign-%d`", t.Day())
	data, err := mdb.DBQuery(sqlStr)
	if err != nil {
		panic(err)
	}

	// MetricsCollector object
	mc := prom.NewMetricsCollector(data)
	reg := prometheus.NewPedanticRegistry()
	reg.MustRegister(mc)

	router := mux.NewRouter().StrictSlash(false)
	// use logging middleware
	router.Use(loggingMiddleware)

	router.HandleFunc("/health", healthzHandler).Methods("GET")
	router.Handle("/metrics", promhttp.HandlerFor(
		//prometheus.DefaultGatherer,
		reg,
		promhttp.HandlerOpts{},
	)).Methods("GET")
	router.NotFoundHandler = router.NewRoute().HandlerFunc(http.NotFound).GetHandler()

	log.Println(fmt.Sprintf("Server listen on: %s", listen))
	log.Fatal(http.ListenAndServe(listen, router))
}

func healthzHandler(w http.ResponseWriter, r *http.Request) {
	res := map[string]map[string]string{
		"status": map[string]string{
			"message": "success",
			"code":    strconv.Itoa(http.StatusOK),
		},
		"detail": map[string]string{
			"url":  r.RequestURI,
			"help": "checkout http server health",
		}}
	response, _ := json.Marshal(res)
	w.Header().Set("Content-Type", "application/json")
	w.Write(response)
}

func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Do stuff here
		logFormat := fmt.Sprintf("{\"url\":\"%s\", \"method\":\"%s\", \"host\":\"%s\", "+
			"\"client\":\"%s\", \"agent\":\"%s\"}", r.RequestURI, r.Method, r.Host, r.RemoteAddr, r.UserAgent())
		log.Println(logFormat)
		// Call the next handler, which can be another middleware in the chain, or the final handler.
		next.ServeHTTP(w, r)
	})
}

func init() {
	flag.BoolVar(&help, "h", false, "help")
	flag.StringVar(&host, "host", "127.0.0.1", "The hostname or IP on which the REST server will listen")
	flag.IntVar(&port, "port", 8090, "The port on which the REST server will listen")
	flag.StringVar(&path, "config", "./conf.d", "config path directory, defaults is current directory ./conf.d")

	flag.Parse()
	listen = fmt.Sprintf("%s:%d", host, port)
	if help {
		flag.Usage()
		os.Exit(0)
	}

	vper := viper.New()
	vper.SetConfigName("metrics")
	vper.AddConfigPath(path)
	vper.AddConfigPath(".")
	vper.AddConfigPath("/tmp")
	vper.SetConfigType("yaml")

	if err := vper.ReadInConfig(); err != nil {
		panic(fmt.Errorf("Read error config file: %s \n", err))
	}

	if err := vper.Unmarshal(&conf); err != nil {
		panic(fmt.Errorf("Unmarshal error file: %s \n", err))
	}
}
