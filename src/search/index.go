package search

import (
	"github.com/blevesearch/bleve/v2"
	"github.com/blevesearch/bleve/v2/analysis/analyzer/custom"
	"github.com/blevesearch/bleve/v2/analysis/analyzer/keyword"
	"github.com/blevesearch/bleve/v2/analysis/token/lowercase"
	"github.com/blevesearch/bleve/v2/mapping"
	"github.com/ikawaha/bleveplugin/analysis/lang/ja"
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

func createIndexMapping() mapping.IndexMapping {
	typeMapping := bleve.NewTextFieldMapping()
	typeMapping.Analyzer = keyword.Name
	japaneseTextFieldMapping := bleve.NewTextFieldMapping()
	japaneseTextFieldMapping.Analyzer = "ja_analyzer"
	documentMapping := bleve.NewDocumentMapping()
	documentMapping.AddFieldMappingsAt("Type", typeMapping)
	documentMapping.AddFieldMappingsAt("Title", japaneseTextFieldMapping)
	documentMapping.AddFieldMappingsAt("OwnerTraqID", typeMapping)
	documentMapping.AddFieldMappingsAt("MessageContent", japaneseTextFieldMapping)

	indexMapping := bleve.NewIndexMapping()
	indexMapping.TypeField = "type"
	indexMapping.AddDocumentMapping("sodan", documentMapping)
	err := indexMapping.AddCustomTokenizer("ja_tokenizer", map[string]interface{}{
		"type":      ja.Name,
		"dict":      ja.DictIPA,
		"base_form": true,
		"stop_tags": true,
	})
	if err != nil {
		log.Printf("[Error from search engine] failed to add custom tokenizer: %v\n", err)
		return nil
	}
	err = indexMapping.AddCustomAnalyzer("ja_analyzer", map[string]interface{}{
		"type":      custom.Name,
		"tokenizer": "ja_tokenizer",
		"token_filters": []string{
			ja.StopWordsName,
			lowercase.Name,
		},
	})
	if err != nil {
		log.Printf("[Error from search engine] failed to add custom analyzer: %v\n", err)
		return nil
	}

	return indexMapping
}

func Indexing(data []IndexData) {
	var index bleve.Index
	if _, err := os.Stat("index.bleve"); err != nil {
		// new
		indexMapping := createIndexMapping()
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

	res := Search("windows", 10, 0)
	log.Println(res)
}
