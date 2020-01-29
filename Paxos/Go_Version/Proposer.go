/*
 * @Author: haha_giraffe
 * @Date: 2020-01-19 09:46:27
 * @Description: 提议者与提议
 */
package main

type Data struct {
	id    int //提议ID
	value int //提议值
}

type Proposer struct {
	m_acceptorCount int //接受者数量
	m_proposerCount int //提议者数量

	m_value                Data //提议者的提议
	m_proposeFinished      bool //完成prepared
	m_isAgree              bool
	m_maxAcceptedSerialNum int
	// m_start
	m_okCount     int
	m_refuseCount int
}

//设置接受者和提议者的数量
func (p *Proposer) SetProposerCount(ac, pc int) {
	p.m_acceptorCount = ac
	p.m_proposerCount = pc
}

//开启propose阶段
func (p *Proposer) StartPropose(value *Data) {
	p.m_value = *value
	p.m_proposeFinished = false
	p.m_isAgree = false
	p.m_maxAcceptedSerialNum = 0
	p.m_okCount = 0
	p.m_refuseCount = 0
}

func (p *Proposer) GetProposal() Data {
	return p.m_value
}

func (p *Proposer) Proposed(ok bool, lastAcceptValue Data) bool {
	if p.m_proposeFinished {
		return true
	}
	if !ok {
		p.m_refuseCount++
		if p.m_refuseCount > p.m_acceptorCount/2 {
			p.m_value.id += p.m_proposerCount
			p.StartPropose(&p.m_value)
			return false
		}

		return true
	}

	p.m_okCount++

	if lastAcceptValue.id > p.m_maxAcceptedSerialNum {
		p.m_maxAcceptedSerialNum = lastAcceptValue.id
		p.m_value.value = lastAcceptValue.value
	}

	if p.m_okCount > p.m_acceptorCount/2 {
		p.m_okCount = 0
		p.m_proposeFinished = true
	}
	return true
}

func (p *Proposer) StartAccept() bool {
	return p.m_proposeFinished
}

func (p *Proposer) Accepted(ok bool) bool {
	if !p.m_proposeFinished {
		return true
	}
	if !ok {
		p.m_refuseCount++
		if p.m_refuseCount > p.m_acceptorCount/2 {
			p.m_value.id += p.m_proposerCount
			p.StartPropose(&p.m_value)
			return false
		}

		return true
	}
	p.m_okCount++
	if p.m_okCount > p.m_acceptorCount/2 {
		p.m_isAgree = true
	}
	return true
}

func (p *Proposer) IsAgree() bool {
	return p.m_isAgree
}
