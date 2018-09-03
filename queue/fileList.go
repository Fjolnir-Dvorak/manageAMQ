package queue

import (
	"errors"
	"sync"
)

type ActiveRunner struct {
	sync.Mutex
	Cond *sync.Cond
	FileList FileList
	Running bool
	Paused bool
}

type FileList struct {
	Files []SingleFile
	CurrentFile int
	TotalFiles int
}

type SingleFile struct {
	ReprenstativeName string
	FullPath string
	ReadingPosition int
	TotalLines int
}

func NewRunner(fileList FileList) ActiveRunner {
	runner := ActiveRunner{
		Running: true,
		Paused: true,
		FileList: fileList,
	}
	runner.Cond = sync.NewCond(&runner)
	return runner
}

func (fl FileList) Append(a *FileList) {
	if a == nil {
		return
	}
	if fl.Files == nil {
		fl.Files = a.Files
		fl.CurrentFile = a.CurrentFile
		fl.TotalFiles = a.TotalFiles
	} else {
		// TODO search for duplicates
		fl.Files = append(fl.Files, a.Files...)
		fl.TotalFiles = len(fl.Files)
	}
}
func (fl FileList) GetNextFile() (*SingleFile, error) {
	if fl.CurrentFile == fl.TotalFiles {
		return nil, errors.New("EOF")
	}
	fl.CurrentFile++
	return &fl.Files[fl.CurrentFile], nil
}

func (fl FileList) HasNextFile() bool {
	return fl.CurrentFile <= fl.TotalFiles
}

func (fl FileList) Current() SingleFile {
	return fl.Files[fl.CurrentFile - 1]
}