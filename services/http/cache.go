package http

import (
	"fmt"
	"strconv"
	"strings"
	"sync"
	"time"

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
	currentTime := time.Now()
	start := currentTime.Add(time.Duration(-10) * time.Minute)
	outdated := isOutDated(start, proxy.Timestamp)
	if ok && !outdated {
		return proxy, nil
	}

	db := c.db
	table := c.tableName
	c.Delete(name)

	results, err := db.Query("SELECT DISTINCT source, endpoint, port, proxyType, user, pass, apiEndpoint FROM "+table+" WHERE source = ?", name)
	if err != nil {
		fmt.Println("Getting from database failed")
	}
	if results.Next() {
		err = results.Scan(&proxy.ProxyName, &proxy.Endpoint, &proxy.Port, &proxy.ProxyType, &proxy.User, &proxy.Pass, &proxy.APIEndpoint)
		if err != nil {
			fmt.Println("Scanning from rows failed")
		}
		c.Update(proxy.ProxyName, proxy)
	}
	if proxy.ProxyType == "dynamic" {
		ip, err := utils.GetIPFromAPI(proxy.APIEndpoint)
		if err != nil {
			fmt.Println("Getting IP from external api failed")
		}
		ipStr := strings.Split(ip, ":")
		proxy.Endpoint = ipStr[0]
		port, err := strconv.Atoi(ipStr[1])
		if err != nil {
			fmt.Println("Invalid Port")
		}
		proxy.Port = port
		currentTime = time.Now()
		proxy.Timestamp = currentTime
		c.Update(proxy.ProxyName, proxy)
	}
	return proxy, nil
}

func (c *CacheImpl) Delete(name string) {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	delete(c.proxyMap, name)
}

func (c *CacheImpl) Update(name string, proxy Proxy) {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	c.proxyMap[name] = proxy
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

func isOutDated(start, check time.Time) bool {
	return check.Before(start)
}
