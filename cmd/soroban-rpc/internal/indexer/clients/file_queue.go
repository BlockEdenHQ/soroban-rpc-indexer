package clients

import (
	supportlog "github.com/stellar/go/support/log"
	"os"
	"sync"
)

type FileQueue struct {
	Filename     string
	PositionFile string
	file         *os.File
	lock         sync.Mutex
	logger       *supportlog.Entry
}

func NewFileQueue(filename, positionFile string, logger *supportlog.Entry) *FileQueue {
	fq := &FileQueue{Filename: filename, PositionFile: positionFile, logger: logger}
	fq.initialize()
	return fq
}

func (fq *FileQueue) initialize() {
	fq.lock.Lock()
	defer fq.lock.Unlock()

	var err error
	fq.file, err = os.OpenFile(fq.Filename, os.O_APPEND|os.O_CREATE|os.O_RDWR, 0644)
	if err != nil {
		fq.logger.Errorf("Error: Failed to open queue file: %v", err)
	}

}

func (fq *FileQueue) Enqueue(data string) {
	fq.lock.Lock()
	defer fq.lock.Unlock()

	if _, err := fq.file.WriteString(data + "\n"); err != nil {
		fq.logger.Errorf("Error: Failed to enqueue data: %v", err)
	}
}

func (fq *FileQueue) Close() {
	if err := fq.file.Close(); err != nil {
		fq.logger.WithError(err).Error("Error: Failed to close queue file")
	}
}
