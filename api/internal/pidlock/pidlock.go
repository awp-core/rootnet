// Package pidlock provides single-instance process locking via flock.
// Each service (api, indexer, keeper) acquires an exclusive lock on a
// pidfile at startup. If another instance is already running, the new
// process exits immediately with a clear error message.
package pidlock

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"syscall"
)

// Lock represents a held pidfile lock.
type Lock struct {
	file *os.File
}

// Acquire attempts to obtain an exclusive lock on the given pidfile path.
// Returns a Lock on success. If another process holds the lock, returns an error.
// The lock is automatically released when the process exits.
func Acquire(name string) (*Lock, error) {
	dir := os.Getenv("PID_DIR")
	if dir == "" {
		dir = os.TempDir()
	}
	path := filepath.Join(dir, fmt.Sprintf("awp-%s.lock", name))

	f, err := os.OpenFile(path, os.O_CREATE|os.O_RDWR, 0644)
	if err != nil {
		return nil, fmt.Errorf("pidlock: open %s: %w", path, err)
	}

	// Set close-on-exec to prevent child processes from inheriting the lock fd
	syscall.CloseOnExec(int(f.Fd()))

	// Try non-blocking exclusive lock
	err = syscall.Flock(int(f.Fd()), syscall.LOCK_EX|syscall.LOCK_NB)
	if err != nil {
		// Read existing PID for diagnostic message
		data := make([]byte, 32)
		n, _ := f.Read(data)
		_ = f.Close()
		existingPID := string(data[:n])
		return nil, fmt.Errorf("pidlock: another %s instance is already running (pid %s, lockfile %s)", name, existingPID, path)
	}

	// Write our PID
	_ = f.Truncate(0)
	_, _ = f.Seek(0, 0)
	_, _ = f.WriteString(strconv.Itoa(os.Getpid()))
	_ = f.Sync()

	return &Lock{file: f}, nil
}

// Release explicitly releases the lock. Usually not needed — the OS releases
// flock automatically when the process exits.
func (l *Lock) Release() {
	if l.file != nil {
		_ = syscall.Flock(int(l.file.Fd()), syscall.LOCK_UN)
		_ = l.file.Close()
		l.file = nil
	}
}
