package main

import (
	"github.com/shogo82148/go-mecab"
	"log"
)

func main() {
	tagger, err := mecab.New(map[string]string{"output-format-type": "dump"})
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
