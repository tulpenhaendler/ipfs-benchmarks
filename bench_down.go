package main

import (
	"fmt"
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
	b.terminateInstances()
	
	b.deleteKeyfiles()
	b.deleteSecurityGroups()
}

func (b *Bench) deleteKeyfiles(){
	regions := b.c.GetRegions()
	wg := sync.WaitGroup{}
	wg.Add(len(regions))

	kl := &sync.Mutex{}
	for _,a := range regions {
		go func(region string,lock *sync.Mutex) {
			log := b.l.WithField("region",region).WithField("step","keygen")
			log.Trace("Start keygen")
			session := b.aws.GetRegion(region)
			client := ec2.New(session)
			input := ec2.DeleteKeyPairInput{
				KeyName: aws.String("ipfsbench"),
			}
			_,e := client.DeleteKeyPair(&input)
			if e != nil {
				log.Error(e)
			}
			wg.Done()
			log.Info("delete key success")

		}(a,kl)
	}

	wg.Wait()
}

func (b *Bench) deleteSecurityGroups(){
	regions := b.c.GetRegions()
	wg := sync.WaitGroup{}
	wg.Add(len(regions))

	kl := &sync.Mutex{}
	for _,a := range regions {
		go func(region string,lock *sync.Mutex) {
			log := b.l.WithField("region",region).WithField("step","make_security_group")
			log.Trace("Start sg")
			session := b.aws.GetRegion(region)
			client := ec2.New(session)
			input := ec2.DeleteSecurityGroupInput{
				GroupName: aws.String("ipfsbench"),
			}
			_,e := client.DeleteSecurityGroup(&input)
			if e != nil {
				log.Error(e)
			}

			wg.Done()
			log.Info("delete SG success")

		}(a,kl)
	}

	wg.Wait()
}


func (b *Bench) terminateInstances(){
	regions := b.c.GetRegions()
	wg := sync.WaitGroup{}
	wg.Add(len(regions))

	kl := &sync.Mutex{}
	for _,a := range regions {
		go func(region  string,lock *sync.Mutex) {
			log := b.l.WithField("region",region).WithField("step","instanceup")
			log.Trace("Start keygen")
			session := b.aws.GetRegion(region)
			client := ec2.New(session)

			di := ec2.DescribeInstancesInput{
				Filters:     []*ec2.Filter{
					&ec2.Filter{
						Name: aws.String("tag:IPFSBENCH"),
						Values: aws.StringSlice([]string{"TRUE"}),
					},
				},
			}

			instances,_ := client.DescribeInstances(&di)

			for _,i := range instances.Reservations {
				fmt.Println(*i.Instances[0].InstanceId)
				_,e := client.TerminateInstances(&ec2.TerminateInstancesInput{
					InstanceIds: aws.StringSlice([]string{*i.Instances[0].InstanceId}),
				})
				log.Info(e)
			}

			wg.Done()
			log.Info("Instances terminated")

		}(a, kl)
	}

	wg.Wait()
}