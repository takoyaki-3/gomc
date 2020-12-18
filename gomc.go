package gomc

import (
	"os"
	"log"
	"encoding/csv"
	"io/ioutil"
	"path/filepath"
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

func WriteCSV(path string,titles map[string]int,records [][]string){
	file, err := os.Create(path)
	if err != nil {
		log.Println(err)
	}
	defer file.Close()

	writer := csv.NewWriter(file) // utf8

	header := []string{}
	for k,v:=range titles{
		for v >= len(header){
			header = append(header, "")
		}
		header[v] = k
	}

	writer.Write(header)

	for _,line:=range records{
		writer.Write(line)
	}

	writer.Flush()
}

func Dirwalk(dir string) []string {
	files, err := ioutil.ReadDir(dir)
	if err != nil {
		panic(err)
	}

	var paths []string
	for _, file := range files {
		if file.IsDir() {
			paths = append(paths, Dirwalk(filepath.Join(dir, file.Name()))...)
			continue
		}
		paths = append(paths, filepath.Join(dir, file.Name()))
	}

	return paths
}