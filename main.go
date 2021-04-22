package main

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"go.uber.org/dig"
	"io"
	"os"
	"time"
	"github.com/olivere/elastic/v7"
	"gopkg.in/sohlich/elogrus.v7"
)

func main(){
	c := dig.New()
	c.Provide(NewConfig)
	c.Provide(NewBench)
	c.Provide(NewAwsWrapper)
	c.Provide(func() *logrus.Entry{

		log := logrus.New()
		client, err := elastic.NewClient(elastic.SetURL("http://localhost:9200"))
		if err != nil {
			log.Panic(err)
		}
		hook, err := elogrus.NewAsyncElasticHook(client, "localhost", logrus.DebugLevel, "mylog")
		if err != nil {
			log.Panic(err)
		}


		l := logrus.New()
		logFile,e := os.OpenFile(".log",os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		if e != nil {
			fmt.Println(e)
		}
		mw := io.MultiWriter(os.Stdout, logFile)
		l.SetOutput(mw)
		l.AddHook(hook)
		l.SetFormatter(&logrus.JSONFormatter{})
		return l.WithField("starttime",time.Now()).WithField("runId",RandStringRunes(8))
	})


	rootCmd := GetRootCommand(c)
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}