package search

import (
	"github.com/blevesearch/bleve/v2"
	"github.com/shogo82148/go-mecab"
	"log"
	"os"
	"strconv"
)

type IndexData struct {
	ID             int
	Type           string
	Title          string
	OwnerTraqID    string
	MessageContent string
}

func Indexing(data []IndexData) {
	tagger, err := mecab.New(map[string]string{"output-format-type": "wakati"})
	if err != nil {
		log.Println("[Error from search engine] failed to create mecab tagger: %v\n", err)
		return
	}
	defer tagger.Destroy()

	for i, d := range data {
		result, err := tagger.Parse(d.MessageContent)
		if err != nil {
			log.Println("[Error from search engine] failed to parse: %v\n", err)
			return
		}
		//log.Println(result)
		data[i].MessageContent = result
	}

	var index bleve.Index
	if _, err := os.Stat("index.bleve"); err != nil {
		// new
		indexMapping := bleve.NewIndexMapping()
		index, err = bleve.New("index.bleve", indexMapping)
		if err != nil {
			log.Printf("[Error from search engine] failed to create index: %v\n", err)
			return
		}
	} else {
		// already exists
		index, err = bleve.Open("index.bleve")
		if err != nil {
			log.Printf("[Error from search engine] failed to open index: %v\n", err)
			return
		}
	}
	for _, d := range data {
		err := index.Index(strconv.Itoa(d.ID), d)
		if err != nil {
			log.Printf("[Error from search engine] failed to index: %v\n", err)
			return
		}
	}

	docCount, err := index.DocCount()
	if err != nil {
		log.Printf("[Error from search engine] failed to get doc count: %v\n", err)
		return
	}
	log.Printf("[From search engine] Finish index successfully. doc count: %d\n\n", docCount)

	err = index.Close()
	if err != nil {
		log.Printf("[Error from search engine] failed to close index: %v\n", err)
		return
	}

	res := Search("windows", 20, 0)
	log.Println(res)
}
