package queue

import (
	"errors"
	"sync"
)

type ActiveRunner struct {
	sync.Mutex
	Channel chan struct{}
	Waiter sync.WaitGroup
	FileList *FileList
	Running bool
	Paused bool
}

type FileList struct {
	Files []SingleFile `json:"files"`
	CurrentFile int `json:"currentFile"`
	TotalFiles int `json:"totalFiles"`
}

type SingleFile struct {
	ReprenstativeName string `json:"representativeName"`
	FullPath string `json:"fullPath"`
	ReadingPosition int `json:"readingPosition"`
	TotalLines int `json:"totalLines"`
}

func NewRunner(fileList *FileList) *ActiveRunner {
	runner := ActiveRunner{
		Running: true,
		Paused: true,
		FileList: fileList,
	}
	runner.Channel = make(chan struct{}, 1)
	return &runner
}

func (fl *FileList) Append(a *FileList) {
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
func (fl *FileList) GetNextFile() (*SingleFile, error) {
	if fl.CurrentFile >= fl.TotalFiles {
		return nil, errors.New("EOF")
	}
	fl.CurrentFile++
	return &fl.Files[fl.CurrentFile - 1], nil
}

func (fl *FileList) HasNextFile() bool {
	return fl.CurrentFile <= fl.TotalFiles
}

func (fl *FileList) CurrentOrFirst() SingleFile {
	if fl.CurrentFile <= 0 {
		return fl.Files[fl.CurrentFile]
	}
	return fl.Files[fl.CurrentFile - 1]
}