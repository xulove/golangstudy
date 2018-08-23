package main

import (
	"encoding/csv"
	"fmt"
	"os"
)
const fileName = "2.csv"
func main() {
	data := make(map[string]string)
	data["name"] = "xuxiaofeng"
	data["age"]= "20"
	err := Write(data)
	fmt.Println(err)


}
func Write(data map[string]string)error{
	//f,err := os.OpenFile(fileName,os.O_APPEND,0666)
	f,err := os.Create(fileName)
	if err != nil {
		return err
	}
	w := csv.NewWriter(f)

	for k,v := range data{
		fmt.Println(k)
		fmt.Println(v)
		w.Write([]string{k,v})
	}
	w.Flush()
	return nil
}
