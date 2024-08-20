package tag

/*
#cgo pkg-config: python3
#cgo LDFLAGS: -L. -lkeyword_extractor -lpython3.11
#include <stdlib.h>
#include "keyword_extractor.h"
*/
import "C"
import (
	"unsafe"
)

type Tag struct {
	TagName string
	Score   float64
}

func KeywordExtractor(text string, num_keyword int) []Tag {
	CText := C.CString(text)
	defer C.free(unsafe.Pointer(CText))

	cppData := C.extract(CText, C.int(num_keyword))
	defer C.free_data_array(cppData)

	tagNames := make([]string, cppData.size)
	scores := make([]float64, cppData.size)

	for i := 0; i < int(cppData.size); i++ {
		tagNamesPtr := (**C.char)(unsafe.Pointer(uintptr(unsafe.Pointer(cppData.tag_names)) + uintptr(i)*unsafe.Sizeof(cppData.tag_names)))
		scoresPtr := (*C.double)(unsafe.Pointer(uintptr(unsafe.Pointer(cppData.scores)) + uintptr(i)*unsafe.Sizeof(cppData.scores)))

		tagNames[i] = C.GoString(*tagNamesPtr)
		scores[i] = float64(*scoresPtr)
	}

	Tags := make([]Tag, cppData.size)
	for i := 0; i < int(cppData.size); i++ {
		Tags[i] = Tag{TagName: tagNames[i], Score: scores[i]}
	}

	return Tags
}
