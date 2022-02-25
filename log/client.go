package log

import (
	"bytes"
	"fmt"
	"go-distributed/registry"
	syslog "log"
	"net/http"
)

type clientLogger struct {
	url string
}

func SetClientLogger(serviceURL string, clientService registry.ServiceName) {
	syslog.SetPrefix(fmt.Sprintf("[ %v ] - ", clientService))
	syslog.SetFlags(0)
	syslog.SetOutput(&clientLogger{
		serviceURL,
	})
}

func (cl clientLogger) Write(data []byte) (int, error) {
	buffer := bytes.NewBuffer(data)
	res, err := http.Post(cl.url+"/log", "text/plain", buffer)
	if err != nil {
		return 0, err
	}
	if res.StatusCode != http.StatusOK {
		return 0, fmt.Errorf("failed to send log message with code %v", res.StatusCode)
	}
	return len(data), nil
}
