package storage

import (
	"os"
	"encoding/csv"
	"xu.com/mykv/router"
)

const filename  = "test.csv"
//csv文件特点：

//每次启动程序时重csv加载map：
func Read(cache *router.Cache)error {
	var f *os.File
	var err error
	_,err = os.Stat(filename)
	if err == nil {
		f,err = os.OpenFile(filename,os.O_RDONLY,0666)
		if err != nil {
			return err
		}
		reader := csv.NewReader(f)
		records,err := reader.ReadAll()
		if err == nil{
			for _,slice := range records{
				if len(slice) ==2 {
					cache.Data[slice[0]] = slice[1]
				}
			}
		}

	}else {
		f,err = os.Create(filename)
		if err != nil {
			return err
		}
	}
	defer f.Close()
	return nil
}

func Write(data map[string]string)error{
	f,err := os.OpenFile(filename,os.O_APPEND,0666)
	if err != nil {
		return err
	}
	w := csv.NewWriter(f)
	for k,v := range data{
		w.Write([]string{k,v})
	}
	w.Flush()
	return nil
}

func Update(data []string)error{
	f,err := os.OpenFile(filename,os.O_RDONLY,0666)
	if err != nil {
		return err
	}
	r := csv.NewReader(f)
	records,_ := r.ReadAll()
	for k,v := range records{
		if v[0] == data[0]{
			records[k][1] = data[1]
		}
	}
	f.Close()

	f,err = os.Create(filename)
	if err != nil {
		return err
	}
	w := csv.NewWriter(f)
	w.WriteAll(records)
	f.Close()
	return nil
}

func Delete(data string)error{
	f,err := os.OpenFile(filename,os.O_RDONLY,0666)
	if err != nil {
		return err
	}
	r := csv.NewReader(f)
	records,_ := r.ReadAll()
	f.Close()
	f,err = os.Create(filename)
	if err != nil {
		return err
	}
	w := csv.NewWriter(f)
	for _,v := range records {
		if v[0] == data{
			continue
		}
		w.Write(v)
	}
	return nil
}