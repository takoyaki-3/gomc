package gomc

import (
	"os"
	"log"
	"encoding/csv"
	"io"
	"bufio"
	"io/ioutil"
	"path/filepath"
	"net/http"
	"archive/zip"
	"fmt"
)

func newCsvReader(r io.Reader) *csv.Reader {
	br := bufio.NewReader(r)
	bs, err := br.Peek(3)
	if err != nil {
		return csv.NewReader(br)
	}
	if bs[0] == 0xEF && bs[1] == 0xBB && bs[2] == 0xBF {
		br.Discard(3)
	}
	return csv.NewReader(br)
}

func ReadCSV(path string)(map[string]int,[][]string){
	file, err := os.Open(path)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	reader := newCsvReader(file)
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

func DownloadURL(url string,path string)error{
	client := &http.Client{}
	req, err := http.NewRequest("GET", url, nil)
	resp, err := client.Do(req)
	defer resp.Body.Close()
	if err != nil {
		return err
	}

	out, err := os.Create(path)
	defer out.Close()
	
	_,err = io.Copy(out, resp.Body)

	return err
}

func Unzip(src, dest string) error {
	r, err := zip.OpenReader(src)
	if err != nil {
		return err
	}
	defer r.Close()

	for _, f := range r.File {
		rc, err := f.Open()
		if err != nil {
			return err
		}
		defer rc.Close()

		path := filepath.Join(dest, f.Name)
		if f.FileInfo().IsDir() {
			os.MkdirAll(path, f.Mode())
		} else {
			f, err := os.OpenFile(
					path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, f.Mode())
			if err != nil {
					return err
			}
			defer f.Close()

			_, err = io.Copy(f, rc)
			if err != nil {
					return err
			}
		}
	}
	return nil
}

func Ziparchive(output string, paths []string) error {
	var compressedFile *os.File
	var err error

	//ZIPファイル作成
	if compressedFile, err = os.Create(output); err != nil {
			return err
	}
	defer compressedFile.Close()

	if err := compress(compressedFile, ".", paths); err != nil {
			return err
	}

	return nil
}

func compress(compressedFile io.Writer, targetDir string, files []string) error {
	w := zip.NewWriter(compressedFile)

	for _, filename := range files {
			filepath := fmt.Sprintf("%s/%s", targetDir, "./onegtfs/" + filename)
			info, err := os.Stat(filepath)
			if err != nil {
					return err
			}

			if info.IsDir() {
					continue
			}

			file, err := os.Open(filepath)
			if err != nil {
					return err
			}
			defer file.Close()

			hdr, err := zip.FileInfoHeader(info)
			if err != nil {
					return err
			}

			hdr.Name = filename

			f, err := w.CreateHeader(hdr)
			if err != nil {
					return err
			}

			contents, _ := ioutil.ReadFile(filepath)
			_, err = f.Write(contents)
			if err != nil {
					return err
			}
	}

	if err := w.Close(); err != nil {
			return err
	}
	return nil
}

func DirwalkPF(dir string) ([]string,[]string) {
	files, err := ioutil.ReadDir(dir)
	if err != nil {
		panic(err)
	}

	var paths,file_names []string
	for _, file := range files {
		paths = append(paths, filepath.Join(dir, file.Name()))
		file_names = append(file_names,file.Name())
	}
	return paths,file_names
}