package main

import (
	"encoding/json"
	"fmt"
	"github.com/Unknwon/goconfig"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
)

type JsonData struct {
	FileType string `json:"type"`
	FileSize string `json:"size"`
	FilePath string `json:"path"`
}

var buf []byte

func isDirExits(path string) bool {
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

type Sizer interface {
	Size() int64
}

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
		if ext == "noExt" {
			return
		}
		//
		fileSizer, _ := file.(Sizer)
		if fileSizer.Size() > customfileSize {
			io.WriteString(w, "maxFileSize "+strconv.FormatInt(customfileSize, 10)+" M")
			return
		}
		fmt.Println(fileSizer.Size())
		FINAL_PATH := "file/" + year + "/" + month + "/" + day + "/"
		if !isDirExits(FINAL_PATH) {
			os.MkdirAll(FINAL_PATH, 0700)
		}
		if bl {
			FINAL_PATH += fileName() + "." + ext
			f, err := os.OpenFile(FINAL_PATH, os.O_WRONLY|os.O_CREATE, 0666)
			io.Copy(f, file)
			checkErr(err)
			defer f.Close()
			defer file.Close()
			fileMSG := JsonData{ext, strconv.FormatInt(file.(Sizer).Size(), 10), DO_MAIN + FINAL_PATH}
			data, err := json.Marshal(fileMSG)
			io.WriteString(w, string(data))
			fmt.Println(r.Host + " SaveFile " + FINAL_PATH)
		} else {
			io.WriteString(w, "The FileExpand No Access!")
		}
	}
}

func expand(fileName string) (bool, string) {
	var arr []string = strings.Split(fileName, ".")
	fileExpand := arr[len(arr)-1]
	var fileType string
	if len(arr)-1 == 0 {
		fmt.Println("No FileExt Name")
		return false, "noExt"
	}
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
	return year + month + day + hour + minute + second
}

var customfileSize int64

func initIniFile() {
	c, err := goconfig.LoadConfigFile("settings.ini")
	checkErr(err)
	domain, err := c.GetValue(goconfig.DEFAULT_SECTION, "domainName")
	checkErr(err)
	listenPort, err := c.GetValue(goconfig.DEFAULT_SECTION, "listenPort")
	checkErr(err)
	expand, err := c.GetValue(goconfig.DEFAULT_SECTION, "uploadType")
	checkErr(err)
	maxfileSize, err := c.Int64(goconfig.DEFAULT_SECTION, "uploadSize")
	customfileSize = maxfileSize
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
	http.HandleFunc("/upload", upload)
	http.Handle("/", http.FileServer(http.Dir("./")))
	err := http.ListenAndServe(":"+LISTEN_PORT, nil)
	if err != nil {
		log.Fatal("listenAndServe: ", err)
	}
}
