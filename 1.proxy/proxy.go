package main

import (
	"fmt"
	"gopkg.in/yaml.v3"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"sync"
)

// structs to parse the yaml

type Proxy struct {
	Configs Configs `yaml:"proxy"`
}

type Configs struct {
	Listen   Host      `yaml:"listen"`
	Services []Service `yaml:"services"`
}

type Host struct {
	Address string `yaml:"address"`
	Port    int    `yaml:"port"`
}

type Service struct {
	Name   string `yaml:"name"`
	Domain string `yaml:"domain"`
	Hosts  []Host `yaml:"hosts"`
}

// structs to manage concurrent access to global variables by server threads

type RequestCache struct {
	mu         sync.Mutex
	statusCode map[string]int
	body       map[string][]byte
}

type RequestsReceived struct {
	mu     sync.Mutex
	amount int
}

// initializes cache maps
func initializeCache() {
	requestCache.statusCode = make(map[string]int)
	requestCache.body = make(map[string][]byte)
}

// reads the yaml file and returns the configs
func readConfigs() (listen Proxy) {

	file, err1 := ioutil.ReadFile("proxy.yaml")
	if err1 != nil {
		fmt.Println("Readfile error:", err1)
	}

	listen = Proxy{}
	err := yaml.Unmarshal(file, &listen)
	if err != nil {
		fmt.Println("Unmarshal error:", err)
	}

	return
}

// returns the list of hosts of a given service
func getServiceHosts(services []Service, serviceId string) []Host {

	for _, service := range services {
		if service.Name == serviceId {
			return service.Hosts
		}
	}

	return nil
}

// checks for a reply to the URI in the cache
func checkCache(reqUri string) (statusCode int, contents []byte) {
	requestCache.mu.Lock()
	defer requestCache.mu.Unlock()
	statusCode = requestCache.statusCode[reqUri]
	contents = requestCache.body[reqUri]
	return
}

// updates the cache for a given request URI
func updateCache(reqUri string, statusCode int, content []byte) {
	requestCache.mu.Lock()
	defer requestCache.mu.Unlock()
	requestCache.statusCode[reqUri] = statusCode
	requestCache.body[reqUri] = content
}

// responds to the client
func respondToClient(w http.ResponseWriter, statusCode int, body []byte) {
	w.WriteHeader(statusCode)
	w.Write(body)
}

// redirects the request to the desired host running the downstream service
func redirectRequest(r *http.Request, host Host) (res *http.Response) {

	reqUrl := fmt.Sprintf("http://%s:%d/", host.Address, host.Port)

	res, err := http.Get(reqUrl)
	if err != nil {
		fmt.Println("error redirecting request to downstream service:", err)
		os.Exit(1)
	}

	return
}

func selectHostRoundRobin(hosts []Host) (host Host) {
	requestsReceived.mu.Lock()
	defer requestsReceived.mu.Unlock()
	host = hosts[requestsReceived.amount%len(hosts)]
	requestsReceived.amount++
	return host
}

func handleRequest(w http.ResponseWriter, r *http.Request) {

	hosts := getServiceHosts(config.Configs.Services, r.Host)

	if hosts == nil {
		respondToClient(w, http.StatusNotFound, []byte("No such service"))
		return
	}

	// check for contents in cache
	cacheKey := r.RequestURI + r.Host
	statusCode, contents := checkCache(cacheKey)
	if contents != nil {
		respondToClient(w, statusCode, contents)
		return
	}

	host := selectHostRoundRobin(hosts)

	// redirect request to the selected host
	resp := redirectRequest(r, host)

	// produce final response
	defer resp.Body.Close()
	body, _ := io.ReadAll(resp.Body)
	updateCache(cacheKey, resp.StatusCode, body)
	respondToClient(w, resp.StatusCode, body)

}

var requestCache = RequestCache{
	statusCode: map[string]int{},
	body:       map[string][]byte{},
}

var requestsReceived = RequestsReceived{
	amount: 0,
}

var config Proxy

func main() {

	initializeCache()

	config = readConfigs()
	listen := config.Configs.Listen

	// register generic handler function
	http.HandleFunc("/", handleRequest)

	// start server
	err := http.ListenAndServe(fmt.Sprintf("%s:%d", listen.Address, listen.Port), nil)

	// server finished with an error warning
	if err != nil {
		fmt.Println("server finished with error:", err)
		os.Exit(1)
	}

}
