package registry

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"net/url"
	"sync"
)

// RegisterService 客户端调用 注册客户端
func RegisterService(r Registration) error {
	parse, err := url.Parse(r.ServiceUpdateURL)
	if err != nil {
		return err
	}
	http.Handle(parse.Path, &serviceUpdateHandler{})
	buf := new(bytes.Buffer)
	enc := json.NewEncoder(buf)
	err = enc.Encode(r)
	if err != nil {
		return err
	}

	resp, err := http.Post(ServicesURL, "application/json", buf)
	if err != nil {
		return err
	}

	//请求错误
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("Failed to register service. Registry service responsed with code %v \n",
			resp.StatusCode)
	}

	return nil
}

type serviceUpdateHandler struct{}

func (s serviceUpdateHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	decoder := json.NewDecoder(r.Body)
	var p patch
	err := decoder.Decode(&p)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	prov.Update(p)
}

// ShutdownService 客户端调用注销服务
func ShutdownService(url string) error {
	request, err := http.NewRequest(http.MethodDelete, ServicesURL, bytes.NewBuffer([]byte(url)))
	if err != nil {
		return err
	}
	request.Header.Add("Content-Type", "text/plain")
	response, err := http.DefaultClient.Do(request)
	if err != nil {
		return err
	}
	if response.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to remove service with code %v", response.StatusCode)
	}
	return nil
}

type providers struct {
	services map[ServiceName][]string
	mutex    *sync.RWMutex
}

var prov = providers{
	services: make(map[ServiceName][]string),
	mutex:    new(sync.RWMutex),
}

func (p *providers) Update(pat patch) {
	p.mutex.RLock()
	defer p.mutex.RUnlock()

	for _, entry := range pat.Added {
		if _, ok := p.services[entry.Name]; !ok {
			p.services[entry.Name] = make([]string, 0)
		}
		p.services[entry.Name] = append(p.services[entry.Name], entry.URL)
	}

	for _, entry := range pat.Removed {
		if providerURLs, ok := p.services[entry.Name]; !ok {
			for i := range providerURLs {
				if providerURLs[i] == entry.URL {
					p.services[entry.Name] = append(providerURLs[:i], providerURLs[i+1:]...)
				}
			}
		}
	}

}

func (p *providers) get(name ServiceName) (string, error) {
	services, ok := p.services[name]
	if !ok {
		return "", fmt.Errorf("no providers available for serive %v", name)
	}
	idx := int(rand.Float32() * float32(len(services)))
	return services[idx], nil
}

func GetProvider(name ServiceName) (string, error) {
	return prov.get(name)
}
