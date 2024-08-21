package scraper

import (
	"log"
	"quickwiki-backend/tag"
)

func (s *Scraper) setSodanTags() {
	var wikis []Wiki
	err := s.db.Select(&wikis, "SELECT * FROM wikis WHERE type = 'sodan'")
	if err != nil {
		log.Println("failed to get wikis")
		log.Println(err)
	}

	s.setTag(wikis)
}

func (s *Scraper) setTag(wikis []Wiki) {
	var input []tag.KeywordExtractorData
	for _, wiki := range wikis {
		text := wiki.Content
		input = append(input, tag.KeywordExtractorData{WikiID: wiki.ID, Text: text, NumKeyword: 5})
	}

	tags := tag.KeywordExtractorMulti(input)
	for _, tg := range tags {
		for _, t := range tg {
			s.insertTag(t)
		}
	}
}

func (s *Scraper) insertTag(t tag.Tag) {
	_, err := s.db.Exec("INSERT INTO tags (wiki_id, name, tag_score) VALUES (?, ?, ?)", t.WikiID, t.TagName, t.Score)
	if err != nil {
		log.Printf("failed to insert tag: %v", err)
	}
}
