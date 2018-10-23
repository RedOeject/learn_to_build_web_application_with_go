/*
基于内存的session存储接口的实现
*/
package provider

import (
	"container/list"
	"github.com/learn_to_build_web_application_with_go/session"
	"sync"
	"time"
)

var provider = &Provider{list: list.New()}

type SessionStore struct {
	sid          string                      //session id
	timeAccessed time.Time                   //最后访问时间
	value        map[interface{}]interface{} //值
}

//实现Session接口
func (st *SessionStore) Set(key, value interface{}) error {
	st.value[key] = value
	provider.SessionUpdate(st.sid)
	return nil
}

func (st *SessionStore) Get(key interface{}) interface{} {

	provider.SessionUpdate(st.sid)
	if v, ok := st.value[key]; ok {
		return v
	} else {
		return nil
	}
}

func (st *SessionStore) Delete(key interface{}) error {
	delete(st.value, key)
	provider.SessionUpdate(st.sid)
	return nil
}

func (st *SessionStore) SessionID() string {
	return st.sid
}

type Provider struct {
	lock     sync.Mutex               //锁
	sessions map[string]*list.Element //用于存储 用Element可以直接加入GC链表中
	list     *list.List               //GC 链表
}

//实现Provider接口
func (provider *Provider) SessionInit(sid string) (session.Session, error) {
	provider.lock.Lock()
	defer provider.lock.Unlock()

	v := make(map[interface{}]interface{}, 0)
	newsess := &SessionStore{sid: sid, timeAccessed: time.Now(), value: v}
	//element:{Value:newsess}
	element := provider.list.PushBack(newsess)
	provider.sessions[sid] = element
	return newsess, nil
}

func (provider *Provider) SessionRead(sid string) (session.Session, error) {
	if element, ok := provider.sessions[sid]; ok {
		return element.Value.(*SessionStore), nil
	} else {
		//获取不到session，初始化一个
		sess, err := provider.SessionInit(sid)
		return sess, err
	}
}

func (provider *Provider) SessionDestroy(sid string) error {
	if element, ok := provider.sessions[sid]; ok {
		delete(provider.sessions, sid)
		provider.list.Remove(element)
		return nil
	}
	return nil
}

func (provider *Provider) SessionGC(maxLifeTime int64) {
	provider.lock.Lock()
	defer provider.lock.Unlock()

	for {
		//如果链表空返回nil，如果非空返回最后一个element
		element := provider.list.Back()
		if element == nil {
			break
		}
		if element.Value.(*SessionStore).timeAccessed.Unix()+maxLifeTime < time.Now().Unix() {
			//超时回收
			provider.list.Remove(element)
			delete(provider.sessions, element.Value.(*SessionStore).sid)
		} else {
			//因为GC表是有序的，所以检测到未超时的Session可以中断
			break
		}
	}
}

func (provider *Provider) SessionUpdate(sid string) error {
	provider.lock.Lock()
	defer provider.lock.Unlock()

	if element, ok := provider.sessions[sid]; ok {
		element.Value.(*SessionStore).timeAccessed = time.Now()
		provider.list.MoveToFront(element)
		return nil
	}
	return nil
}

func init() {
	provider.sessions = make(map[string]*list.Element, 0)
	session.Register("memory", provider)
}
