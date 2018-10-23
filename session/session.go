package session

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"net/http"
	"net/url"
	"sync"
	"time"
)

/**
目前Go标准包没有为session提供任何支持，我们将会自己动手来实现go版本的session管理和创建。
session的基本原理是由服务器为每个会话维护一份信息数据，
客户端和服务端依靠一个全局唯一的标识来访问这份数据，以达到交互的目的。
当用户访问Web应用时，服务端程序会随需要创建session，这个过程可以概括为三个步骤：
1.生成全局唯一标识符（sessionid）；
2.开辟数据存储空间。一般会在内存中创建相应的数据结构，但这种情况下，系统一旦掉电，所有的会话数据就会丢失，如果是电子商务类网站，这将造成严重的后果。所以为了解决这类问题，你可以将会话数据写到文件里或存储在数据库中，当然这样会增加I/O开销，但是它可以实现某种程度的session持久化，也更有利于session的共享；
3.将session的全局唯一标示符发送给客户端。

如何发送这个session的唯一标识这一步上,一般来说会有两种常用的方式：cookie和URL重写。
1.Cookie 服务端通过设置Set-cookie头就可以将session的标识符传送到客户端，而客户端此后的每一次请求都会带上这个标识符，另外一般包含session信息的cookie会将失效时间设置为0(会话cookie)，即浏览器进程有效时间。至于浏览器怎么处理这个0，每个浏览器都有自己的方案，但差别都不会太大(一般体现在新建浏览器窗口的时候)；
2.URL重写 所谓URL重写，就是在返回给用户的页面里的所有的URL后面追加session标识符，这样用户在收到响应之后，无论点击响应页面里的哪个链接或提交表单，都会自动带上session标识符，从而就实现了会话的保持。虽然这种做法比较麻烦，但是，如果客户端禁用了cookie的话，此种方案将会是首选。

- 全局session管理器
- 保证sessionid 的全局唯一性
- 为每个客户关联一个session
- session 的存储(可以存储到内存、文件、数据库等)
- session 过期处理
*/

//定义一个全局的session管理器
type Manager struct {
	cookieName  string     //保存sessionid的cookie名
	lock        sync.Mutex //lock
	provider    Provider   //Session服务提供者
	maxLifeTime int64      //保存时间
}

func NewManger(provideName, cookieName string, maxLifeTime int64) (*Manager, error) {
	provider, ok := providers[provideName]
	if !ok {
		return nil, fmt.Errorf("session: unknown provide %q ", provideName)
	}
	return &Manager{provider: provider, cookieName: cookieName, maxLifeTime: maxLifeTime}, nil
}

//定义一个Provider接口，处理session管理器底层存储结构，提供以下四个服务，底层存储实现这个接口
type Provider interface {
	//SessionInit函数实现Session的初始化，操作成功则返回此新的Session变量
	SessionInit(sid string) (Session, error)
	//SessionRead函数返回sid所代表的Session变量，如果不存在，那么将以sid为参数调用SessionInit函数创建并返回一个新的Session变量
	SessionRead(sid string) (Session, error)
	//SessionDestroy函数用来销毁sid对应的Session变量
	SessionDestroy(sid string) error
	//SessionGC根据maxLifeTime来删除过期的数据
	SessionGC(maxLifeTime int64)
}

var providers = make(map[string]Provider)

//Session接口四个操作，Session服务
type Session interface {
	Set(key, value interface{}) error
	Get(key interface{}) interface{}
	Delete(key interface{}) error
	SessionID() string
}

//注册服务提供者
func Register(name string, provider Provider) {
	if provider == nil {
		panic("session: Register provider is nil")
	}
	if _, dup := providers[name]; dup {
		panic("session: Register called twice for provider " + name)
	}

	providers[name] = provider
}

//全局唯一的Session ID
func (manager *Manager) sessionId() string {
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		return ""
	}

	return base64.URLEncoding.EncodeToString(b)
}

//创建Session
func (manager *Manager) SessionStart(w http.ResponseWriter, r *http.Request) (session Session) {
	manager.lock.Lock()
	defer manager.lock.Unlock()
	cookie, err := r.Cookie(manager.cookieName)

	if err != nil || cookie.Value == "" {
		//创建一个Session
		sid := manager.sessionId()
		session, _ = manager.provider.SessionInit(sid)
		//QueryEscape转义字符串，以便可以安全地放在URL查询中
		cookie := http.Cookie{Name: manager.cookieName, Value: url.QueryEscape(sid), Path: "/", HttpOnly: true, MaxAge: int(manager.maxLifeTime)}
		http.SetCookie(w, &cookie)
	} else {
		//Session已存在，读取
		sid, _ := url.QueryUnescape(cookie.Value)
		session, _ = manager.provider.SessionRead(sid)
	}

	return
}

//session重置
func (manager *Manager) SessionDestroy(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie(manager.cookieName)
	if err != nil || cookie.Value == "" {
		return
	} else {
		manager.lock.Lock()
		defer manager.lock.Unlock()
		manager.provider.SessionDestroy(cookie.Value)
		expiration := time.Now()
		cookie := http.Cookie{Name: manager.cookieName, Path: "/", HttpOnly: true, Expires: expiration, MaxAge: -1}
		http.SetCookie(w, &cookie)
	}
}

//session销毁
func (manager *Manager) GC() {
	manager.lock.Lock()
	defer manager.lock.Unlock()
	manager.provider.SessionGC(manager.maxLifeTime)
	//maxLifeTime纳秒运行一次GC
	time.AfterFunc(time.Duration(manager.maxLifeTime), func() { manager.GC() })
}
