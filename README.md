# GoLive
## Table of contents
* [General info](#general-info)
* [Problem](#problem)
* [Approach](#approach)
* [Cache Package](#cache-package)
* * [Class Diagram](#class-diagram)
* * [Expiry](#expiry)
* * [CacheObject](#cacheobject)
* * [Cache](#cache)
* * [CacheManager](#cachemanager)
* * [Technologies](#technologies)
* [Setup](#setup)
* * [E-commerce](#e-commerce)
* * [Cache Manager](#cache-manager)
* [Further Improvements](#further-improvements)

## General info
This was written for our Capstone Project in GoSch. Cache Manager is a package that stores data for our e-commerce site in memory. This is to improve rendering times and client requests by reducing pings to the database and processes data through a read-through and write-through strategy.

## Problem
We looked at a traditional E-commerce site and identified the following issues with a typical MERN or go-react app.

Pings to the Database is frequent, Database being stored in harddrives makes cache misses the primary contributing factor to read times.
Database logic tends to be tightly coupled with the application layer making troubleshooting problematic as the application code base scales.
Single points of failure, should any of the microservice goes down, the application is unable to continuously serve its function.

## Approach
In order to address the first problem, we decided to store data in memory which has an internal expiry time such that it deletes itself if it isn't accessed again within a specific timeframe. This is a simple and effective solution to the first problem as processing speed in memory is significantly faster. However, as we proceeded to expand on the code base, we noticed the code blocks growing bigger and bigger. This was attributed to the increasing amounts of Database logic that is added into the code on top of the in memory data structures.

We went with a two pronged approach to rectify the issue:

Decouple our database into its own package which deals with all of the queries and returns data types for the application to consume

Decouple the cache manager such that it acts as an effective buffer, exposing functions and methods that the application needs to access the data within without exposing how it functions.

Upon decoupling the cache manager, we noticed copious amounts of code repeated as we had multiple data types to work with. Seeing that each cache exhibited very similar behaviour in that they need to be able to read,write and expire based on duration, we decided to go with dependency injection.

In order to address the 3rd problem, we intended to implement a load balancer and host multiple isntances of the application with multiple databases all interlinked through read and write requests. This in itself would present to be a pretty huge project on its own which is why we decided to focus on the Cache Package.

## Cache Package
### Class Diagram
![image](https://user-images.githubusercontent.com/73516102/118154022-21589d80-b449-11eb-8d3f-b0a821a06108.png)

### Expiry
The expiry type is the building block for all our cached data.
```Go
type expiry struct{
expiry time.Time
}
```
It contains the time in which the data was cached and methods to self regulate. Access to the data would call on the updateExpiryTime method which refreshes the time to the current time.
```Go
//monitor sleeps for the difference between the expiration time and the current time.
//eg expires at 4pm, current time is 4.30pm. it sleeps for 30mins.
//after sleep ends, it checks if the expiration time has exceeded the current time it returns.
//or else it loops.
func (e *expiry) monitor() {
	for {
		sleeptime := time.Until(e.expiry)
		time.Sleep(sleeptime)
		if e.expiry.Before(time.Now()) {
			return
		}
	}
}

//updateExpiryTime updates the session's expiry time.
func (e *expiry) updateExpiryTime(updatedTime time.Time) {
	e.expiry = updatedTime
}
```
When data is added into a cache, the caches self regulate through the tidy function which is fired off in go routines.
```Go
func (c *cache) add(payLoad CacheObject) {
	key := payLoad.getKey()
	if ok := c.check(key); !ok {
		(*c)[key] = payLoad
		fmt.Println((*c)[key], "added into the cache")
		fmt.Println(key, "key value during storage")
		go c.tidy(key, payLoad)
	}
}
```
The tidy function calls on monitor which blocks and eventually deletes the item off the map when the monitor function returns.
```Go
func (c *cache) tidy(key string, session CacheObject) {
	session.monitor()
	delete(*c, key)
}
```

### Cache Object
The Cache Object interface is the behaviour we require data that is stored in our cache to implement.
```Go
type CacheObject interface {
	monitor()
	updateExpiryTime(time.Time)
	getKey()
}
```
This is implemented by the various data types through embeddeding the expiry type and declaration of their own getKey methods as they may have different ways of storing the keys in their internal data types.

### Cache
The cache type then stores each of these cacheobjects in a map and exposes methods to be called upon by the cache manager.
```Go
type cache map[string]CacheObject

//eg. add CacheObject
func (c *cache) add(payLoad CacheObject) {
	key := payLoad.getKey()
	if ok := c.check(key); !ok {
		(*c)[key] = payLoad
		fmt.Println((*c)[key], "added into the cache")
		fmt.Println(key, "key value during storage")
		go c.tidy(key, payLoad)
	}
}
```
### Cache Manager
The Cache Manager is essentially a wrapper struct that stores each of the caches in its fields and routes to the respective cache depending on their type. By storing a pointer to the db, the cache manager is able to interact directly with the DB, abstracting the whole implementation from the app.
```Go
type CacheManager struct {
	activeUserCache activeUserCache 
	itemsCache      itemCache       
	database        *db.Database    
}

func (c *CacheManager) AddtoCache(payLoad CacheObject) {
	switch v := payLoad.(type) {
	case ActiveSession:
		c.activeUserCache.add(v)
	case *cachedItem:
		c.itemsCache.add(v)

	}
}
```
## Technologies
Project is created with:

* Bootstrap
* BCrypt
* Go
* Gorilla Mux
* HTML

## Setup
### E-commerce
enter go run . in the root folder

### Cache Manager
This hasn't been fully modularised to adapt all data types. Should anyone wish to use, you would have to add the package to your project and supply your fundamental types.

## Afterthought
The cache package gives us ways to instantiate data types that we require to store within the Cache Manager. These data types then implement the CacheObject interface allowing the majority of the methods and functions to be reused within the package. To a huge extent we have decouple our packages but the fundamental types of the data still remains tightly coupled which isn't a bad thing as our primary objectives to achiveve both faster read speeds and maintainable code has been achived.

## Further Improvements
Although data stored as hashedmaps in Go is already very efficient, the routing mechanism for the cache manager is through type switches and assertion.

Based on the benchmark carried out in this post this is an overhead that is still better than reading directly from the hard disk. https://stackoverflow.com/questions/28024884/does-a-type-assertion-type-switch-have-bad-performance-is-slow-in-go/38245293

This may possibly be rectified if the interfaces are better defined to route the data more effectively OR perhaps use a trie.
