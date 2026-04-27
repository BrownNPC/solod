package os

import (
	"unsafe"

	"solod.dev/so/c"
	"solod.dev/so/mem"
	"solod.dev/so/slices"
	"solod.dev/so/strings"
)

// ReadDir reads the named directory, returning
// all its directory entries sorted by filename.
// If an error occurs reading the directory, returns the entries it was
// able to read before the error, along with the error.
//
// If the allocator is nil, uses the system allocator.
// The returned slice and entry names are allocated; the caller owns them.
// Use [FreeDirEntry] to free the result.
func ReadDir(a mem.Allocator, name string) ([]DirEntry, error) {
	entries, err := readDir(a, name)
	slices.SortFunc(entries, compareEntry)
	return entries, err
}

// FreeDirEntry frees a slice of DirEntry previously returned by [ReadDir].
// It frees each entry's Name string and the slice itself.
//
// If the allocator is nil, uses the system allocator.
func FreeDirEntry(a mem.Allocator, entries []DirEntry) {
	for i := range entries {
		mem.FreeString(a, entries[i].Name)
	}
	slices.Free(a, entries)
}

// readDir reads the named directory and returns
// all its directory entries without sorting.
func readDir(a mem.Allocator, name string) ([]DirEntry, error) {
	dir := opendir(name)
	if dir == nil {
		return []DirEntry{}, mapError()
	}

	entries := slices.MakeCap[DirEntry](a, 0, 16)
	var nameBuf [MaxNameLen]byte
	cname := (*c.Char)(&nameBuf[0])

	for {
		r := os_readdir_next(dir, cname, MaxNameLen)
		if !r.ok {
			break
		}

		entryName := unsafe.String((*byte)(cname), int(r.nameLen))

		// Skip "." and "..".
		if entryName == "." || entryName == ".." {
			continue
		}

		dm := dtypeToMode(r.dtype)
		isDir := dm.isDir
		mode := dm.mode

		// DT_UNKNOWN: fall back to Lstat.
		if r.dtype == dtUnknown {
			fi, err := Lstat(name + "/" + entryName)
			if err == nil {
				isDir = fi.IsDir()
				mode = fi.Mode() & (ModeDir | ModeSymlink | ModeNamedPipe | ModeSocket | ModeDevice | ModeCharDevice | ModeIrregular)
			}
		}

		clonedName := strings.Clone(a, entryName)
		entries = slices.Append(a, entries, DirEntry{
			Name:  clonedName,
			IsDir: isDir,
			Type:  mode,
		})
	}

	// Check for read error (errno set by readdir).
	readErr := errno
	closedir(dir)

	if readErr != 0 {
		errno = readErr // restore errno from readdir
		return entries, mapError()
	}
	return entries, nil
}

// compareEntry compares two DirEntry values by their name for sorting.
func compareEntry(a, b any) int {
	e1 := c.PtrAs[DirEntry](a)
	e2 := c.PtrAs[DirEntry](b)
	return strings.Compare(e1.Name, e2.Name)
}

// dtypeModeResult holds the result of dtypeToMode.
type dtypeModeResult struct {
	isDir bool
	mode  FileMode
}

// dtypeToMode converts a dirent d_type value to isDir and FileMode type bits.
func dtypeToMode(dtype uint8) dtypeModeResult {
	switch dtype {
	case dtDIR:
		return dtypeModeResult{isDir: true, mode: ModeDir}
	case dtLNK:
		return dtypeModeResult{mode: ModeSymlink}
	case dtFIFO:
		return dtypeModeResult{mode: ModeNamedPipe}
	case dtSOCK:
		return dtypeModeResult{mode: ModeSocket}
	case dtBLK:
		return dtypeModeResult{mode: ModeDevice}
	case dtCHR:
		return dtypeModeResult{mode: ModeCharDevice}
	default:
		return dtypeModeResult{}
	}
}
