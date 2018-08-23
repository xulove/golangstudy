package router
import "net/http"
import "encoding/json"
import (
	"io/ioutil"
	"strings"
	"sync"
	"fmt"
)

type Cache struct{
	Data map[string]string
	Mutex *sync.Mutex
	
}


func (cache Cache) Post(r *http.Request,postChan chan <- map[string]string) string{
	var response map[string]string
	body,err := ioutil.ReadAll(r.Body)
	fmt.Println(body)
	if err != nil {
		panic(err)
	}
	err = json.Unmarshal(body,&response)

	if err != nil {
		panic (err)
	}
	if len(response) == 0 {
		return "none json in body"
	}

	cache.Mutex.Lock()
	defer cache.Mutex.Unlock()
	for k,v := range response {
		if _,ok := cache.Data[k];ok {
			return k + "%s already exists"
		}
		cache.Data[k] = v
	}
	fmt.Println("here1")
	postChan <- response
	fmt.Println("here2")
	return "success"
	
}

func (cache Cache) Get(r *http.Request) string{
	uri := r.RequestURI
	fmt.Println(uri)
	pathInfo := strings.Split(uri,"/")
	fmt.Println(pathInfo)
	if len(pathInfo) == 0 {
		return " no path info"
	}
	key := pathInfo[1]
	fmt.Println(key)
	if _,ok := cache.Data[key];ok {
		return cache.Data[key]
	}else {
		return key + " is not exist"
	}
}

func (cache Cache) Put(r *http.Request,putChan chan <-[]string) string{
	var response map[string]string
	body,err := ioutil.ReadAll(r.Body)
	if err != nil{
		panic(err)
	}
	if err := json.Unmarshal(body,&response);err != nil {
		panic(err)
	}
	if len(response) == 0 {
		return "no json in body"
	}
	cache.Mutex.Lock()
	defer cache.Mutex.Unlock()

	for k,v := range response{
		if _,ok := cache.Data[k] ;ok {
			cache.Data[k] = v
		}
	}
	return "success"
}

func (cache Cache)Delete(r *http.Request,delChan chan <- string) string{
	uri := r.RequestURI
	pathInfo := strings.Split(uri,"/")
	if len(pathInfo) == 0 {
		return "no path info"
	}
	key := pathInfo[1]
	cache.Mutex.Lock()
	defer cache.Mutex.Unlock()

	fmt.Println(cache.Data[key])
	if _,ok := cache.Data[key];ok {
		delete(cache.Data,key)
		delChan <- key
	}else {
		return key + "is not exist"
	}

	return "success"
}






















































