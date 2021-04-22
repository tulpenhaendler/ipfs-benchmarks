package main

import (
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/sirupsen/logrus"
	"math/rand"
)

type Bench struct {
	aws *AWSWrapper
	c *Config
	l *logrus.Entry
	id string

	// step1: make keyfiles
	names map[string]string
	keys map[string]*ec2.CreateKeyPairOutput
}


func NewBench(l *logrus.Entry, c *Config, aws *AWSWrapper) *Bench {
	b := Bench{}
	b.aws = aws
	b.c = c
	b.id = RandStringRunes(12)
	b.l = l.WithField("source","bench")
	b.keys = map[string]*ec2.CreateKeyPairOutput{}
	return &b
}




func RandStringRunes(n int) string {
	var letterRunes = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")
	b := make([]rune, n)
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}
	return string(b)
}