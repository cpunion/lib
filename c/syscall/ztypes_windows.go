//go:build windows

package syscall

const (
	ERROR_FILE_NOT_FOUND Errno = 2
	ERROR_PATH_NOT_FOUND Errno = 3
)

// Timespec is kept for source compatibility with packages that share Unix and
// Windows declarations. Windows APIs use FILETIME or the UCRT timespec layout
// at their ABI boundaries instead.
type Timespec struct {
	Sec  int64
	Nsec int64
}

// These Unix-only resource structures have no Windows ABI. Empty definitions
// allow shared declaration packages to expose unsupported calls without
// inventing a binary layout for them.
type Rusage struct{}
type Rlimit struct{}

// Stat_t matches the Universal CRT's 64-bit _stat64 layout. The layout is the
// same for Windows AMD64 and ARM64.
type Stat_t struct {
	Dev   uint32
	Ino   uint16
	Mode  uint16
	Nlink int16
	Uid   int16
	Gid   int16
	Rdev  uint32
	Size  int64
	Atime int64
	Mtime int64
	Ctime int64
}
