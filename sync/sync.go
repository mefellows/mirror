package sync

import (
	"errors"
	"fmt"
	"github.com/mefellows/mirror/filesystem"
	utils "github.com/mefellows/mirror/filesystem/utils"
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
	if fromFile.IsDir() {
		toFile, toFs, err := utils.MakeFile(destRaw)

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

				if err == nil {
					if file.IsDir() {
						log.Printf("Mkdir: %s -> %s\n", file.Path(), toFile.Path())
						toFs.MkDir(toFile)
					} else {
						log.Printf("Copying file: %s -> %s\n", file.Path(), toFile.Path())
						bytes, err := fromFs.Read(file)
						err = toFs.Write(toFile, bytes, file.Mode())
						if err != nil {
							log.Printf("Error copying file %s: %v", file.Path(), err)
						}
					}
				}
			}
		} else {
			log.Printf("Error: %v\n", err)
		}
	} else {
		toFile := utils.MkToFile(src, dest, fromFile)
		toFs, err := utils.GetFileSystemFromFile(destRaw)
		if err != nil {
			log.Printf("Error opening dest file: %v", err)
			return errors.New(fmt.Sprintf("Error opening dest file: %v", err))
		}

		bytes, err := fromFs.Read(fromFile)
		if err != nil {
			log.Printf("Error reading from source file: %s", err.Error())
			return errors.New(fmt.Sprintf("Error reading source file: %v", err))
		}
		err = toFs.Write(toFile, bytes, toFile.Mode())
		if err != nil {
			log.Printf("Error writing to remote path: %s", err.Error())
			return errors.New(fmt.Sprintf("Error write to remote path: %v", err))
		}
	}
	return err
}
