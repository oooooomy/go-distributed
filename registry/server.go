package registry

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"sync"
)

const (
	//ServerPort 注册中心地址
	ServerPort  = ":3000"
	ServicesURL = "http://localhost" + ServerPort + "/services"
)

type registry struct {
	registrations []Registration
	mutex         *sync.RWMutex
}

var reg = registry{
	registrations: make([]Registration, 0),
	mutex:         new(sync.RWMutex),
}

// 添加服务数据
func (r *registry) add(reg Registration) error {
	r.mutex.Lock()
	r.registrations = append(r.registrations, reg)
	r.mutex.Unlock()
	return nil
}

// 删除服务数据
func (r *registry) remove(url string) error {
	for i := range r.registrations {
		if r.registrations[i].ServiceURL == url {
			r.mutex.Lock()
			r.registrations = append(reg.registrations[:i], r.registrations[i+1:]...)
			r.mutex.Unlock()
			return nil
		}
	}
	return fmt.Errorf("service at URL %s not found", url)
}

func (r *registry) sendRequiredServices(reg Registration) error {
	r.mutex.RLock()
	defer r.mutex.RUnlock()

	var p patch
	for _, serviceReg := range r.registrations {
		for _, serviceRequired := range reg.RequiredServices {
			if serviceReg.ServiceName == serviceRequired {
				p.Added = append(p.Added, patchEntry{
					Name: serviceReg.ServiceName,
					URL:  serviceReg.ServiceURL})
			}
		}
	}

	return r.sendPatch(p, reg.ServiceUpdateURL)
}

func (r *registry) sendPatch(p patch, url string) error {
	marshal, err := json.Marshal(p)
	if err != nil {
		return err
	}
	_, err = http.Post(url, "application/json", bytes.NewBuffer(marshal))
	return err
}

type Service struct{}

//ServeHTTP 注册中心提供的HTTP接口
func (s Service) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	log.Println("Request received")
	switch r.Method {

	//注册服务
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

	//删除服务
	case http.MethodDelete:
		data, err := ioutil.ReadAll(r.Body)
		if err != nil {
			log.Println(err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		url := string(data)
		err = reg.remove(url)
		if err != nil {
			log.Println(err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		log.Printf("Remove service at URL: %s\n", url)

	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
}
