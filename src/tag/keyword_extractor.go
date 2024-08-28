package tag

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

type Tag struct {
	WikiID  int
	TagName string
	Score   float64
}

type KeywordExtractorData struct {
	WikiID int
	Text   string
}

func KeywordExtractorMulti(data []KeywordExtractorData, numKeyword int) [][]Tag {
	var res [][]Tag

	if _, err := os.Stat("/src/tag/tmp.txt"); err == nil {
		if err := os.Remove("/src/tag/tmp.txt"); err != nil {
			log.Printf("failed to remove file: %v", err)
			return res
		}
	}
	f, err := os.Create("/src/tag/tmp.txt")
	if err != nil {
		log.Printf("failed to create file: %v", err)
		return res
	}

	f.WriteString(fmt.Sprintf("%d\n", numKeyword))
	for _, d := range data {
		f.WriteString(fmt.Sprintf("%s\n", d.Text))
	}

	f.Close()

	out, err := exec.Command("python3", "/src/tag/keyword_extractor.py").Output()
	if err != nil {
		log.Printf("failed to run python script: %v", err)
		return res
	}

	fmt.Println(string(out))

	if string(out) != "" {
		f, err := os.Open("/src/tag/tmp.txt")
		if err != nil {
			log.Printf("failed to open file: %v", err)
			return res
		}
		var pyData []string
		scanner := bufio.NewScanner(f)
		for scanner.Scan() {
			pyData = append(pyData, scanner.Text())
		}
		f.Close()

		for i, tagsStr := range pyData {
			tagsData := strings.Split(tagsStr, ",")

			var tags []Tag

			for _, tagData := range tagsData {
				tagDataSplit := strings.Split(tagData, ":")
				if len(tagDataSplit) != 2 {
					log.Printf("tagDataSplit length is not 2: %v", tagDataSplit)
					continue
				}

				tagName := tagDataSplit[0]
				tagScoreStr := tagDataSplit[1]
				tagScore, err := strconv.ParseFloat(tagScoreStr, 64)
				if err != nil {
					log.Printf("failed to convert tagScore to float: %v", err)
					continue
				}

				tag := Tag{
					WikiID:  data[i].WikiID,
					TagName: tagName,
					Score:   tagScore,
				}

				log.Printf("tag: %v\n", tag)

				tags = append(tags, tag)
			}

			res = append(res, tags)
		}
	}

	return res
}
