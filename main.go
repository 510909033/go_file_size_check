package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"sort"
	"strings"
	"sync/atomic"
	"time"
	"unicode/utf8"
)

var Logname = `C:\Users\Administrator\Desktop\temp.log`

type DirList []PathInfo

var FileStats = make(map[string]PathInfo)
var sepa = "\\"
var processCount int64

type PathInfo struct {
	Dir         string
	TotalSize   int64
	FolderCount int
	FileCount   int
}

func main() {
	log.SetFlags(log.Lshortfile)
	log.Println("系统：", runtime.GOOS)

	//copyDir(`E:\aaaaaa\8.63`, `E:\aaaaaa\8.63_bak`)
	//return

	var configDir = new(string)
	fmt.Fprint(os.Stdout, "输入对应的字母后回车， 如e盘输入e 然后回车\n或输入具体目录如c:/test/\n")
	fmt.Fscanln(os.Stdin, configDir)

	if *configDir == "" {
		log.Println("未配置--dir")
		flag.Usage()
		return
	}
	if runtime.GOOS != "windows" {
		Logname = "/tmp/dirchecklog.log"
	}
	Logname = Logname + time.Now().Format("_2006-01-02_15-04-05")
	if utf8.RuneCountInString(*configDir) == 1 && runtime.GOOS == "windows" {
		*configDir = (*configDir) + ":/"
	}
	log.Println("待统计目录：", *configDir)
	log.Println("结果日志存储文件：", Logname)
	WinRun(*configDir)
	//win.WinRun("e:/")
	//win.WinRun("e:/bbt/")
	//win.WinRun(`E:\bbt\test_auto_pregnancy\data`)
	//log.Printf("%+v\n",win.FileStats)
	//log.Printf("%+v\n",win.Sort())

	//bytes, _ := json.Marshal(win.FileStats)
	//log.Printf("%+v\n", string(bytes))

	fileList := Sort()
	//bytes1, _ := json.Marshal(fileList)
	//log.Printf("%+v\n", string(bytes1))

	fileList.Dump()

}

func (d DirList) Len() int {
	return len(d)
}

func (d DirList) Less(i, j int) bool {
	if d[i].TotalSize == d[j].TotalSize {
		return len(d[i].Dir) < len(d[j].Dir)
	}
	return d[i].TotalSize > d[j].TotalSize
}

func (d DirList) Swap(i, j int) {
	d[i], d[j] = d[j], d[i]
}

func Sort() DirList {
	var dirList = make(DirList, len(FileStats))
	i := 0
	for _, v := range FileStats {
		v := v
		dirList[i] = v
		i++
	}
	sort.Sort(dirList)
	return dirList
}

func (d DirList) Dump() {

	f, err := os.OpenFile(Logname, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0755)
	if err != nil {
		log.Fatal("创建文件失败：" + Logname)
	}
	defer f.Close()
	for _, v := range d {
		sizeN := v.TotalSize / 1024 / 1024
		if sizeN < 1 {
			continue
		}
		size := fmt.Sprintf("%d M", v.TotalSize/1024/1024)
		f.WriteString(size + " " + v.Dir + "\n")
	}
	f.Sync()

	s, _ := filepath.Abs(Logname)
	fmt.Println("已保存到文件：", s)
}

func fn(path string, info os.FileInfo, err error) error {
	newProcessCount := atomic.AddInt64(&processCount, 1)
	if newProcessCount%5000 == 0 {
		log.Println("process... 已统计文件数：", newProcessCount)
	}
	//log.Printf("fn path=%s, err=%+v", path,err)
	if err != nil {
		//log.Printf("d.info err=%+v", err)
		return nil
	}
	if info.IsDir() {
		return nil
	}

	var key string
	for _, v := range strings.Split(filepath.Dir(path), sepa) {
		if key == "" {
			key = v
		} else {
			key = key + sepa + v
		}
		tmpKey := strings.TrimRight(key, "\\")

		pathInfo := FileStats[tmpKey]
		pathInfo.Dir = tmpKey
		pathInfo.FileCount++
		pathInfo.TotalSize += info.Size()
		FileStats[tmpKey] = pathInfo
	}

	return nil
}

func WinRun(dir string) {

	err := filepath.Walk(dir, fn)
	if err != nil {
		log.Fatalf("filepath.WalkDir err=%+v", err)
		return
	}
}
