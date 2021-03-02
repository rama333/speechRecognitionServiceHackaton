package main

//TODO while everything is in one place, later I will make architecture

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/grpclog"
	"net"
	"speechRecognitionService/decoder"
	"net/http"
	"path/filepath"
	"speechRecognitionService/kaldi_go"
	pb "speechRecognitionService/proto"
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

func main()  {


	go func() {
		listener, err := net.Listen("tcp", ":5300")

		if err != nil{
			grpclog.Fatal("fatal %v", err)
		}


		opts := []grpc.ServerOption{}

		grpcServer := grpc.NewServer(opts...)

		pb.RegisterReverseServer(grpcServer, &server{})

		grpcServer.Serve(listener)

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

	router := gin.Default()

	router.MaxMultipartMemory = 8 << 20 // 8 MiB

	router.Use(CORSMiddleware())
	router.LoadHTMLGlob("./public/*")
	router.Static("/assets", "./assets")

	router.GET("/", func(c *gin.Context) {

		c.HTML(
			http.StatusOK,
			"index.html",
			gin.H{
				"title": "Home Page",
			})
	})

	router.POST("/upload", func(c *gin.Context) {

		file, err := c.FormFile("voice")

		logrus.Println(file.Filename)
		logrus.Println(file.Header)
		logrus.Println(file.Size)

		if err != nil {
			c.String(http.StatusBadRequest, fmt.Sprintf("get form err: %s", err.Error()))
			return
		}

		filename := filepath.Base(file.Filename)
		if err := c.SaveUploadedFile(file, filename); err != nil {
			c.String(http.StatusBadRequest, fmt.Sprintf("upload file err: %s", err.Error()))
			return
		}

		dec := decoder.NewDecoder()

		err = dec.Dec(filename, filename + ".wav")
		if err != nil{
			logrus.Panic(err)
		}

		text, err  := kaldi.Recognition(filename + ".wav")

		if err != nil {
			logrus.Print(err)
			c.JSON(200, gin.H{
				"result": "ok",
				"data": "error",
			})
			return
		}

		c.JSON(200, gin.H{
			"result": "ok",
			"data": text,
		})
		return
	})
	router.Run(":4080")
}