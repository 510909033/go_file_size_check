package main

import (
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

//copyDir("E:\\STUDY", "E:\\abc")
func copyDir(src string, dest string) {
	srcOriginal := src
	err := filepath.Walk(src, func(src string, f os.FileInfo, err error) error {
		if f == nil {
			return err
		}
		if f.IsDir() {
			copyDir(f.Name(), dest+"/"+f.Name())
		} else {
			destNew := strings.Replace(src, srcOriginal, dest, -1)
			_, err := CopyFile(src, destNew)
			if err != nil {
				log.Println(err)
			}
		}
		return nil
	})
	if err != nil {
		log.Printf("filepath.Walk() returned %v\n", err)
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
	dstSlices := strings.Split(dst, "\\")
	dstSlicesLen := len(dstSlices)
	destDir := ""
	for i := 0; i < dstSlicesLen-1; i++ {
		destDir = destDir + dstSlices[i] + "\\"
	}
	//dest_dir := getParentDirectory(dst)
	b, err := PathExists(destDir)
	if b == false {
		err := os.Mkdir(destDir, os.ModePerm) //在当前目录下生成md目录
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
