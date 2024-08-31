package scraper

import (
	"log"
	"quickwiki-backend/model"
	"quickwiki-backend/search"
)

func (s *Scraper) SetIndexing() {
	var wikis []model.WikiContent_fromDB
	err := s.db.Select(&wikis, "SELECT * FROM wikis WHERE type = 'sodan'")
	if err != nil {
		log.Println("failed to get wikis")
		log.Println(err)
	}

	var IndexData []search.IndexData
	for _, wiki := range wikis {
		IndexData = append(IndexData, search.IndexData{
			ID:             wiki.ID,
			Type:           wiki.Type,
			Title:          wiki.Name,
			OwnerTraqID:    wiki.OwnerTraqID,
			MessageContent: ProcessLink(removeNewLine(removeCodeBlock(removeTeX(wiki.Content)))),
		})
	}

	search.Indexing(IndexData)
}
