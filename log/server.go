package log

import (
	"io/ioutil"
	syslog "log"
	"net/http"
	"os"
)

var log *syslog.Logger

type fileLog string

func (fl fileLog) Write(data []byte) (int, error) {
	f, err := os.OpenFile(string(fl), os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0600)
	if err != nil {
		return 0, err
	}
	defer func(f *os.File) {
		err := f.Close()
		if err != nil {
			log.Fatal("File close error: ", err)
		}
	}(f)
	return f.Write(data)
}

func Run(destination string) {
	log = syslog.New(fileLog(destination), "[ LogService ]: ", syslog.LstdFlags)
}

func RegisterHandler() {
	http.HandleFunc("/log", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		//http post request
		case http.MethodPost:
			bytes, err := ioutil.ReadAll(r.Body)
			if err != nil || len(bytes) == 0 {
				w.WriteHeader(http.StatusBadRequest)
				return
			}
			write(string(bytes))
		default:
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}
	})
}

func write(message string) {
	log.Printf("%v\n", message)
}
