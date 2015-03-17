// File: dbHandler.go

package controllers

import (
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

type DBController struct {
	session    *mgo.Session
	url        string
	name       string
	collection string
}

func (c DBController) errorHandler(err error) {
	if err != nil {
		panic(err)
	}
}

// --- DB initializer ---
func (c *DBController) Init(info map[string]string) {
	var err error

	c.url = info["url"]
	c.name = info["name"]
	c.collection = info["collection"]

	if c.session == nil {
		c.session, err = mgo.Dial(c.url)
		c.errorHandler(err)
	}
}

// --- query  ---
func (c *DBController) queryCollection(query func(col *mgo.Collection) error) error {
	obj := c.session.DB(c.name).C(c.collection)

	return query(obj)
}

// --- Insert object ---
func (c *DBController) Insert(obj interface{}) {
	var err error

	err = c.session.DB(c.name).C(c.collection).Insert(obj)
	c.errorHandler(err)
}

// --- Remove object ---
func (c *DBController) Remove(obj interface{}) {
	var err error

	err = c.session.DB(c.name).C(c.collection).Remove(obj)
	c.errorHandler(err)
}

// --- Find object ---
func (c *DBController) Find(q interface{}) (results []interface{}) {
	var err error

	query := func(c *mgo.Collection) error {
		fn := c.Find(q).Skip(DB_SKIP).Limit(DB_LIMIT).All(&results)
		if DB_LIMIT < 0 {
			fn = c.Find(q).Skip(DB_SKIP).All(&results)
		}
		return fn
	}

	search := func() error {
		return c.queryCollection(query)
	}

	err = search()
	c.errorHandler(err)

	return
}

// --- Update ---
func (c *DBController) Update(q interface{}, phrase map[string]interface{}) {
	var err error

	change := bson.M{"$set": phrase}
	err = c.session.DB(c.name).C(c.collection).Update(q, change)
	c.errorHandler(err)
}

// --- Update for Array ---
func (c *DBController) UpdateArray(q interface{}, phrase map[string]interface{}) {
	var err error

	change := bson.M{"$addToSet": phrase}
	err = c.session.DB(c.name).C(c.collection).Update(q, change)
	c.errorHandler(err)
}

// --- Find All ---
func (c *DBController) All() (results []interface{}) {
	var err error

	err = c.session.DB(c.name).C(c.collection).Find(bson.M{}).All(&results)
	c.errorHandler(err)

	return
}

// --- DB Close ---
func (c *DBController) Close() {
	c.session.Close()
}
