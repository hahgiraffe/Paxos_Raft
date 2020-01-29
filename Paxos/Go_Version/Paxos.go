/*
 * @Author: haha_giraffe
 * @Date: 2020-01-19 09:44:10
 * @Description: 模拟Paxos算法
 */
package main

import (
	"fmt"
	"sync"
	"sync/atomic"
	"time"
)

//设置接受者和提议者的人数
const (
	PROPOSERNUM = 5
	ACCEPTNUM   = 11
)

var (
	p           [5]Proposer    //提议者
	a           [11]Acceptor   //接受者
	mu          [11]sync.Mutex //每一个接受者对应一个mu
	wg          sync.WaitGroup //goroutine中同步使用
	isFinished  bool           = false
	finalValue  int
	finishcount int32
)

func mypropose(id int) {
	fmt.Printf("the Proposor %d is beginning\n", id)
	mypro := p[id]
	var value = mypro.GetProposal()
	var lastValue Data

	var acceptorId [11]int
	var count int = 0

	for {
		value = mypro.GetProposal()
		fmt.Printf("Proposer %d 开始Propose阶段：提议 = [编号:%d, 提议:%d]\n", id, value.id, value.value)
		count = 0
		for i := 0; i < ACCEPTNUM; i++ {
			time.Sleep(1)
			mu[i].Lock()
			ok := a[i].propose(value.id, &lastValue)
			mu[i].Unlock()
			time.Sleep(1)
			if !mypro.Proposed(ok, lastValue) {
				time.Sleep(1)
				break
			}
			curValue := mypro.GetProposal()
			if curValue.value != value.value {
				fmt.Printf("Proposer%d号修改了提议:提议=[编号:%d，提议:%d]\n", id, curValue.id, curValue.value)
				break
			}
			acceptorId[count] = i
			count++
			if mypro.StartAccept() {
				break
			}
		}

		if !mypro.StartAccept() {
			continue
		}

		//开始Accept
		value = mypro.GetProposal()
		fmt.Printf("Proposer%d号开始Accept阶段:提议=[编号:%d，提议:%d]\n", id, value.id, value.value)
		for i := 0; i < count; i++ {
			time.Sleep(1)
			mu[i].Lock()
			ok := a[acceptorId[i]].accept(&value)
			mu[i].Unlock()
			time.Sleep(1)
			if !mypro.Accepted(ok) {
				time.Sleep(1)
				break
			}

			if mypro.IsAgree() {
				fmt.Printf("%d号提议被批准, 最终提议=[编号:%d，提议:%d]\n", id, value.id, value.value)
				if finalValue == -1 {
					finalValue = value.value
				} else if finalValue != value.value {
					finalValue = 0
				}

				atomic.AddInt32(&finishcount, 1)
				if PROPOSERNUM == atomic.LoadInt32(&finishcount) {
					isFinished = true
					if finalValue > 0 {
						fmt.Printf("Paxos完成，最终提议值为:%d\n", finalValue)
					} else {
						fmt.Printf("Paxos完成，最终提议值不一致\n")
					}
				}
				wg.Done()
				return
			}
		}
	}
	// return
	// wg.Done()
}

func main() {
	fmt.Println("Paxos开始")
	var d Data
	// atomic.StoreInt32(finishcount, 0)
	for i := 0; i < 5; i++ {
		p[i].SetProposerCount(5, 11)
		d.id = i + 1
		d.value = i + 1
		p[i].StartPropose(&d)
	}
	for num := 0; num < PROPOSERNUM; num++ {
		wg.Add(1)
		go mypropose(num)
	}
	wg.Wait()
	fmt.Println("Paxos结束")

}
