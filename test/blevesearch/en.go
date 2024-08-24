package main

import (
	"fmt"
	"github.com/blevesearch/bleve/v2"
	"github.com/blevesearch/bleve/v2/analysis/analyzer/custom"
	"github.com/blevesearch/bleve/v2/analysis/analyzer/keyword"
	"github.com/blevesearch/bleve/v2/analysis/token/lowercase"
	"github.com/ikawaha/bleveplugin/analysis/lang/ja"
	"log"
)

type IndexDataA struct {
	Type    string
	Title   string
	Author  string
	Content string
}

func main() {
	typeMapping := bleve.NewTextFieldMapping()
	typeMapping.Analyzer = keyword.Name
	japaneseTextFieldMapping := bleve.NewTextFieldMapping()
	japaneseTextFieldMapping.Analyzer = "ja_analyzer"
	documentMapping := bleve.NewDocumentMapping()
	documentMapping.AddFieldMappingsAt("Type", typeMapping)
	documentMapping.AddFieldMappingsAt("Title", japaneseTextFieldMapping)
	documentMapping.AddFieldMappingsAt("Author", typeMapping)
	documentMapping.AddFieldMappingsAt("Content", japaneseTextFieldMapping)

	indexMapping := bleve.NewIndexMapping()
	indexMapping.TypeField = "Type"
	indexMapping.AddDocumentMapping("sodan", documentMapping)
	err := indexMapping.AddCustomTokenizer("ja_tokenizer", map[string]interface{}{
		"type":      ja.Name,
		"dict":      ja.DictIPA,
		"base_form": true,
		"stop_tags": true,
	})
	if err != nil {
		log.Printf("failed to add custom tokenizer: %v\n", err)
		panic(err)
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
		log.Printf("failed to add custom analyzer: %v\n", err)
		panic(err)
	}

	// open a new index
	mapping := bleve.NewIndexMapping()
	index, err := bleve.New("example.bleve", mapping)
	// index, err := bleve.Open("example.bleve")
	if err != nil {
		fmt.Println(err)
		return
	}

	data := IndexDataA{
		Type:    "sodan",
		Title:   "sample",
		Author:  "ha",
		Content: "魚を食べます",
	}

	data2 := IndexDataA{
		Type:    "sodan",
		Title:   "sample2",
		Author:  "hassss",
		Content: "魚は食べられます",
	}

	// index some data
	err = index.Index("id", data)
	if err != nil {
		fmt.Println(err)
	}
	err = index.Index("id2", data2)
	if err != nil {
		fmt.Println(err)
	}

	// search for some text
	query := bleve.NewMatchQuery("食べる")
	search := bleve.NewSearchRequest(query)
	searchResults, err := index.Search(search)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(searchResults)
}
