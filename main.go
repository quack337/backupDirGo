package main

import (
	"log"
	"os"
	"slices"
	"sort"

	"github.com/quack337/goLib/fs"
	"golang.org/x/text/language"
	"golang.org/x/text/message"
)

var prt = message.NewPrinter(language.English)

func main() {
	if len(os.Args) != 3 {
		log.Fatal("Usage: backupDirGo <srcDir> <dstDir>")
	}
	var srcDirPath = os.Args[1]
	var dstDirPath = os.Args[2]

	BackupDir(srcDirPath, dstDirPath)
}

func BackupDir(srcDirPath, dstDirPath string) {
	flag, _ := fs.DirExists(srcDirPath)
	if !flag { 
		log.Fatal("srcDir not exists") 
	}
	flag, err := fs.IsHidden(srcDirPath)
	if flag || err != nil {
		return
	}
	flag, _ = fs.DirExists(dstDirPath)
	if !flag { 
		var err = os.MkdirAll(dstDirPath, os.ModePerm)
		if err != nil {
			log.Fatal(err) 
		}
	}
	srcFiles, srcDirs, err := fs.GetEntries(srcDirPath)
	if err != nil {
		log.Fatal(err)
	}
	dstFiles, dstDirs, err := fs.GetEntries(dstDirPath)
	if err != nil {
		log.Fatal(err)
	}
	backupFiles(srcDirPath, dstDirPath, srcFiles, dstFiles)

	sort.Strings(srcDirs)
	sort.Strings(dstDirs)
	for _, dstDir := range dstDirs {
		_, found := slices.BinarySearch(srcDirs, dstDir)
		if !found {
			os.RemoveAll(dstDirPath + "/" + dstDir)
		}
	}
	for _, srcDir := range srcDirs {
		srcSubDirPath := srcDirPath + "/" + srcDir
		dstSubDirPath := dstDirPath + "/" + srcDir
		BackupDir(srcSubDirPath, dstSubDirPath)
	}
}

func backupFiles(srcDir, dstDir string, srcFiles, dstFiles []fs.FileInfo) {
	fs.SortFileInfos(srcFiles)
	fs.SortFileInfos(dstFiles)

	for _, dstFile := range dstFiles {
		var index = fs.BinarySearchFileInfoByName(srcFiles, dstFile.Name)
		if index < 0 {
			var err = os.Remove(dstDir + "/" + dstFile.Name)
			if err != nil {
				log.Fatal(err)
			}
		}	
	}
	for _, srcFile := range srcFiles {
		var index = fs.BinarySearchFileInfoByName(dstFiles, srcFile.Name)
		if index >= 0 &&
		   srcFile.ModTime == dstFiles[index].ModTime &&
		   srcFile.Size == dstFiles[index].Size {
			continue
		}
		var srcFilePath = srcDir + "/" + srcFile.Name
		var dstFilePath = dstDir + "/" + srcFile.Name

		flag, err := fs.IsHidden(srcFilePath)
		if flag || err != nil {
			continue
		}
		_, err = fs.CopyFile(srcFilePath, dstFilePath)
		if err != nil {
			log.Fatal(err)
		}
		err = fs.CopyFileTime(srcDir + "/" + srcFile.Name, dstFilePath)
		if err != nil {
			log.Fatal(err)				
		}
		prt.Printf("  %s   %d bytes\n", srcDir + "/" + srcFile.Name, srcFile.Size)
	}
}