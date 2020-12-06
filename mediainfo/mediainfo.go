package mediainfo

/*
 #cgo linux LDFLAGS: -ldl
 #cgo darwin LDFLAGS: -framework CoreFoundation
 #include <stdlib.h>
 #include "c/mediainfo_wrapper.c"
 #include "c/buffer_wrapper.c"
*/
import "C"

import (
	"errors"
	"io"
	"strings"
	"unsafe"
)

const (
	General = C.MediaInfo_Stream_General
	Video   = C.MediaInfo_Stream_Video
	Audio   = C.MediaInfo_Stream_Audio
	Image   = C.MediaInfo_Stream_Image
)

/* Dont expose this ugliness. */
type MediaInfo struct {
	ptr unsafe.Pointer
}

/* Loas the shared library. */
func Init() {
	C.mediainfo_c_init()
}

/*
 * Opens and parses the file.
 *
 * Takes a full path or reltaive path as an argument,
 * and returns a MediaInfo handler.
 */
func Open(file string) (MediaInfo, error) {
	var ret MediaInfo

	cfile := C.CString(file)
	defer C.free(unsafe.Pointer(cfile))

	cptr := C.mediainfo_c_open(cfile)
	ret.ptr = cptr
	if cptr == nil {
		return ret, errors.New("cannot open file")
	}

	return ret, nil
}

/*
 * Initializes MediaInfo structure for buffer read.
 */
func InitBuffer() (MediaInfo, error) {
	var ret MediaInfo
	cptr := C.buffer_c_init()

	ret.ptr = cptr
	if cptr == nil {
		return ret, errors.New("cannot init MediaInfo")
	}

	return ret, nil
}

/*
 * Assigns opened buffer to the MediaInfo structure and reads its info.
 * Returns true if everything is OK.
 */
func (handle MediaInfo) ReadBuffer(file io.ReadSeeker) bool {
	mi := unsafe.Pointer(handle.ptr)
	sz, _ := file.Seek(0, io.SeekEnd)
	C.buffer_c_open(mi, C.ulong(sz), C.ulong(0))

	_, _ = file.Seek(0, io.SeekStart)
	for {
		bt := make([]byte, 1024)
		rd, err := file.Read(bt)
		if err != nil {
			if err != io.EOF {
				return false
			}

			break
		}

		if rd == 0 {
			break
		}

		if C.buffer_c_continue(mi, (*C.uchar)(&bt[0]), C.ulong(rd)) != 0 {
			break
		}
	}

	return true
}

/*
 * Get audio or video info for a key.
 *
 * Matches up with the list available via:
 *     mediainfo --Info-Parameters
 *
 * Only handles one video and audio stream currently.
 *
 * Takes a key, a stream number, and a stream type as
 * arguments.
 */
func (handle MediaInfo) Get(key string, stream int, typ uint32) (string, error) {
	ckey := C.CString(key)
	cptr := unsafe.Pointer(handle.ptr)
	defer C.free(unsafe.Pointer(ckey))

	cret := C.mediainfo_c_get(cptr, ckey, C.size_t(stream), typ)
	ret := C.GoString(cret)
	if len(ret) == 0 {
		return "", errors.New("cannot get value for key")
	}

	return ret, nil
}

/*
 * Set Option
 *
 * Takes key and value strings
 */
func (handle MediaInfo) Option(key string, value string) string {
	ckey := C.CString(key)
	defer C.free(unsafe.Pointer(ckey))

	cvalue := C.CString(value)
	defer C.free(unsafe.Pointer(cvalue))

	cptr := unsafe.Pointer(handle.ptr)

	return C.GoString(C.mediainfo_c_option(cptr, ckey, cvalue))
}

/*
 * Get complete information for the file
 *
 * Takes a stream number
 */
func (handle MediaInfo) Inform(stream int) (string, error) {
	cptr := unsafe.Pointer(handle.ptr)
	cret := C.mediainfo_c_inform(cptr, C.size_t(stream))
	ret := C.GoString(cret)
	if len(ret) == 0 {
		return "", errors.New("cannot get information")
	}

	return ret, nil
}

type Info map[string]map[string]string

/*
 * Get Parsed version of file info
 *
 * Takes a stream number
 */
func (handle MediaInfo) Info(stream int) (Info, error) {
	info := make(Info)

	handle.Option("Complete", "1")
	handle.Option("Output", "CSV")
	val, err := handle.Inform(stream)
	if err != nil {
		return info, err
	}

	// log.Printf("CSV %+v\n", val)

	var section string
	lines := strings.Split(val, "\n")
	for i := range lines {
		line := lines[i]
		lineSplit := strings.SplitN(line, ",", 2)

		splitLength := len(lineSplit)

		if splitLength == 1 {
			section = lineSplit[0]
			info[section] = make(map[string]string)
		} else if splitLength == 2 {
			subsectionNoSlashes := strings.Replace(lineSplit[0], "/", " ", -1)
			subsectionTitle := strings.Title(subsectionNoSlashes)
			subsectionNoSpaces := strings.Replace(subsectionTitle, " ", "_", -1)

			if _, ok := info[section][subsectionNoSpaces]; !ok {
				if strings.Contains(subsectionNoSpaces, "Extensions") {
					info[section][subsectionNoSpaces] = strings.Split(lineSplit[1], " ")[0]
				} else {
					info[section][subsectionNoSpaces] = lineSplit[1]
				}
			}
		}
	}

	// log.Printf("INFO %+v\n", info)

	return info, nil
}

/* Close a handle. */
func (handle MediaInfo) Close() {
	cptr := unsafe.Pointer(handle.ptr)

	C.mediainfo_c_close(cptr)
}
