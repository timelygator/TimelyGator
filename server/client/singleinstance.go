package client

import (
	"log"
	"os"
	"path/filepath"
	"syscall"
    "fmt"

	"timelygator/server/utils"
)

// The file singleinstance.go ensures that only one instance of a client
// can run at a time by using a lock file mechanism.

type SingleInstance struct {
    lockfile string
    fd       *os.File
}


func NewSingleInstance(clientName string) (*SingleInstance, error) {
    cachedir, err := utils.GetDir("cache")
    if err != nil {
        return nil, fmt.Errorf("failed to get cache dir: %w", err)
    }
    lockfile := filepath.Join(cachedir, clientName)
    log.Printf("SingleInstance lockfile: %s", lockfile)

    instance := &SingleInstance{lockfile: lockfile}

    if os.PathSeparator == '\\' { // Windows
        if _, err := os.Stat(lockfile); err == nil {
            if err := os.Remove(lockfile); err != nil {
                return nil, fmt.Errorf("failed to remove existing lockfile: %w", err)
            }
        }
        fd, err := os.OpenFile(lockfile, os.O_CREATE|os.O_EXCL|os.O_RDWR, 0600)
        if err != nil {
            if os.IsPermission(err) {
                return nil, fmt.Errorf("another instance is already running, quitting: %w", err)
            }
            return nil, fmt.Errorf("failed to open lockfile: %w", err)
        }
        instance.fd = fd
    } else { // non-Windows
        fd, err := os.OpenFile(lockfile, os.O_CREATE|os.O_WRONLY, 0600)
        if err != nil {
            return nil, fmt.Errorf("failed to open lockfile: %w", err)
        }
        instance.fd = fd
        if err := syscall.Flock(int(fd.Fd()), syscall.LOCK_EX|syscall.LOCK_NB); err != nil {
            return nil, fmt.Errorf("another instance is already running, quitting: %w", err)
        }
    }

    return instance, nil
}

func (si *SingleInstance) Close() {
    if si.fd != nil {
        si.fd.Close()
        os.Remove(si.lockfile)
    }
}