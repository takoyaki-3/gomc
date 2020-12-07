package gomc

import (
	"os"
	"encoding/csv"
)

func ReadCSV(path string)(map[string]int,[][]string){
	file, err := os.Open(path)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	reader := csv.NewReader(file)
	reader.FieldsPerRecord = -1
	var line []string

	counter:=-1
	titles :=map[string]int{}
	data := [][]string{}
	for {
		counter++
		line, err = reader.Read()
		if err != nil {
			break
		}
		if counter==0{
			for k,v:=range line{
				titles[v]=k
			}
			continue
		}
		data = append(data,line)
	}

	return titles,data
}