package file

import (
	"context"
	"fmt"
	"io"
	"os"

	log "github.com/sirupsen/logrus"
	"golang.org/x/sync/semaphore"
)

var Sem *semaphore.Weighted

func init() {
	Sem = semaphore.NewWeighted(1000)
}

func PathExists(filename string) bool {
	_, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false
	}
	return true
}

func DirIsEmpty(name string) (bool, error) {
	f, err := os.Open(name)
	if err != nil {
		return false, err
	}
	defer f.Close()

	_, err = f.Readdirnames(1) // Or f.Readdir(1)
	if err == io.EOF {
		return true, nil
	}
	return false, err // Either not empty or error, suits both cases
}

func DirExists(filename string) bool {
	info, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false
	}
	return info.IsDir()
}

func FileExists(filename string) bool {
	info, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
}

func ReadBytes(filename string) ([]byte, error) {
	Sem.Acquire(context.Background(), 1)
	defer Sem.Release(1)
	if handle, err := os.Open(filename); err == nil {
		defer handle.Close()
		stat, err := handle.Stat()
		if err != nil {
			log.Debugf("Could not stat() %v", filename)
			return []byte{}, err
		}
		result := make([]byte, stat.Size())
		bytes_read, err := handle.Read(result)
		if err != nil {
			log.Debugf("Could not Read() %v", filename)
			return []byte{}, err
		}
		if int64(bytes_read) != stat.Size() {
			return []byte{}, fmt.Errorf("Did not read expected number of bytes: %v", bytes_read)
		}
		return result, nil
	} else {
		return []byte{}, fmt.Errorf("File %v could not be opened; maybe as it doesn't exist: %v", filename, err)
	}
}
