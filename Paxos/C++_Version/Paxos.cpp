#include <stdlib.h>
#include <stdio.h>
#include <time.h>
#include "./src/Acceptor.h"
#include "./src/Proposer.h"
#include <pthread.h>
#include <unistd.h>
#include <atomic>
paxos::Proposer p[5];
paxos::Acceptor a[11];
pthread_mutex_t l[11];

std::atomic<int> finishedCount;
int finalValue = -1;
bool isFinished = false;
int g_start;

void* Proposer(void *id)
{
	paxos::Proposer &proposer = p[(long)id];
	paxos::PROPOSAL value = proposer.GetProposal();
	paxos::PROPOSAL lastValue;


	int acceptorId[11];
	int count = 0;

	int start = time(NULL);
    while ( true )
	{
		value = proposer.GetProposal();//拿到提议
		printf("Proposer%d号开始(Propose阶段):提议=[编号:%d，提议:%d]\n", (long)id, value.serialNum, value.value);
		count = 0;
		int i = 0;
		for (i = 0; i < 11; i++ )
		{
		/*
			* 发送消息到第i个acceptor
			* 经过一定时间达到acceptor，sleep(随机时间)模拟
			* acceptor处理消息，mAcceptors[i].Propose()
			* 回应proposer
			* 经过一定时间proposer收到回应，sleep(随机时间)模拟
			* proposer处理回应mProposer.proposed(ok, lastValue)
		*/
			sleep(1);
            //处理消息
            pthread_mutex_lock(&l[i]);
			bool ok = a[i].Propose(value.serialNum, lastValue);
            pthread_mutex_unlock(&l[i]);
			sleep(1);
			//处理Propose回应
			if ( !proposer.Proposed(ok, lastValue) ) //重新开始Propose阶段
			{
			    sleep(1);//经过随机时间,消息到达Proposer
				break;
			}
			paxos::PROPOSAL curValue = proposer.GetProposal();//拿到提议
			if ( curValue.value != value.value )//acceptor本次回应可能推荐了一个提议
			{
				printf("Proposer%d号修改了提议:提议=[编号:%d，提议:%d]\n", (long)id, curValue.serialNum, curValue.value);
				break;
			}
			acceptorId[count++] = i;//记录愿意投票的acceptor
			if ( proposer.StartAccept() ) break;
		}
		//检查有没有达到Accept开始条件，如果没有表示要重新开始Propose阶段
		if ( !proposer.StartAccept() ) continue;

		//开始Accept阶段
		//发送Accept消息到所有愿意投票的acceptor
		value = proposer.GetProposal();
		printf("Proposer%d号开始(Accept阶段):提议=[编号:%d，提议:%d]\n", (long)id, value.serialNum, value.value);
		for (i = 0; i < count; i++ )
		{
			//发送accept消息到acceptor
			//减少accept阶段等待时间，加快收敛
			sleep(1);
            //处理accept消息
            pthread_mutex_lock(&l[acceptorId[i]]);
			bool ok = a[acceptorId[i]].Accept(value);
            pthread_mutex_unlock(&l[acceptorId[i]]);
			sleep(1);
            //处理accept回应
			if ( !proposer.Accepted(ok) ) //重新开始Propose阶段
			{
				sleep(1);
                break;
			}
			if ( proposer.IsAgree() )//成功批准了提议
			{
				start = time(NULL) - start;
                printf("Proposer%d号的提议被批准,用时%lluMS:最终提议 = [编号:%d，提议:%d]\n", (long)id, start, value.serialNum, value.value);
				if(finalValue == -1) finalValue = value.value;
				else if(finalValue != value.value) finalValue = 0;
				if ( 4 == std::atomic_fetch_add(&finishedCount, 1))
				{
					isFinished = true;
					g_start = time(NULL) - g_start;
                    if(finalValue > 0){
						printf("Paxos完成，用时%lluMS，最终通过提议值为：%d\n", g_start, finalValue);
					}
					else{
						printf("Paxos完成，用时%lluMS，最终结果不一致！\n", g_start);
					}
				}
				return NULL;
			}
		}
	}
	return NULL;
}

//Paxos过程模拟演示程序
int main(int argc, char* argv[])
{
	int i = 0;
	printf("5个Proposer, 11个Acceptor准备进行Paxos\n"
		"每个Proposer独立线程，Acceptor不需要线程\n"
		"Proposer编号从0-10,编号为i的Proposer初始提议编号和提议值是（i+1, i+1）\n"
		"Proposer每次重新提议会将提议编号增加5\n"
		"Proposer被批准后结束线程,其它线程继续投票最终，全部批准相同的值，达成一致。\n");
    printf("Paxos开始\n" );
	paxos::PROPOSAL value;

	for ( i = 0; i < 5; i++ ) 
	{
		p[i].SetPlayerCount(5, 11);
		value.serialNum = value.value = i + 1;
		p[i].StartPropose(value);
	}
    pthread_t t[5];
    for(i = 0; i < 11 ; ++i){
        pthread_mutex_init(&l[i], NULL);
    }
    for(i = 0; i < 5 ; ++i){
        pthread_create(&t[i], NULL, Proposer, (void*)i);
    }
    for(i = 0; i < 5 ; ++i){
        pthread_join(t[i],NULL);
    }
	while(true){
		if(isFinished) break;
        sleep(0.5);
	}
    for(i = 0; i < 11 ; ++i){
        pthread_mutex_destroy(&l[i]);
    }
	return 0;
}
