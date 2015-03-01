package main

import (
	"fmt"
	//"html/template"
	"github.com/Unknwon/goconfig"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
)

var buf []byte

var SAVE_PATH = "./images/"

func isDirExits(path string) bool {
	fmt.Println(path)
	fi, err := os.Stat(path)
	checkErr(err)
	if err != nil {
		return os.IsExist(err)
	} else {
		return fi.IsDir()
	}
}

var DO_MAIN, LISTEN_PORT string
var EXPAND []string

func upload(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	if r.Method == "GET" {
		io.WriteString(w, "Only Support Post Method!")
	} else {
		file, handle, err := r.FormFile("file")
		checkErr(err)
		year := strconv.Itoa(time.Now().Year())
		month := strconv.Itoa(int(time.Now().Month()))
		day := strconv.Itoa(time.Now().Day())
		bl, ext := expand(handle.Filename)
		FINAL_PATH := SAVE_PATH + year + "/" + month + "/" + day + "/"
		fmt.Println(FINAL_PATH)
		if !isDirExits(FINAL_PATH) {
			os.MkdirAll(FINAL_PATH, 0700)
		}
		if bl {
			f, err := os.OpenFile(FINAL_PATH+fileName()+"."+ext, os.O_WRONLY|os.O_CREATE, 0666)
			io.Copy(f, file)
			checkErr(err)
			defer f.Close()
			defer file.Close()
			//io.WriteString(w, DO_MAIN+SAVE_PATH+handle.Filename)
			fmt.Println("upload success")
		} else {
			io.WriteString(w, "The FileExpand No Access!")
		}
	}
}

func expand(fileName string) (bool, string) {
	var arr []string = strings.Split(fileName, ".")
	fileExpand := arr[len(arr)-1]
	var fileType string
	for i := 0; i < len(EXPAND); i++ {
		if EXPAND[i] == fileExpand {
			fileType = fileExpand
			return true, fileType
		}
	}
	return false, fileType
}

func fileName() string {
	t := time.Now()
	year := strconv.Itoa(t.Year())
	month := strconv.Itoa(int(t.Month()))
	day := strconv.Itoa(t.Day())
	hour := strconv.Itoa(t.Hour())
	minute := strconv.Itoa(t.Minute())
	second := strconv.Itoa(t.Second())
	fmt.Println(year + month + day + hour + minute + second)
	return year + month + day + hour + minute + second
}

func initIniFile() {
	c, err := goconfig.LoadConfigFile("settings.ini")
	checkErr(err)
	domain, err := c.GetValue(goconfig.DEFAULT_SECTION, "domainName")
	checkErr(err)
	listenPort, err := c.GetValue(goconfig.DEFAULT_SECTION, "listenPort")
	checkErr(err)
	expand, err := c.GetValue(goconfig.DEFAULT_SECTION, "uploadType")
	checkErr(err)
	DO_MAIN = domain
	LISTEN_PORT = listenPort
	EXPAND = strings.Split(expand, ",")
}

func checkErr(err error) {
	if err != nil {
		err.Error()
	}
}

func main() {
	initIniFile()
	fmt.Println(fileName())
	http.HandleFunc("/upload", upload)
	err := http.ListenAndServe(":"+LISTEN_PORT, nil)
	if err != nil {
		log.Fatal("listenAndServe: ", err)
	}
}
