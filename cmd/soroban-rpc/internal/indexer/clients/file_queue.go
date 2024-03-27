package clients

import (
	"bufio"
	supportlog "github.com/stellar/go/support/log"
	"os"
	"strconv"
	"sync"
)

type FileQueue struct {
	Filename     string
	PositionFile string
	file         *os.File
	position     int
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

	fq.position = fq.readPosition()
}

func (fq *FileQueue) Enqueue(data string) {
	fq.lock.Lock()
	defer fq.lock.Unlock()

	if _, err := fq.file.WriteString(data + "\n"); err != nil {
		fq.logger.Errorf("Error: Failed to enqueue data: %v", err)
	}
}

func (fq *FileQueue) Dequeue() string {
	fq.lock.Lock()
	defer fq.lock.Unlock()

	// Reset to start if the current position is not set
	if fq.position == 0 {
		if _, err := fq.file.Seek(0, 0); err != nil {
			fq.logger.Errorf("Error: Failed to seek to start of the file: %v", err)
		}
	}

	scanner := bufio.NewScanner(fq.file)
	for i := 0; i < fq.position; i++ {
		if !scanner.Scan() && scanner.Err() != nil {
			fq.logger.Errorf("Error during file scan: %v", scanner.Err())
		}
	}

	if scanner.Scan() {
		line := scanner.Text()
		fq.position++
		fq.savePosition(fq.position) // Update position for the next read
		return line
	}

	return "" // Reached EOF or no more lines
}

func (fq *FileQueue) Close() {
	if err := fq.file.Close(); err != nil {
		fq.logger.WithError(err).Error("Error: Failed to close queue file")
	}
}

func (fq *FileQueue) readPosition() int {
	data, err := os.ReadFile(fq.PositionFile)
	if err != nil {
		return 0 // Assume start position if file doesn't exist or can't be read
	}
	position, err := strconv.Atoi(string(data))
	if err != nil {
		return 0 // Default to start position on error
	}
	return position
}

func (fq *FileQueue) savePosition(position int) {
	err := os.WriteFile(fq.PositionFile, []byte(strconv.Itoa(position)), 0644)
	if err != nil {
		fq.logger.WithError(err).Error("Error: Failed to save read position")
	}
}
