package main

import "net/http"
import (
	"xu.com/mykv/router"
	"sync"
	"xu.com/mykv/storage"
	"fmt"
)

func main() {
	postChan := make(chan map[string]string)
	putChan := make(chan []string)
	delChan := make(chan string)
	cache := router.Cache{Data: make(map[string]string), Mutex: &sync.Mutex{}}

	storage.Read(&cache)

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case "POST":
			fmt.Println("post")
			result := cache.Post(r, postChan)
			fmt.Println(result)
			fmt.Fprint(w, result)
			break
		case "GET":
			fmt.Println("get")
			result := cache.Get(r)
			fmt.Fprint(w, result)
			break
		case "PUT":
			fmt.Println("put")
			result := cache.Put(r,putChan)
			fmt.Fprint(w, result)
			break
		case "DELETE":
			fmt.Println("delete")
			result := cache.Delete(r,delChan)
			fmt.Fprint(w, result)
			break
		default:
			fmt.Println("/")
			fmt.Fprint(w, "does not support this method")
		}
	})
	go func() {
		for {
			fmt.Println("又一轮通道接受")
			select {
			case addData :=  <-postChan:
				fmt.Println("postChan中接受到了数据")
				storage.Write(addData)
			case updateData := <-putChan:
				storage.Update(updateData)
			case key := <-delChan:
				storage.Delete(key)
			}
		}
	}()
	http.ListenAndServe(":8070", nil)

}

/*
method: get url: http://127.0.0.1:8070/user 获取user的值
method: post url: http://127.0.0.1:8070/ {"user": "1"} 新增user的值
method: put url: http://127.0.0.1:8070/ {"user": "2"} 修改user的值
method: delete url: http://127.0.0.1:8070/user 删除user的值
*/
//func Router(w http.ResponseWriter, r *http.Request, cache router.Cache) {
//	switch r.Method {
//	case "POST":
//		result := cache.Post(r, postChan)
//		fmt.Sprint(w, result)
//		break
//	case "GET":
//		result := cache.Get(r)
//		fmt.Sprint(w, result)
//		break
//	case "PUT":
//		result := cache.Put(r)
//		fmt.Sprint(w, result)
//		break
//	case "DELETE":
//		result := cache.Delete(r)
//		fmt.Sprint(w, result)
//		break
//	default:
//		fmt.Sprint(w, "does not support this method")
//	}
//}
