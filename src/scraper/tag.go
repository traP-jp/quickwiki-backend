package scraper

import (
	"fmt"
	"log"
	"os"
	"quickwiki-backend/model"
	"quickwiki-backend/tag"
	"regexp"
	"strconv"
	"strings"
)

func (s *Scraper) setSodanTags() {
	var wikis []model.WikiContent_fromDB
	err := s.db.Select(&wikis, "SELECT * FROM wikis WHERE type = 'sodan'")
	if err != nil {
		log.Println("failed to get wikis")
		log.Println(err)
	}

	s.setTag(wikis)
}

func (s *Scraper) setTag(wikis []model.WikiContent_fromDB) {
	var input []tag.KeywordExtractorData
	for _, wiki := range wikis {
		text := ProcessMentionAll(ProcessLink(removeNewLine(wiki.Content)))
		input = append(input, tag.KeywordExtractorData{WikiID: wiki.ID, Text: text})
	}

	numKeyword, err := strconv.Atoi(os.Getenv("NUM_KEYWORD"))
	if err != nil {
		log.Println("please set NUM_KEYWORD correctly")
		numKeyword = 5
	}
	tags := tag.KeywordExtractorMulti(input, numKeyword)
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

func ProcessLink(content string) string {
	re := regexp.MustCompile("(http|https)://[^ ]*")
	return re.ReplaceAllString(content, "")
}

func ProcessMention(content string) string {
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

func ProcessMentionAll(content string) string {
	re := regexp.MustCompile(`!{"type":([^!]*)}`)
	res := re.ReplaceAllString(content, "")
	return res
}

func removeNewLine(content string) string {
	return strings.Replace(content, "\n", "", -1)
}
