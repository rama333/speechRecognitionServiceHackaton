package main

//TODO while everything is in one place, later I will make architecture

import (
	"bufio"
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"net"
	"net/http"
	bolt "speechRecognitionService/db"
	"speechRecognitionService/kaldi_go"
	"speechRecognitionService/processing"
	pb "speechRecognitionService/proto"
	"strings"
	"time"
)

func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}

type server struct {}

func (s server) Do(c context.Context, request *pb.Request) (*pb.Response, error) {

	res := &pb.Response{Message: "ok"}

	return res,nil
}


const sampleRate = 44100
const seconds = 2

func maint()  {


	fmt.Println("Launching server...")

	// Устанавливаем прослушивание порта
	ln, _ := net.Listen("tcp", ":8081")

	// Открываем порт
	conn, _ := ln.Accept()

	// Запускаем цикл
	for {
		// Будем прослушивать все сообщения разделенные \n
		message, _ := bufio.NewReader(conn).ReadString('\n')
		// Распечатываем полученое сообщение
		fmt.Print("Message Received:", string(message))
		// Процесс выборки для полученной строки
		newmessage := strings.ToUpper(message)
		// Отправить новую строку обратно клиенту
		conn.Write([]byte(newmessage + "\n"))
	}

	//portaudio.Initialize()
	//defer portaudio.Terminate()
	//buffer := make([]float32, sampleRate * seconds)
	//stream, err := portaudio.OpenDefaultStream(1, 0, sampleRate, len(buffer), func(in []float32) {
	//	for i := range buffer {
	//		buffer[i] = in[i]
	//		logrus.Info(in[i])
	//	}
	//})
	//chk(err)
	//chk(stream.Start())
	//defer stream.Close()

	//go func() {
	//
	//	for {
	//		logrus.Info(buffer)
	//	}
	//}()

	http.HandleFunc("/audio", func(w http.ResponseWriter, r *http.Request) {
		flusher, ok := w.(http.Flusher)
		if !ok {
			panic("expected http.ResponseWriter to be an http.Flusher")
		}

		w.Header().Set("Connection", "Keep-Alive")
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("X-Content-Type-Options", "nosniff")
		w.Header().Set("Transfer-Encoding", "chunked")
		w.Header().Set("Content-Type", "audio/wave")
		for true {
			//binary.Write(w, binary.BigEndian, )//&buffer)
			//logrus.Info(buffer)
			flusher.Flush() // Trigger "chunked" encoding and send a chunk...
			return
		}
	})

	logrus.Info("start")

	http.ListenAndServe(":8080", nil)
}

func chk(err error) {
	if err != nil {
		panic(err)
	}
}

func main()  {

	var st time.Time

	defer func() {
		if !st.IsZero() {
			logrus.Infof("stopped in %s second, exiting", time.Now().Sub(st))
		}
	}()

	kaldi := kaldi_go.NewConfig(&kaldi_go.Config{
		"/home/ubuntu/kaldi/src/online2bin/online2-wav-nnet3-latgen-faster",
		false,
		3,
		1.0,
		13.0,
		6.0,
		7000,
		"/home/ubuntu/test/kaldi-ru-0.9/exp/tdnn/conf/online.conf",
		"/home/ubuntu/test/kaldi-ru-0.9/data/lang_test_rescore/words.txt",
		"/home/ubuntu/test/kaldi-ru-0.9/exp/tdnn/final.mdl",
		"/home/ubuntu/test/kaldi-ru-0.9/exp/tdnn/graph/HCLG.fst",
	})

	db, err := bolt.NewStrorage()

	if err != nil{
		logrus.WithError(err).Fatal("failed to create bolt_db_storage")
	}

	defer func() {
		err := db.Close()

		if err != nil {
			logrus.WithError(err).Fatal("failed to close bolt_db_storage")
		} else{
			logrus.Info("bolt_db_storage closed")
		}

	}()


	proc, err := processing.New(kaldi, db)

	if err != nil{
		logrus.Panic(err)
	}

	proc.Processing()

	//router := gin.Default()
	//
	//router.MaxMultipartMemory = 8 << 20 // 8 MiB
	//
	//router.Use(CORSMiddleware())
	//router.LoadHTMLGlob("../../public/*")
	//router.Static("/assets", "../../assets")
	//
	//router.GET("/", func(c *gin.Context) {
	//
	//	c.HTML(
	//		http.StatusOK,
	//		"index.html",
	//		gin.H{
	//			"title": "Home Page",
	//		})
	//})
	//
	//router.POST("/upload", func(c *gin.Context) {
	//
	//	file, err := c.FormFile("voice")
	//
	//	logrus.Println(file.Filename)
	//	logrus.Println(file.Header)
	//	logrus.Println(file.Size)
	//
	//	if err != nil {
	//		c.String(http.StatusBadRequest, fmt.Sprintf("get form err: %s", err.Error()))
	//		return
	//	}
	//
	//	filename := filepath.Base(file.Filename)
	//	if err := c.SaveUploadedFile(file, filename); err != nil {
	//		c.String(http.StatusBadRequest, fmt.Sprintf("upload file err: %s", err.Error()))
	//		return
	//	}
	//
	//	dec := decoder.NewDecoder()
	//
	//	err = dec.Dec(filename, filename + ".wav")
	//	if err != nil{
	//		logrus.Panic(err)
	//	}
	//
	//	text, err  := kaldi.Recognition(filename + ".wav")
	//
	//	if err != nil {
	//		logrus.Print(err)
	//		c.JSON(200, gin.H{
	//			"result": "ok",
	//			"data": "error",
	//		})
	//		return
	//	}
	//
	//	c.JSON(200, gin.H{
	//		"result": "ok",
	//		"data": text,
	//	})
	//	return
	//})
	//router.Run(":4080")

	st = time.Now()
}