package scraper

import (
	"log"
	"quickwiki-backend/search"
)

func (s *Scraper) setIndexing() {
	var wikis []Wiki
	err := s.db.Select(&wikis, "SELECT * FROM wikis WHERE type = 'sodan'")
	if err != nil {
		log.Println("failed to get wikis")
		log.Println(err)
	}

	var IndexData []search.IndexData
	for _, wiki := range wikis {
		IndexData = append(IndexData, search.IndexData{
			ID:          wiki.ID,
			Type:        wiki.Type,
			Title:       wiki.Name,
			OwnerTraqID: wiki.OwnerTraqID,
			Content:     wiki.Content,
		})
	}

	search.Indexing(IndexData)
}
