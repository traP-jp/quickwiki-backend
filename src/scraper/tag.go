package scraper

import (
	"fmt"
	"log"
	"quickwiki-backend/tag"
	"regexp"
	"strings"
)

func (s *Scraper) setSodanTags() {
	var wikis []Wiki_fromDB
	err := s.db.Select(&wikis, "SELECT * FROM wikis WHERE type = 'sodan'")
	if err != nil {
		log.Println("failed to get wikis")
		log.Println(err)
	}

	s.setTag(wikis)
}

func (s *Scraper) setTag(wikis []Wiki_fromDB) {
	var input []tag.KeywordExtractorData
	for _, wiki := range wikis {
		text := processMention(processLink(wiki.Content))
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

func processLink(content string) string {
	re := regexp.MustCompile("(http|https)://[^ ]*")
	return re.ReplaceAllString(content, "")
}

func processMention(content string) string {
	re := regexp.MustCompile(`!{"type":([^!]*)}`)
	mentions := re.FindAllString(content, -1)
	res := content
	for _, mention := range mentions {
		fmt.Println(mention)
		re = regexp.MustCompile(`"raw":"(.*)",( *)"id"`)
		mentionRaw := re.FindString(mention)
		mentionRaw = mentionRaw[8 : len(mentionRaw)-1]
		quoteIndex := strings.Index(mentionRaw, "\"")
		mentionRaw = mentionRaw[:quoteIndex]
		res = strings.Replace(res, mention, mentionRaw, 1)
	}
	return res
}
