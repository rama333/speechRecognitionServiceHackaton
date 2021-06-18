package api

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"net/http"
	"path/filepath"
	"speechRecognitionService/decoder"
	"speechRecognitionService/kaldi_go"
)

//TODO not yet required
func Api(kaldi *kaldi_go.Encoder)  {
	router := gin.Default()

	router.MaxMultipartMemory = 8 << 20 // 8 MiB

	//router.Use(CORSMiddleware())
	router.LoadHTMLGlob("../../public/*")
	router.Static("/assets", "../../assets")

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
