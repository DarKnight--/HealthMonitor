package target

/*
#cgo linux LDFLAGS:-L/usr/local/lib/ -lfuzzy -ldl -I/usr/local/include/
#include <stdlib.h>
#include <fuzzy.h>
*/
import "C"

import (
	"errors"
	"unsafe"
)

// HashFilename generates the fuzzy hash of the content stored in the file.
func HashFilename(filename string) (string, error) {
	outputHash := (*C.char)(C.calloc(C.FUZZY_MAX_RESULT, 1))
	defer C.free(unsafe.Pointer(outputHash))
	cfileName := C.CString(filename)
	defer C.free(unsafe.Pointer(cfileName))

	if C.fuzzy_hash_filename(cfileName, outputHash) != 0 {
		return "", errors.New("")
	}

	return C.GoString(outputHash), nil
}

// HashString generates fuzzy hash for the given byte array if no error occures
func HashString(str []byte) (string, error) {
	buf := (*C.char)(C.calloc(C.FUZZY_MAX_RESULT, 1))
	defer C.free(unsafe.Pointer(buf))

	length := C.uint32_t(len(str))
	if C.fuzzy_hash_buf((*C.uchar)(unsafe.Pointer(&str[0])), length, buf) != 0 {
		return "", errors.New("")
	}

	return C.GoString(buf), nil
}

// CompareHash compares two fuzzy hashes and return the value between 0-100.
// In case of any error -1 will be returned
func CompareHash(str1, str2 string) int {
	cstr1 := C.CString(str1)
	defer C.free(unsafe.Pointer(cstr1))

	cstr2 := C.CString(str2)
	defer C.free(unsafe.Pointer(cstr2))

	if score := C.fuzzy_compare(cstr1, cstr2); score >= 0 {
		return int(score)
	}

	return -1
}
