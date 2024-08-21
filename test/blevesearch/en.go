package main

import (
	"fmt"
	"github.com/blevesearch/bleve/v2"
)

func mai() {
	// open a new index
	// mapping := bleve.NewIndexMapping()
	// index, err := bleve.New("example.bleve", mapping)
	index, err := bleve.Open("example.bleve")
	if err != nil {
		fmt.Println(err)
		return
	}

	data := struct {
		Name string
	}{
		Name: "text",
	}

	// index some data
	index.Index("id", data)

	// search for some text
	query := bleve.NewMatchQuery("text")
	search := bleve.NewSearchRequest(query)
	searchResults, err := index.Search(search)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(searchResults)
}
