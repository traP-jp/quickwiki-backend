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
	"strconv"
	"strings"
	"time"
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

	tagsStr := C.GoString(C.extract(CText, C.int(num_keyword)))

	log.Println(tagsStr)

	tagsData := strings.Split(tagsStr, ",")

	var tags []Tag

	for _, tagData := range tagsData {
		tagDataSplit := strings.Split(tagData, ":")
		if len(tagDataSplit) != 2 {
			continue
		}

		tagName := tagDataSplit[0]
		tagScoreStr := tagDataSplit[1]
		tagScore, err := strconv.ParseFloat(tagScoreStr, 64)
		if err != nil {
			log.Printf("failed to convert tagScore to float: %v", err)
			continue
		}

		tag := Tag{
			WikiID:  wikiID,
			TagName: tagName,
			Score:   tagScore,
		}

		log.Printf("tag: %v\n", tag)

		tags = append(tags, tag)
	}

	return tags
}

func KeywordExtractorMulti(data []KeywordExtractorData) [][]Tag {
	var res [][]Tag
	C.initialize_python()

	for _, d := range data {
		log.Printf("d: %v\n", d)
		keywords := KeywordExtractor(d.Text, d.NumKeyword, d.WikiID)
		res = append(res, keywords)
		time.Sleep(500 * time.Millisecond)
	}

	//C.finalize_python()
	return res
}
