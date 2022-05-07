package main

import (
	"bufio"
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"gopkg.babytree-inc.com/bgf/awesomeProject1/cody_dir"
	"io"
	"log"
	"os"
	"path"
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
	var configSwitch = new(string)

	fmt.Fprint(os.Stdout, "请选择功能，输入对应数字\n")
	fmt.Fprint(os.Stdout, "1：目录文件大小检查\n")
	fmt.Fprint(os.Stdout, "2：目录文件复制\n")

	fmt.Fscanln(os.Stdin, configSwitch)

	switch *configSwitch {
	case "1":
		case1()
	case "2":
		case2()
	default:
		fmt.Fprint(os.Stdout, "错误的选项\n")
		return
	}

}

func case1() {
	var configDir = new(string)
	var maxLevel = new(int)
	fmt.Fprint(os.Stdout, "输入对应的字母后回车， 如e盘输入e 然后回车\n或输入具体目录如c:/test/\n")
	fmt.Fscanln(os.Stdin, configDir)

	fmt.Fprint(os.Stdout, "输入最大展示目录层级，最小为1\n")
	fmt.Fscanln(os.Stdin, maxLevel)

	if *configDir == "" {
		log.Println("未配置--dir")
		return
	}

	if *maxLevel < 1 {
		log.Println("最大展示目录层级不能小于1")
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

	WinRun(*configDir, fnRepeatFile)

	fileList := Sort()

	fileList.Dump(*maxLevel)
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

func (d DirList) Dump(maxLevel int) {

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
		if strings.Count(v.Dir, `\`) > maxLevel {
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

func WinRun(dir string, fn filepath.WalkFunc) {
	err := filepath.Walk(dir, fn)
	if err != nil {
		log.Fatalf("filepath.WalkDir err=%+v", err)
		return
	}
}

//查找相同重复的图片、文件等， 并且列出位置
func fnRepeatFile(path string, info os.FileInfo, err error) error {
	newProcessCount := atomic.AddInt64(&processCount, 1)
	if newProcessCount%5000 == 0 {
		log.Println("process... 已统计文件数：", newProcessCount)
	}
	if err != nil {
		log.Printf("d.info err=%+v", err)
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

type Case2 struct {
	Cnt     int
	Files   []string
	Md5File string
}

func case2() {
	var baseDir string
	//baseDir := "e:/test111/test222"
	//baseDir := `e:\test114`
	dirList := make([]string, 0)
	//dirList = []string{`e:\pprof`}
	//dirList = []string{`E:\1-aa`}
	//dirList = []string{`C:\Users\Administrator\Documents\WeChat Files\baotian0506\FileStorage\Image`}
	//dirList = []string{`C:\Users\Administrator\Documents\WeChat Files\baotian0506\FileStorage\Video`}
	//dirList = []string{`e:\chls`}

	var ptr = new(string)

	rd := bufio.NewReader(os.Stdin)

	for {
		fmt.Fprint(os.Stdout, "输入待复制目录，例如c:/test/\n")
		str, err := rd.ReadString('\n')
		str = strings.TrimSpace(str)
		cody_dir.Debug(func() {
			log.Println(str, err)
		})
		if str == "" {
			fmt.Fprint(os.Stdout, "请输入回车开始执行\n")
			str, _ = rd.ReadString('\n')
			str = strings.TrimSpace(str)
			if str == "" {
				break
			}
			continue
		}
		dirList = append(dirList, str)
	}

	fmt.Fprint(os.Stdout, "输入要将文件复制到哪个目录，例如c:/test/\n")
	str, _ := rd.ReadString('\n')
	str = strings.TrimSpace(str)
	baseDir = str

	fmt.Fprint(os.Stdout, "请确认信息\n")
	fmt.Fprint(os.Stdout, "待复制目录列表\n")
	index := 0
	for _, v := range dirList {
		index++
		fmt.Fprintf(os.Stdout, "\t%d： %s\n", index, v)
	}

	fmt.Fprint(os.Stdout, "将文件复制到的目录\n")
	fmt.Fprintf(os.Stdout, " -- \t %s\n", baseDir)

	fmt.Fprint(os.Stdout, "输入回车确认，其他字符+回车取消\n")
	fmt.Fscanln(os.Stdin, ptr)
	if *ptr != "" {
		fmt.Fprint(os.Stdout, "已取消\n")
		os.Exit(1)
	}

	var allFiles = make(map[string]*Case2)

	callback := func(src, desc string) bool {
		cody_dir.Debug(func() {
			log.Printf("将要复制一个文件, src=%s, dest=%s\n", src, desc)
		})

		//stat, err := os.Stat(src)
		//if err!=nil{
		//	log.Println(cody_dir.ErrMsg(fmt.Sprintf("os.Stat报错，src=%s, err=%+v", src, err)))
		//	return false
		//}

		fileMd5, err := HashFileMd5(src)
		if err != nil {
			msg := fmt.Sprintf("HashFileMd5报错, err=%+v, src=%s", err, src)
			log.Println(cody_dir.ErrMsg(msg))
			return false
		}
		if allFiles[fileMd5] == nil {
			allFiles[fileMd5] = &Case2{
				Files: make([]string, 0),
			}
		}
		allFiles[fileMd5].Cnt++
		allFiles[fileMd5].Md5File = fileMd5
		allFiles[fileMd5].Files = append(allFiles[fileMd5].Files, src)

		if allFiles[fileMd5].Cnt > 1 {
			existFile(src, allFiles[fileMd5])
		}

		return true
	}

	for _, v := range dirList {
		log.Printf("%s的path.Base=%s, filepath.Base=%s\n", v, path.Base(v), filepath.Base(v))
		if exists, err := cody_dir.PathExists(v); err != nil {
			log.Println(cody_dir.ErrMsg(fmt.Sprintf("代码报错，path=%s, err=%+v\n", v, err)))
			return
		} else if !exists {
			log.Println(cody_dir.ErrMsg("目录不存在，" + v + "\n"))
			return
		}

		cody_dir.CopyDir(v, baseDir+`\`+filepath.Base(v), callback)
	}

}

func existFile(src string, case2 *Case2) {
	//todo
	return
	log.Printf("文件重复%d次, src=%s\n", case2.Cnt, src)
	log.Print("重复文件：\n")
	for _, v := range case2.Files {
		log.Printf("\t\t%s\n", v)
	}
	log.Print("\n")
}

func HashFileMd5(filePath string) (string, error) {
	var returnMD5String string
	file, err := os.Open(filePath)
	if err != nil {
		return returnMD5String, err
	}
	defer file.Close()
	hash := md5.New()
	if _, err := io.Copy(hash, file); err != nil {
		return returnMD5String, err
	}
	hashInBytes := hash.Sum(nil)[:16]
	returnMD5String = hex.EncodeToString(hashInBytes)
	return returnMD5String, nil
}
