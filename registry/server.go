package registry

import (
	"encoding/json"
	"log"
	"net/http"
	"sync"
)

const (
	//ServerPort 注册中心地址
	ServerPort  = "3000"
	ServicesURL = "http://localhost:" + ServerPort + "/services"
)

type registry struct {
	registrations []Registration
	mutex         *sync.Mutex
}

func (r *registry) add(reg Registration) error {
	r.mutex.Lock()
	r.registrations = append(r.registrations, reg)
	r.mutex.Unlock()
	return nil
}

var reg = registry{
	registrations: make([]Registration, 0),
	mutex:         new(sync.Mutex),
}

type Service struct{}

//ServeHTTP 注册服务
func (s Service) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	log.Println("Request received")
	switch r.Method {
	case http.MethodPost:
		decoder := json.NewDecoder(r.Body)
		var r Registration
		err := decoder.Decode(&r)
		if err != nil {
			log.Println(err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		log.Printf("Add new service: %v with URL: %s\n", r.ServiceName, r.ServiceURL)
		err = reg.add(r)
		if err != nil {
			log.Println(err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
}
