package registry

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

// RegisterService 客户端调用 注册客户端
func RegisterService(r Registration) error {
	buf := new(bytes.Buffer)
	enc := json.NewEncoder(buf)
	err := enc.Encode(r)
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
