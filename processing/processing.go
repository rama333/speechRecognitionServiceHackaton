package processing

import (
	"bytes"
	"encoding/json"
	"encoding/xml"
	"github.com/sirupsen/logrus"
	"io/ioutil"
	"net/http"
	"time"

	"os"
	"regexp"
	"speechRecognitionService/decoder"
	"speechRecognitionService/entity"
	"strings"
	"sync"
	"sync/atomic"

)

type FileRecognition interface {
	Recognition(file string) (string, error)
}

type Process struct {
	log *logrus.Entry
	wg *sync.WaitGroup
	fileRecognition FileRecognition
}

var totalFile int32
var totalErrors int32
var totalSucces int32

func New(fileRecognition FileRecognition) (*Process, error)  {

	//if err
	p := &Process{
		log: logrus.WithField("subsystem", "processing"),
		fileRecognition: fileRecognition,
	}

	return p, nil
}

func (p *Process) Processing()  {

	p.readDir("/home/temp/")

	logrus.Info(totalFile)

}

func (p *Process) readDir(dir string)  {

	files, err := ioutil.ReadDir(dir)

	if err != nil{
		logrus.Info(err)
	}

	for _,f := range files{
		if !f.IsDir() {
			//logrus.Info(f.Name())
			atomic.AddInt32(&totalFile, 1)

			math, err := regexp.MatchString(".xml", f.Name())

			if err != nil{
				atomic.AddInt32(&totalErrors, 1)
			}

			if math {

				p.log.Info("start proccesing for ", f.Name())

				file, err := os.Open(dir + f.Name())

				if err != nil{
					atomic.AddInt32(&totalErrors, 1)
				}

				read, err := ioutil.ReadAll(file)

				var rec entity.Recording
				xml.Unmarshal(read, &rec)

				//logrus.Infof("%v", rec)

				fileName := strings.Split(f.Name(), ".xml")[0]

				p.log.Info(fileName)

				dec := decoder.NewDecoder()

				err = dec.Dec(dir + fileName + ".wav", "/home/test/" + fileName + ".wav")
				if err != nil{
					logrus.Panic(err)
				}

				recogn, err := p.fileRecognition.Recognition("/home/test/" + fileName + ".wav")

				p.log.Info(recogn)
				atomic.AddInt32(&totalSucces, 1)

				data := entity.Data{recogn, rec}

				j,err := json.Marshal(data)
				if err != nil{
					p.log.Info("error", err)
					atomic.AddInt32(&totalErrors, 1)
				}

				req, err := http.NewRequest("POST","http://192.168.143.91:5088", bytes.NewBuffer(j))

				if err != nil {
					p.log.Info(err)
				}
				req.Header.Set("Content-Type", "application/json")

				client := &http.Client{
					Timeout: time.Second * 10,
				}
				resp, err := client.Do(req)
				if err != nil {
					panic(err)
				}
				defer resp.Body.Close()

			}


		} else {
			p.readDir(dir + f.Name() + "/")
		}
	}

}

