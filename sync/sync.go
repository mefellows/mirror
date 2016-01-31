package sync

import (
	"errors"
	"fmt"
	"github.com/mefellows/mirror/filesystem"
	utils "github.com/mefellows/mirror/filesystem/utils"
	"gopkg.in/fsnotify.v1"
	"log"
	"sync"
)

func Sync(srcRaw string, destRaw string) error {

	// Remove from src/dest strings
	src := utils.ExtractURL(srcRaw).Path
	dest := utils.ExtractURL(destRaw).Path

	// Get a handle on the to/from Filesystems
	fromFile, fromFs, err := utils.MakeFile(srcRaw)

	if err != nil {
		return err
	}

	toFile, toFs, err := utils.MakeFile(destRaw)

	if fromFile.IsDir() {

		// Build file trees on both sides of the operation, and compare
		var leftMap filesystem.FileMap
		var rightMap filesystem.FileMap
		var done sync.WaitGroup
		done.Add(2)
		go func() {
			leftMap = fromFs.FileMap(fromFile)
			done.Done()
		}()
		go func() {
			rightMap = toFs.FileMap(toFile)
			done.Done()
		}()
		done.Wait()
		diff, err := filesystem.FileMapDiff(
			leftMap, rightMap, filesystem.ModifiedComparator)

		if err == nil {
			for _, file := range diff {
				toFile = utils.MkToFile(src, dest, file)
				CopySingle(fromFs, src, toFs, toFile.Path())
			}
		} else {
			log.Printf("Error: %v\n", err)
		}
	} else {
		CopySingle(fromFs, src, toFs, dest)
	}
	return err
}

func DeleteSingle(destFs filesystem.FileSystem, destRaw string) error {
	return destFs.Delete(destRaw)
}

// Copy a single file/dir (mkdir).
func CopySingle(srcFs filesystem.FileSystem, srcRaw string, destFs filesystem.FileSystem, destRaw string) error {
	fromFile, _, err := utils.MakeFile(srcRaw)
	toFile := utils.MkToFile(srcRaw, destRaw, fromFile)

	if err != nil {
		log.Printf("Error opening dest file: %v", err)
		return errors.New(fmt.Sprintf("Error opening dest file: %v", err))
	}

	if fromFile.IsDir() {
		log.Printf("Mkdir %s -> %s\n", fromFile.Path(), toFile.Path())
		destFs.MkDir(toFile)
	} else {
		log.Printf("Copying file: %s -> %s\n", fromFile.Path(), toFile.Path())
		bytes, err := srcFs.Read(fromFile)
		err = destFs.Write(toFile, bytes, fromFile.Mode())
		if err != nil {
			log.Printf("Error copying file %s: %v", fromFile.Path(), err)
		}
	}

	return nil
}

func Watch(srcRaw string, destRaw string) error {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Fatal(err)
	}
	defer watcher.Close()

	fromFs, err := utils.GetFileSystemFromFile(srcRaw)
	toFs, err := utils.GetFileSystemFromFile(destRaw)
	src := utils.ExtractURL(srcRaw).Path
	dest := utils.ExtractURL(destRaw).Path

	done := make(chan bool)
	go func() {
		for {
			select {
			case event := <-watcher.Events:
				if event.Op&fsnotify.Remove == fsnotify.Remove {
					path := utils.RelativeFilePath(src, dest, event.Name)
					err := DeleteSingle(toFs, path)
					if err != nil {
						log.Printf("Error: %v", err)
					}
				} else {
					path := utils.RelativeFilePath(src, dest, event.Name)
					CopySingle(fromFs, event.Name, toFs, path)
				}
			case err := <-watcher.Errors:
				log.Println("Error:", err)
			}
		}
	}()

	err = watcher.Add(srcRaw)
	if err != nil {
		log.Fatal(err)
	}
	<-done
	return nil
}
