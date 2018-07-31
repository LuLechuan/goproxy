package http

import (
	"fmt"
	"strings"
	"sync"

	"github.com/jmoiron/sqlx"
)

type Pattern struct {
	Patt      string
	ProxyName string
	Priority  int
}

func NewPattern(pattern string, proxyName string, priority int) *Pattern {
	return &Pattern{
		Patt:      pattern,
		ProxyName: proxyName,
		Priority:  priority,
	}
}

type PatternTable interface {
	Get(url string) (string, bool)
	Delete(patternString string)
}

type PatternTableImpl struct {
	mutex     sync.Mutex
	db        *sqlx.DB
	tableName string
}

func NewPatternTable(sqlConn string, tableName string) (*PatternTableImpl, error) {
	// assume that tableName is valid
	fmt.Println(sqlConn)
	var db *sqlx.DB
	db, err := sqlx.Connect("mysql", sqlConn)
	if db == nil {
		return nil, err
	}

	return &PatternTableImpl{
		db:        db,
		tableName: tableName,
		mutex:     sync.Mutex{},
	}, nil
}

func (pt *PatternTableImpl) Get(url string) (string, bool) {
	results, err := pt.db.Query("SELECT DISTINCT * FROM " + pt.tableName)
	var tempPattern *Pattern
	resultPattern := NewPattern("", "", 0)
	var max int = 0
	if err != nil {
		fmt.Println("Getting from database failed")
	}
	if results.Next() {
		err = results.Scan(&tempPattern.Patt, &tempPattern.ProxyName, &tempPattern.Priority)
		if err != nil {
			fmt.Println("Scanning from rows failed")
		}
		if strings.Contains(url, tempPattern.Patt) && tempPattern.Priority > max {
			resultPattern = tempPattern
			max = tempPattern.Priority
		}
	}
	if resultPattern.ProxyName == "" {
		return "", false
	}
	return resultPattern.ProxyName, true
}

func (pt *PatternTableImpl) Delete(patternString string) {
	_, err := pt.db.NamedExec("DELETE FROM "+pt.tableName+" WHERE pattern = :pattern", map[string]interface{}{
		"pattern": patternString,
	})
	if err != nil {
		fmt.Println("Delete pattern failed")
	}
}
