package sync

import (
	"errors"
	"fmt"
	"github.com/mefellows/mirror/filesystem"
	utils "github.com/mefellows/mirror/filesystem/utils"
	"log"
)

func Sync(srcRaw string, destRaw string) error {

	// Remove from src/dest strings
	src := utils.ExtractURL(srcRaw).Path
	dest := utils.ExtractURL(destRaw).Path
	fmt.Printf("Src: %s, Dest: %s", src, dest)

	// Get a handle on the to/from Filesystems
	fromFile, fromFs, err := utils.MakeFile(srcRaw)

	if err != nil {
		log.Printf("Error opening src file: %v", err)
		return err
	}
	if fromFile.IsDir() {
		toFile, toFs, err := utils.MakeFile(destRaw)
		diff, err := filesystem.FileTreeDiff(
			fromFs.FileTree(fromFile), toFs.FileTree(toFile), filesystem.ModifiedComparator)

		if err == nil {
			for _, file := range diff {
				toFile = utils.MkToFile(src, dest, file)

				if err == nil {
					if file.IsDir() {
						//log.Printf("Mkdir: %s -> %s\n", file.Path(), toFile.Path())
						toFs.MkDir(toFile)
					} else {
						//log.Printf("Copying file: %s -> %s\n", file.Path(), toFile.Path())
						bytes, err := fromFs.Read(file)
						//log.Printf("Read bytes: %s\n", len(bytes))
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
		log.Printf("Current mode: %s", fromFile.Mode())
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
