/*
 * @Author: haha_giraffe
 * @Date: 2020-01-17 19:59:47
 * @Description: file content
 */
#ifndef PAXOS_DATA_H
#define PAXOS_DATA_H

namespace paxos
{
	//提议数据结构
	typedef struct PROPOSAL
	{
		unsigned int	serialNum;      //提议ID，1开始递增，保证全局唯一
		unsigned int	value;          //提议内容（提议值）
	}PROPOSAL;
}

#endif //PAXOS_DATA_H