package main

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"time"
)

func asd(){
	l := logrus.New().WithField("test","true")
	n := NewIpfs("localhost:550" +
		"01","testnode",l)

	n.Init()
	t := time.Now()
	id := n.Addr()
	fmt.Println(time.Since(t))
	fmt.Println(id)

	for {
		fmt.Println(n.MakeRandomObject(10))
	}
}