package main

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"go.uber.org/dig"
	"os"
	"time"
)

func main(){
	c := dig.New()
	c.Provide(NewConfig)
	c.Provide(NewBench)
	c.Provide(NewAwsWrapper)
	c.Provide(func() *logrus.Entry{
		l := logrus.New()
		l.SetFormatter(&logrus.JSONFormatter{})
		return l.WithField("starttime",time.Now())
	})


	rootCmd := GetRootCommand(c)
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}