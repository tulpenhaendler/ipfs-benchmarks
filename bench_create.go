package main

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/sirupsen/logrus"
	"sync"
)

func (b *Bench) Run(){
	b.l.Log(logrus.InfoLevel, "Starting benchmark")
	b.l.
		WithField("count_regions",b.c.GetNumRegions()).
		WithField("count_node",b.c.getNumNodes()).
		Log(logrus.InfoLevel,"Start config")
	b.makeKeyfiles()
	b.makeSecurityGroups()
}

func (b *Bench) makeKeyfiles(){
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
			input := ec2.CreateKeyPairInput{
				KeyName: aws.String(name),
			}
			keyz,e := client.CreateKeyPair(&input)
			if e != nil {
				log.Error(e)
			}
			kl.Lock()
			defer kl.Unlock()
			b.keys[name] = keyz
			wg.Done()
			log.Info("create key success")

		}(a.Region,a.Name,kl)
	}

	wg.Wait()
}

func (b *Bench) makeSecurityGroups(){
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
			result, err := client.DescribeVpcs(nil)
			if err != nil {
				log.Error("Unable to describe VPCs, %v", err)
			}
			if len(result.Vpcs) == 0 {
				log.Error("No VPCs found to associate security group with.")
			}
			vpcid := aws.StringValue(result.Vpcs[0].VpcId)

			input := ec2.CreateSecurityGroupInput{
				GroupName: aws.String(name),
				Description: aws.String("ipfs-benchark default group"),
				VpcId: &vpcid,
			}
			o,e := client.CreateSecurityGroup(&input)
			if e != nil {
				log.Error(e)
			}
			_, err = client.AuthorizeSecurityGroupIngress(&ec2.AuthorizeSecurityGroupIngressInput{
				GroupId: o.GroupId,
				IpPermissions: []*ec2.IpPermission{
					// Can use setters to simplify seting multiple values without the
					// needing to use aws.String or associated helper utilities.
					(&ec2.IpPermission{}).
						SetIpProtocol("tcp").
						SetFromPort(22).
						SetToPort(22).
						SetIpRanges([]*ec2.IpRange{
							{CidrIp: aws.String("0.0.0.0/0")},
						}),
					(&ec2.IpPermission{}).
						SetIpProtocol("tcp").
						SetFromPort(4001).
						SetToPort(4001).
						SetIpRanges([]*ec2.IpRange{
							(&ec2.IpRange{}).
								SetCidrIp("0.0.0.0/0"),
						}),
					(&ec2.IpPermission{}).
						SetIpProtocol("udp").
						SetFromPort(4001).
						SetToPort(4001).
						SetIpRanges([]*ec2.IpRange{
							(&ec2.IpRange{}).
								SetCidrIp("0.0.0.0/0"),
						}),
				},
			})
			wg.Done()
			log.Info("create SG success")

		}(a.Region,a.Name,kl)
	}

	wg.Wait()
}


func (b *Bench) makeKeyfiles(){
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
			input := ec2.CreateKeyPairInput{
				KeyName: aws.String(name),
			}
			keyz,e := client.CreateKeyPair(&input)
			if e != nil {
				log.Error(e)
			}
			kl.Lock()
			defer kl.Unlock()
			b.keys[name] = keyz
			wg.Done()
			log.Info("create key success")

		}(a.Region,a.Name,kl)
	}

	wg.Wait()
}
