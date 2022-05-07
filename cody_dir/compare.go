package cody_dir

import (
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"
)

//
//func compare(rawSourceDir, dst string) {
//	err := filepath.Walk(rawSourceDir, compareFn)
//	if err != nil {
//		log.Fatalf("filepath.WalkDir err=%+v", err)
//		return
//	}
//
//}
//
//func compareFn(path string, info os.FileInfo, err error) error {
//
//	if err != nil {
//		log.Printf("err=%+v", err)
//		return nil
//	}
//	if info.IsDir() {
//		return nil
//	}
//
//	return nil
//}
var debug = false

func Debug(fn func()) {
	if !debug {
		return
	}
	fn()
}

func ErrMsg(msg string) string {
	return msg
}

//copyDir("E:\\STUDY", "E:\\abc")
var debugWalkAll = make(map[string]int)

func CopyDir(root string, dest string, callback func(src, dest string) bool) {
	srcOriginal := root
	err := filepath.Walk(root, func(src string, f os.FileInfo, err error) error {
		if f == nil {
			log.Println("err=", err)
			return err
		}

		Debug(func() {
			log.Printf("EVERY WALK每一个walk回调, src=%s\n", src)
		})

		debugWalkAll[src]++
		if debugWalkAll[src] > 1 {
			log.Println(ErrMsg(fmt.Sprintf("src=%+v, walk次数超过一次， cnt=%d", src, debugWalkAll[src])))
			log.Fatal("....")
		}
		if f.IsDir() {
			//if root == src {
			//	Debug(func() {
			//		log.Printf("fipath.Walk方法的root和fn的src相同， 跳过 root=%s\n", root)
			//	})
			//	return nil
			//}
			//Debug(func() {
			//	log.Printf("DIR root=%s, src=%s\n", root, src)
			//})
			//CopyDir(src, dest+`\`+f.Name(), callback)
		} else {
			Debug(func() {
				//log.Println(src, srcOriginal, dest)
			})
			destNew := strings.Replace(src, srcOriginal, dest, -1)
			if callback(src, destNew) {
				_, err := CopyFile(src, destNew)
				if err != nil {
					log.Println(err)
				}
			}
		}
		return nil
	})
	if err != nil {
		log.Printf("filepath.Walk() returned %v, root=%s, desc=%s\n", err, root, dest)
	}
}

//
////egodic directories
//func getFilelist(path string) {
//	err := filepath.Walk(path, func(path string, f os.FileInfo, err error) error {
//		if f == nil {
//			return err
//		}
//		if f.IsDir() {
//			return nil
//		}
//		println(path)
//		return nil
//	})
//	if err != nil {
//		fmt.Printf("filepath.Walk() returned %v\n", err)
//	}
//}

func PathExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}

//copy file
func CopyFile(src, dst string) (w int64, err error) {
	srcFile, err := os.Open(src)
	if err != nil {
		log.Println(err.Error())
		return
	}

	defer func() {
		err := srcFile.Close()
		if err != nil {
			log.Println(err)
		}
	}()

	srcFileInfo, err := srcFile.Stat()
	if err != nil {
		log.Println(err.Error())
		return
	}

	//
	//dstReadFile, err := os.Open(dst)
	//if err != nil {
	//	log.Println(err.Error())
	//	return
	//}
	//dstReadFile
	dstFileInfo, err := os.Stat(dst)
	if err != nil {
		if !os.IsNotExist(err) {
			log.Println(err.Error())
			return
		}
		err = nil
	}
	if dstFileInfo != nil && srcFileInfo != nil && dstFileInfo.Size() == srcFileInfo.Size() {
		Debug(func() {
			//log.Println("文件大小一致，跳过", dstFileInfo.Size(), src, dst)
		})
		return
	}

	dstSlices := strings.Split(dst, "\\")
	dstSlicesLen := len(dstSlices)
	destDir := ""
	for i := 0; i < dstSlicesLen-1; i++ {
		destDir = destDir + dstSlices[i] + "\\"
	}
	//dest_dir := getParentDirectory(dst)
	b, err := PathExists(destDir)
	if b == false {
		err := os.MkdirAll(destDir, os.ModePerm) //在当前目录下生成md目录
		if err != nil {
			log.Println(err)
		}
	}

	dstFile, err := os.Create(dst)

	if err != nil {
		log.Println(err.Error())
		return
	}

	defer func() {
		err := dstFile.Close()
		if err != nil {
			log.Println(err.Error())
			return
		}
	}()

	return io.Copy(dstFile, srcFile)
}
