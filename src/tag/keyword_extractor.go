package tag

/*
#cgo pkg-config: python3
#cgo LDFLAGS: -L. -lkeyword_extractor -lpython3.11
#include <stdlib.h>
#include "keyword_extractor.h"
*/
import "C"
import (
	"log"
	"unsafe"
)

type Tag struct {
	WikiID  int
	TagName string
	Score   float64
}

type KeywordExtractorData struct {
	WikiID     int
	Text       string
	NumKeyword int
}

func KeywordExtractor(text string, num_keyword int, wikiID int) []Tag {
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

	log.Println("-----------------")
	log.Println("num_keyword: ", num_keyword)
	Tags := make([]Tag, cppData.size)
	for i := 0; i < int(cppData.size); i++ {
		Tags[i] = Tag{WikiID: wikiID, TagName: tagNames[i], Score: scores[i]}
		log.Printf("tag: %s, score: %f", Tags[i].TagName, Tags[i].Score)
	}

	return Tags
}

func KeywordExtractorMulti(data []KeywordExtractorData) [][]Tag {
	var res [][]Tag
	C.initialize_python()

	for _, d := range data {
		res = append(res, KeywordExtractor(d.Text, d.NumKeyword, d.WikiID))
	}

	C.finalize_python()
	return res
}
