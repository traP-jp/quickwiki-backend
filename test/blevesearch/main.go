package main

import (
	"github.com/blevesearch/bleve/v2"
	"github.com/blevesearch/bleve/v2/analysis/analyzer/custom"
	"github.com/blevesearch/bleve/v2/analysis/analyzer/keyword"
	"github.com/blevesearch/bleve/v2/analysis/token/lowercase"
	"github.com/ikawaha/bleveplugin/analysis/lang/ja"
	"github.com/shogo82148/go-mecab"
	"log"
)

type IndexData struct {
	Type    string
	Title   string
	Author  string
	Content string
}

func main() {
	setting()
	index, err := bleve.Open("index.bleve")
	if err != nil {
		log.Println(err)
		return
	}

	// search for some text
	query := bleve.NewMatchQuery("作る")
	query.SetField("Content")
	search := bleve.NewSearchRequest(query)
	searchResults, err := index.Search(search)
	if err != nil {
		log.Println(err)
		return
	}
	log.Println(searchResults)
	for _, hit := range searchResults.Hits {
		log.Printf("ID: %s, Score: %f\n", hit.ID, hit.Score)
	}
	index.Close()

	settinrMecab()
}

func settinrMecab() {
	tagger, err := mecab.New(map[string]string{"output-format-type": "wakati"})
	if err != nil {
		log.Println(err)
		return
	}
	defer tagger.Destroy()

	text := "GoはGoogleによって作られたオープンソースのプログラミング言語です。"
	result, err := tagger.Parse(text)
	if err != nil {
		log.Println(err)
		return
	}
	log.Println(result)
}

func setting() {
	// open a new index
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

	//var index bleve.Index
	//if _, err := os.Stat("index.bleve"); err != nil {
	//	index, err = bleve.New("index.bleve", indexMapping)
	//} else {
	//	index, err = bleve.Open("index.bleve")
	//}
	index, err := bleve.New("index.bleve", indexMapping)
	if err != nil {
		panic(err)
	}

	data1 := IndexData{
		Type:    "sodan",
		Title:   "Go言語入門",
		Author:  "山田 太郎",
		Content: "GoはGoogleによって作られたオープンソースのプログラミング言語です。",
	}

	data2 := IndexData{
		Type:    "sodan",
		Title:   "Goの並行処理",
		Author:  "佐藤 花子",
		Content: "Goの並行処理はゴルーチンとチャネルを使用して実現されます。",
	}

	data3 := IndexData{
		Type:    "sodan",
		Title:   "Go言語の学習",
		Author:  "鈴木 一郎",
		Content: "この本はGoプログラミング言語の基本をカバーしています。",
	}

	data4 := IndexData{
		Type:    "sodan",
		Title:   "Goのマスター",
		Author:  "田中 次郎",
		Content: "この本はGoプログラミングの高度な技術に深く掘り下げています。",
	}

	data5 := IndexData{
		Type:    "sodan",
		Title:   "GoでのWeb開発",
		Author:  "高橋 三郎",
		Content: "Goは効率的なWebアプリケーションを作るために使用できます。",
	}
	// index some data
	index.Index("id1", data1)
	index.Index("id2", data2)
	index.Index("id3", data3)
	index.Index("id4", data4)
	index.Index("id5", data5)

	docCount, err := index.DocCount()
	if err != nil {
		log.Printf("failed to get doc count: %v\n", err)
		panic(err)
	}
	log.Printf("doc count: %d\n\n", docCount)
	index.Close()
}
