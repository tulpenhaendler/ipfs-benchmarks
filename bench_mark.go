package main

import (
	"math/rand"
	"time"
)

func (b *Bench) DoBench(count int){
	// when this gets executed, all nodes are already provisioned
	// to specify where ndoes should be take a look at the config.yml
	// every 5 seconds, we pick a random node
	// give it 1MB of data, and see how long it takes all others
	// to resolve the CID of it
	interval := 5*time.Second
	size := 10
	for i:=0;i<=count;i++{ // do "count" benchmarks
		numnodes := len(b.nodes)
		source := rand.Int() % numnodes
		var sourceNode *IPFS
		others := []*IPFS{}
		i :=0
		// pick a random node as source
		for _,a := range b.nodes {
			if i == source {
				sourceNode = a
			} else {
				others = append(others, a)
			}
			i++
		}
		// on the source node, make a random datablock, and get the CID
		cid := sourceNode.MakeRandomObject(size)
		b.l.Info("cid generated: ", cid)
		b.countsLock.Lock()
		// to keepl track of what is synced
		b.counts[cid] = 0
		b.countsLock.Unlock()
		for i,_ := range others {
			go func(i int,cid string,l int) {
				// wait till that node can resolve that cid
				for {
					if v := others[i].CanGetCid(cid); v == true {
						break
					}
				}
					// got it at this point
					b.countsLock.Lock()
					b.counts[cid] += 1
					b.countsLock.Unlock()

			}(i,cid,len(others))
		}

		go func(cid string,l int) {
			for {
				// got it at this point
				b.countsLock.Lock()
				b.l.WithField("source","bench").WithField("cid",cid).
					WithField("targetCount",l).WithField("actualCount",b.counts[cid]).Info("Count update")
				if l == b.counts[cid] {
					b.l.WithField("source","bench").WithField("cid",cid).
						WithField("targetCount",l).WithField("actualCount",b.counts[cid]).Info("Target count reached!!!")
					b.countsLock.Unlock()
					break
				}
				b.countsLock.Unlock()
				time.Sleep(1*time.Second)
			}
		}(cid,len(others))


		time.Sleep(interval)
	}

}

