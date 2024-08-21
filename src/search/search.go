package search

import (
	"github.com/blevesearch/bleve/v2"
	"log"
	"os"
	"strconv"
)

// return array of ID
// return empty array if no result or error
func Search(query string, limit int, offset int) []int {
	log.Println("[Info from search] Searching by query: ", query)
	if _, err := os.Stat("index.bleve"); err != nil {
		log.Println("[Error from search] Index file does not exist. Please create index file first.")
	}

	index, err := bleve.Open("index.bleve")
	if err != nil {
		log.Printf("[Error from search] failed to open index: %v\n", err)
	}

	bleveQuery := bleve.NewMatchQuery(query)
	bleveQuery.SetField("Content")
	search := bleve.NewSearchRequest(bleveQuery)
	searchResults, err := index.Search(search)
	if err != nil {
		log.Printf("[Error from search] failed to search by query\"%s\": %v\n", query, err)
		return []int{}
	}

	err = index.Close()
	if err != nil {
		log.Printf("[Error from search] failed to close index: %v\n", err)
	}

	returnCount := min(limit, int(searchResults.Total)-offset)
	res := make([]int, returnCount)
	for i := offset; i < offset+returnCount; i++ {
		res[i-offset], err = strconv.Atoi(searchResults.Hits[i].ID)
		if err != nil {
			log.Printf("[Error from search] failed to convert ID to int: %v\n", err)
			return []int{}
		}
	}

	return res
}
