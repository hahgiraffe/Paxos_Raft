/*
 * @Author: haha_giraffe
 * @Date: 2020-01-19 09:45:02
 * @Description: 接受者
 */
package main

type Acceptor struct {
	m_lastAcceptValue Data //最后接受的提议
	m_maxSerialID     int  //Propose提交的最大id
}

func (a *Acceptor) accept(d *Data) bool {
	if d.id == 0 {
		return false
	}
	if a.m_maxSerialID > d.id {
		return false
	}

	a.m_lastAcceptValue = *d
	return true
}

func (a *Acceptor) propose(numID int, d *Data) bool {
	if numID == 0 {
		return false
	}
	if numID < a.m_maxSerialID {
		return false
	}

	a.m_maxSerialID = numID
	d = &a.m_lastAcceptValue
	return true

}
