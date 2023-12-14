package main

import (
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/gocolly/colly"
)

type Comment struct {
	UserID     string
	Content    string
	IPDatetime string
}

func main() {

	maxComments := flag.Int("max", 10, "Max number of comments to show")
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage of %s:\n", os.Args[0])
		flag.PrintDefaults()
	}
	flag.Parse()

	url := "https://www.ptt.cc/bbs/joke/M.1481217639.A.4DF.html"

	c := colly.NewCollector()

	var comments []Comment

	c.OnHTML(".bbs-screen .article-metaline, .bbs-screen .article-content, .bbs-screen .push", func(e *colly.HTMLElement) {

		if strings.Contains(e.Attr("class"), "article-metaline") {
			return
		}

		if strings.Contains(e.Attr("class"), "push") {

			userID := e.ChildText(".push-userid")
			content := e.ChildText(".push-content")
			ipdatetime := e.ChildText(".push-ipdatetime")
			comments = append(comments, Comment{UserID: userID, Content: content, IPDatetime: ipdatetime})
		}
	})
	err := c.Visit(url)
	if err != nil {
		fmt.Println("Error:", err)
	}

	if len(comments) > *maxComments {
		comments = comments[:*maxComments]
	}

	for i, comment := range comments {
		fmt.Printf("%d. 名字：%s，留言%s，時間： %s\n", i+1, comment.UserID, comment.Content, comment.IPDatetime)
	}
}
