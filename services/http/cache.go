package http

import (
	"fmt"
	"sync"

	"github.com/LuLechuan/goproxy/utils"
	"github.com/jmoiron/sqlx"
)

type Cache interface {
	GetProxy(name string) (Proxy, error)
	Delete(name string)
	Update(name string, ip string)
}

type CacheImpl struct {
	mutex     sync.Mutex
	proxyMap  map[string]Proxy
	db        *sqlx.DB
	tableName string
}

func (c *CacheImpl) GetProxy(name string) (Proxy, error) {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	proxy, ok := c.proxyMap[name]
	// TODO: check the timestamp of the proxy
	if ok {
		return proxy, nil
	}

	db := c.db
	table := c.tableName
	c.Delete(name)

	results, err := db.Query("SELECT DISTINCT * FROM "+table+" WHERE name = ?", name)
	if err != nil {
		fmt.Println("Getting from database failed")
	}
	if results.Next() {
		err = results.Scan(&proxy.ProxyName, &proxy.IP, &proxy.IsDynamic, &proxy.User, &proxy.Pass, &proxy.Timestamp)
		if err != nil {
			fmt.Println("Scanning from rows failed")
		}
		c.Update(proxy.ProxyName, proxy.IP)
	}
	if proxy.IsDynamic == true {
		// TODO: Remove hard code here
		ip, err := utils.GetIPFromAPI("http://dly.134t.com/query.txt?key=NP6AB48D4F&word=&count=10")
	}
	return proxy, nil
}

func (c *CacheImpl) Delete(name string) {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	delete(c.proxyMap, name)
}

func (c *CacheImpl) Update(name string, ip string) {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	proxy := c.proxyMap[name]
	proxy.SetIP(ip)
}

func NewCache(sqlConn string, tableName string) (*CacheImpl, error) {
	// assume that tableName is valid
	fmt.Println(sqlConn)
	var db *sqlx.DB
	db, err := sqlx.Connect("mysql", sqlConn)
	if db == nil {
		return nil, err
	}

	return &CacheImpl{
		db:        db,
		tableName: tableName,
		mutex:     sync.Mutex{},
		proxyMap:  map[string]Proxy{},
	}, nil
}
