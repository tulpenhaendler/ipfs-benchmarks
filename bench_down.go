package main

import (
"github.com/aws/aws-sdk-go/aws"
"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/sirupsen/logrus"
	"sync"
)

func (b *Bench) Delete(){
	b.l.Log(logrus.InfoLevel, "Starting benchmark")
	b.l.
		WithField("count_regions",b.c.GetNumRegions()).
		WithField("count_node",b.c.getNumNodes()).
		Log(logrus.InfoLevel,"Start config")
	b.deleteKeyfiles()
	b.deleteSecurityGroups()
}

func (b *Bench) deleteKeyfiles(){
	instances := b.c.Nodes.Instances
	wg := sync.WaitGroup{}
	wg.Add(len(instances))

	kl := &sync.Mutex{}
	for _,a := range instances {
		go func(region, name  string,lock *sync.Mutex) {
			log := b.l.WithField("region",region).WithField("step","keygen")
			log.Trace("Start keygen")
			session := b.aws.GetRegion(region)
			client := ec2.New(session)
			input := ec2.DeleteKeyPairInput{
				KeyName: aws.String(name),
			}
			_,e := client.DeleteKeyPair(&input)
			if e != nil {
				log.Error(e)
			}
			wg.Done()
			log.Info("delete key success")

		}(a.Region,a.Name,kl)
	}

	wg.Wait()
}

func (b *Bench) deleteSecurityGroups(){
	instances := b.c.Nodes.Instances
	wg := sync.WaitGroup{}
	wg.Add(len(instances))

	kl := &sync.Mutex{}
	for _,a := range instances {
		go func(region, name  string,lock *sync.Mutex) {
			log := b.l.WithField("region",region).WithField("step","make_security_group")
			log.Trace("Start keygen")
			session := b.aws.GetRegion(region)
			client := ec2.New(session)
			input := ec2.DeleteSecurityGroupInput{
				GroupName: aws.String(name),
			}
			_,e := client.DeleteSecurityGroup(&input)
			if e != nil {
				log.Error(e)
			}

			wg.Done()
			log.Info("delete SG success")

		}(a.Region,a.Name,kl)
	}

	wg.Wait()
}
