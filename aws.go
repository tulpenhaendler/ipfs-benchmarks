package main

import (
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/sirupsen/logrus"
	"os"
	"sync"
)


type AWSWrapper struct {
	log *logrus.Entry
	sesions map[string]*session.Session
	sl *sync.Mutex
	ds *session.Session
}


func NewAwsWrapper(c *Config, l *logrus.Entry) *AWSWrapper {
	a := AWSWrapper{}
	a.sl = &sync.Mutex{}
	a.log = l.WithField("source","AWSWrapper")
	a.sesions = map[string]*session.Session{

	}
	// default session
	sess := session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	}))
	a.ds = sess
	return &a
}

func (a *AWSWrapper) GetRegion(name string) *session.Session {
	if val,ok := a.sesions[name]; ok {
		return val
	}
	a.sl.Lock()
	defer a.sl.Unlock()
	sess, err := session.NewSession(&aws.Config{
		Region:      aws.String(name),
		Credentials: credentials.NewSharedCredentials("", "ipfsbench"),
	})
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	a.sesions[name] = sess
	return sess
}

