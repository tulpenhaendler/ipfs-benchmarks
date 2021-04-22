package main

import (
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/aws/aws-sdk-go/service/ssm"
	"github.com/sirupsen/logrus"
	"os"
	"sync"
	"time"
	"golang.org/x/crypto/ssh"
)

func (b *Bench) Run(){
	b.l.Log(logrus.InfoLevel, "Starting benchmark")
	b.l.
		WithField("count_regions",b.c.GetNumRegions()).
		WithField("count_node",b.c.getNumNodes()).
		Log(logrus.InfoLevel,"Start config")
	b.makeKeyfiles()
	b.makeSecurityGroups()
	b.makeInstances()
	b.installIPFSNode()
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
	regions := b.c.GetRegions()
	wg := sync.WaitGroup{}
	wg.Add(len(regions))

	kl := &sync.Mutex{}
	for _,a := range regions {
		go func(region  string,lock *sync.Mutex) {
			log := b.l.WithField("region",region).WithField("step","make_security_group")
			log.Trace("Start make security group")
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
				GroupName: aws.String("ipfsbench"),
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
			kl.Lock()
			defer kl.Unlock()
			b.sgs[region] = o.GroupId
			wg.Done()
			log.Info("create SG success")

		}(a,kl)
	}

	wg.Wait()
}

func (b *Bench) makeInstances(){
	instances := b.c.Nodes.Instances
	wg := sync.WaitGroup{}
	wg.Add(len(instances))

	kl := &sync.Mutex{}
	for _,a := range instances {
		go func(region, name, instance  string,lock *sync.Mutex) {
			log := b.l.WithField("region",region).WithField("step","instanceup").WithField("instance",name)
			log.Trace("Start make instance")
			session := b.aws.GetRegion(region)
			client := ec2.New(session)
			ssmc := ssm.New(session)

			pi := ssm.GetParameterInput{
				Name: aws.String("/aws/service/ami-amazon-linux-latest/amzn2-ami-hvm-x86_64-gp2"),
			}
			ami,_ := ssmc.GetParameter(&pi)

			// Specify the details of the instance that you want to create.
			runResult, err := client.RunInstances(&ec2.RunInstancesInput{
				ImageId:      ami.Parameter.Value,
				InstanceType: aws.String(instance),
				MinCount:     aws.Int64(1),
				MaxCount:     aws.Int64(1),
				SecurityGroups: []*string{aws.String("ipfsbench")},
				BlockDeviceMappings: []*ec2.BlockDeviceMapping{
					&ec2.BlockDeviceMapping{
						DeviceName: aws.String("/dev/xvda"),
						Ebs: &ec2.EbsBlockDevice{
							DeleteOnTermination: aws.Bool(true),
							VolumeSize: aws.Int64(128),
						},
					},
				},
			})

			if err != nil {
				log.Error("Could not create instance", err)
				os.Exit(1)
				return
			}

			log.Info("Created instance ", *runResult.Instances[0].InstanceId)


			// Add tags to the created instance
			_, errtag := client.CreateTags(&ec2.CreateTagsInput{
				Resources: []*string{runResult.Instances[0].InstanceId},
				Tags: []*ec2.Tag{
					{
						Key:   aws.String("IPFSBENCH"),
						Value: aws.String("TRUE"),
					},
				},
			})
			if errtag != nil {
				log.Error("Could not create tags for instance", runResult.Instances[0].InstanceId, errtag)
				os .Exit(1)
			}

			// get public iIP
			ip := ""
			for {
				dii := ec2.DescribeInstancesInput{
					InstanceIds: aws.StringSlice([]string{*runResult.Instances[0].InstanceId}),
				}
				info,_ := client.DescribeInstances(&dii)
				if info.Reservations[0].Instances[0].PublicIpAddress == nil {
					time.Sleep(500*time.Millisecond)
				} else {
					ip = *info.Reservations[0].Instances[0].PublicIpAddress
					break
				}
			}

			kl.Lock()
			defer kl.Unlock()
			b.instances[name] = ip
			wg.Done()
			log.Info("Instance created with IP ",ip)

		}(a.Region,a.Name,b.c.Nodes.Instance, kl)
	}

	wg.Wait()
}

func (b *Bench) installIPFSNode(){
	instances := b.c.Nodes.Instances
	wg := sync.WaitGroup{}
	wg.Add(len(instances))

	kl := &sync.Mutex{}
	for _,a := range instances {
		go func(region, name, instance  string,lock *sync.Mutex) {
			log := b.l.WithField("region",region).WithField("step","ipfsup").WithField("instance",name)
			log.Trace("Start Install IPFS")

			ip := b.instances[name]
			key := b.keys[name]
			signer,_ := signerFromPem([]byte(*key.KeyMaterial),[]byte{})


			sshConfig := &ssh.ClientConfig{
				User: "ec2-user",
				Auth: []ssh.AuthMethod{
					ssh.PublicKeys(signer),
				},
				HostKeyCallback: ssh.InsecureIgnoreHostKey(),
			}

			sshClient := &SshClient{
				Config: sshConfig,
				Server: fmt.Sprintf("%v:%v", ip, 22),
			}

			log.Info("Wait for ssh access...")
			for {
				time.Sleep(100*time.Second)
				_,e := sshClient.RunCommand("ls /")
				if e != nil {
					break
				}
			}

			log.Info("Got ssh access")

			sshClient.RunCommand("curl https://get.docker.com | bash")
			sshClient.RunCommand("usermod -aG docker ec2-user")



			wg.Done()
			log.Info("IPFS installed and running on: ",ip)

		}(a.Region,a.Name,b.c.Nodes.Instance, kl)
	}

	wg.Wait()
}