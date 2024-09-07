package main

import (
	"github.com/blevesearch/bleve/v2"
	"github.com/shogo82148/go-mecab"
	"log"
	"os"
	"strconv"
	"time"
)

func main() {
	tagger, err := mecab.New(map[string]string{"output-format-type": "wakati"})
	if err != nil {
		log.Println(err)
		return
	}
	defer tagger.Destroy()

	data := []IndexDatas{
		{
			ID:             1,
			Type:           "type",
			Title:          "title1",
			OwnerTraqID:    "owner1",
			MessageContent: "GoはGoogleによって作られた，オープンソースのプログラミング言語です。",
			CreatedAt:      time.Now(),
		},
		{
			ID:             2,
			Type:           "type",
			Title:          "title2",
			OwnerTraqID:    "owner2",
			MessageContent: "swiftはGoogleによって作られた，プロプライエタリです。",
			CreatedAt:      time.Now().Add(-48 * time.Hour),
		},
		{
			ID:             3,
			Type:           "type",
			Title:          "title3",
			OwnerTraqID:    "owner3",
			MessageContent: "PythonはGoogleによって作られた，オープンソースのプログラミング言語です。",
			CreatedAt:      time.Now().Add(-24 * time.Hour),
		},
		{
			ID:             4,
			Type:           "type",
			Title:          "title4",
			OwnerTraqID:    "owner4",
			MessageContent: "RubyはYukihiro Matsumotoによって作られた，オープンソースのプログラミング言語です。",
			CreatedAt:      time.Now().Add(-72 * time.Hour),
		},
		{
			ID:             5,
			Type:           "type",
			Title:          "title5",
			OwnerTraqID:    "owner5",
			MessageContent: "JavaはSun Microsystemsによって作られた，オープンソースのプログラミング言語です。",
			CreatedAt:      time.Now().Add(-96 * time.Hour),
		},
		{
			ID:             6,
			Type:           "type",
			Title:          "title6",
			OwnerTraqID:    "owner6",
			MessageContent: "C++はBjarne Stroustrupによって作られたオープンソースのプログラミング言語です。",
			CreatedAt:      time.Now().Add(-120 * time.Hour),
		},
	}

	for i, d := range data {
		result, err := tagger.Parse(d.MessageContent)
		if err != nil {
			log.Println(err)
			return
		}
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

	bleveQuery := bleve.NewMatchQuery("プログラミング言語")
	bleveQuery.SetField("MessageContent")
	search := bleve.NewSearchRequestOptions(bleveQuery, 10, 0, false)
	search.SearchBefore = []string{"Google"}
	//search.SortBy([]string{"CreatedAt", "-_score"})
	searchResults, err := index.Search(search)
	if err != nil {
		log.Printf("[Error from search] failed to search by query\"%s\": %v\n", "プログラミング言語", err)
		return
	}

	log.Printf("[Info from search] Found %d results\n", len(searchResults.Hits))
	for _, hit := range searchResults.Hits {
		log.Printf("ID: %s, Value: %+v\n", hit.ID, hit)
	}

	err = index.Close()
	if err != nil {
		log.Printf("[Error from search engine] failed to close index: %v\n", err)
		return
	}
}

type IndexDatas struct {
	ID             int
	Type           string
	Title          string
	OwnerTraqID    string
	MessageContent string
	CreatedAt      time.Time
}
