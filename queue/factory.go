package queue

import (
	"errors"
	"fmt"
	"github.com/Fjolnir-Dvorak/manageAMQ/utils"
	"io/ioutil"
	"os"
	"path"
)


var (
	ErrNoFile         = errors.New("file does not exist")
	ErrNoDir          = errors.New("directory does not exist")
	ErrNoPermission   = errors.New("no permissions to access file or directory")
	ErrDunno          = errors.New("I seriously do not know what just happened")
	ErrNoDictionaries = errors.New("dictionaries are not supported")
	ErrDirIsFile      = errors.New("directory is file")
	ErrNotOpened      = errors.New("could not open file or directory")
)

func factoryError(why error, file string) (files *FileList, err error) {
	return nil, &os.PathError{Err: why, Path: file, Op: "creatingFileList"}
}

func fileListToStruct(inputFiles []string) (*FileList, error) {
	var fileList = make([]SingleFile, 0, len(inputFiles))
	for _, filename := range inputFiles {
		file, err := os.OpenFile(filename, os.O_RDONLY, os.ModePerm)
		if err != nil {
			return factoryError(ErrNotOpened, filename)
		}
		lineCount, err := utils.CountLinesFromFile(file)
		if err != nil {
			file.Close()
			return factoryError(ErrDunno, filename)
		}
		fi, err := file.Stat()
		if err != nil {
			fmt.Println(err)
			file.Close()
			return factoryError(ErrDunno, filename)
		}
		singleFile := SingleFile{
			FullPath:filename,
			ReadingPosition:0,
			ReprenstativeName:fi.Name(),
			TotalLines:lineCount,
		}
		fileList = append(fileList, singleFile)
		file.Close()
	}
	return &FileList{
		Files:       fileList,
		CurrentFile: 0,
		TotalFiles:  len(fileList),
	}, nil
}

func BuildFromFileList(inputFiles []string) (*FileList, error) {

	// At first test for permissions
	for _, filename := range inputFiles {
		fileinfo, err := os.Stat(filename)
		if err != nil {
			if os.IsNotExist(err) {
				return factoryError(ErrNoFile, filename)
			} else if os.IsPermission(err) {
				return factoryError(ErrNoPermission, filename)
			} else {
				return factoryError(ErrDunno, filename)
			}
		}
		if fileinfo.IsDir() {
			return factoryError(ErrNoDictionaries, filename)
		}
	}

	// And then do the expensive stuff
	return fileListToStruct(inputFiles)
}

func BuildFromDirectory(inputDirectory string) (*FileList, error) {
	fileinfos, err := os.Stat(inputDirectory)
	if err != nil {
		if os.IsNotExist(err) {
			return factoryError(ErrNoDir, inputDirectory)
		} else if os.IsPermission(err) {
			return factoryError(ErrNoPermission, inputDirectory)
		} else {
			return factoryError(ErrDunno, inputDirectory)
		}
	}
	if !fileinfos.IsDir() {
		return factoryError(ErrDirIsFile, inputDirectory)
	}
	fileInfos, err := ioutil.ReadDir(inputDirectory)
	if err != nil {
		return factoryError(ErrNotOpened, inputDirectory)
	}

	fileNames := make([]string, len(fileInfos), len(fileInfos))
	for index, fileInfo := range fileInfos {
		fileNames[index] = path.Join(inputDirectory, fileInfo.Name())
	}
	return BuildFromFileList(fileNames)
}
