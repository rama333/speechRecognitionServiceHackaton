package main

//TODO while everything is in one place, later I will make architecture

import (
	"github.com/sirupsen/logrus"
	"speechRecognitionService/api"
	"speechRecognitionService/kaldi_go"
	"time"
)

func main()  {

	var st time.Time

	defer func() {
		if !st.IsZero() {
			logrus.Infof("stopped in %s second, exiting", time.Now().Sub(st))
		}
	}()

	kaldi := kaldi_go.NewConfig(&kaldi_go.Config{
		"/home/ramil/kaldi/src/online2bin/online2-wav-nnet3-latgen-faster",
		false,
		3,
		1.0,
		13.0,
		6.0,
		7000,
		"/home/ramil/kaldi-ru-0.9/exp/tdnn/conf/online.conf",
		"/home/ramil/kaldi-ru-0.9/data/lang_test_rescore/words.txt",
		"/home/ramil/kaldi-ru-0.9/exp/tdnn/final.mdl",
		"/home/ramil/kaldi-ru-0.9/exp/tdnn/graph/HCLG.fst",
	})


	api.Api(kaldi)

	//db, err := bolt.NewStrorage()
	//
	//if err != nil{
	//	logrus.WithError(err).Fatal("failed to create bolt_db_storage")
	//}
	//
	//defer func() {
	//	err := db.Close()
	//
	//	if err != nil {
	//		logrus.WithError(err).Fatal("failed to close bolt_db_storage")
	//	} else {
	//		logrus.Info("bolt_db_storage closed")
	//	}
	//}()

	//proc, err := processing.New(kaldi, db)
	//
	//if err != nil{
	//	logrus.Panic(err)
	//}
	//
	//proc.Processing()

	st = time.Now()
}