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

type DBStorage interface {
	AddProcessedFile(file string) (error)
	CheckFileIsProssed(file string) (prossed bool, err error)
}

type Process struct {
	log *logrus.Entry
	wg *sync.WaitGroup
	fileRecognition FileRecognition
	dbStorage DBStorage
}

var (
	totalFile int32
	totalErrors int32
	totalSucces int32
	totalSkipped int32
)

func New(fileRecognition FileRecognition, dbStorage DBStorage) (*Process, error)  {

	//if err
	p := &Process{
		log: logrus.WithField("subsystem", "processing"),
		fileRecognition: fileRecognition,
		dbStorage: dbStorage,
	}

	return p, nil
}

func (p *Process) Processing()  {

	defer func() {
		log := p.log.WithFields(logrus.Fields{
				"total succes": totalSucces,
				"total error": totalErrors,
				"total skipped": totalSkipped,
			})
		if totalErrors == 0{
			log.Info("files processed")
		} else {
			log.Warning("files processed with eroors")
		}
	}()

	p.readDir("/home/temp/")

}

func (p *Process) readDir(dir string)  {

	files, err := ioutil.ReadDir(dir)

	if err != nil{
		logrus.Info(err)
	}

	for _,f := range files {

		if !f.IsDir() {

			atomic.AddInt32(&totalFile, 1)

			prossed, err := p.dbStorage.CheckFileIsProssed(f.Name())

			if err != nil{
				p.log.WithError(err).Error("failed to check file")
			}

			if prossed {
				atomic.AddInt32(&totalSkipped, 1)
				continue
			}

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

				fileName := strings.Split(f.Name(), ".xml")[0]

				dec := decoder.NewDecoder()

				err = dec.Dec(dir + fileName + ".wav", "/home/test/" + fileName + ".wav")
				if err != nil{
					logrus.Panic(err)
				}

				wavInfo, err := os.Stat("/home/test/" + fileName + ".wav")

				l := logrus.WithFields(
					logrus.Fields{
						"name": wavInfo.Name(),
						"size":wavInfo.Size(),
					})

				l.Info()

				// if size file > 30 MiB => coun
				if wavInfo.Size() >= 30<<20  {
					atomic.AddInt32(&totalSkipped, 1)
					continue
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

				//send request to elasticsearch
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

				err = p.dbStorage.AddProcessedFile(f.Name())
				if err != nil {
					p.log.WithError(err).Error("failed add in db_storage")
				}

				log := p.log.WithFields(logrus.Fields{
					"total succes": totalSucces,
					"total error": totalErrors,
					"total skipped": totalSkipped,
				})

				log.Info("file processed")
				}
		} else {
			p.readDir(dir + f.Name() + "/")
		}
	}

}

