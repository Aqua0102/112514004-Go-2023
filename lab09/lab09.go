package main

import (
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/gocolly/colly"
)

// Comment 表示一條推文的結構
type Comment struct {
	UserID     string
	Content    string
	IPDatetime string
}

func main() {
	// 定義命令行參數
	maxComments := flag.Int("max", 10, "Max number of comments to show")
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage of %s:\n", os.Args[0])
		flag.PrintDefaults()
	}
	flag.Parse()

	// PTT 網頁的 URL
	url := "https://www.ptt.cc/bbs/joke/M.1481217639.A.4DF.html"

	// 創建 colly 收集器
	c := colly.NewCollector()

	// 創建變數來保存推文資訊
	var comments []Comment

	// 設置 OnHTML 事件處理程序
	c.OnHTML(".bbs-screen .article-metaline, .bbs-screen .article-content, .bbs-screen .push", func(e *colly.HTMLElement) {
		// 忽略 metaline 的內容
		if strings.Contains(e.Attr("class"), "article-metaline") {
			return
		}

		// 如果是推文
		if strings.Contains(e.Attr("class"), "push") {
			// 推文處理邏輯...
			userID := e.ChildText(".push-userid")         // 使用者ID
			content := e.ChildText(".push-content")       // 推文內容
			ipdatetime := e.ChildText(".push-ipdatetime") // 推文的時間戳

			// 將推文資訊添加到切片中
			comments = append(comments, Comment{UserID: userID, Content: content, IPDatetime: ipdatetime})
		}
	})

	// 開始爬取
	err := c.Visit(url)
	if err != nil {
		fmt.Println("Error:", err)
	}

	// 如果留言數量超過設定的上限，只顯示前面幾條
	if len(comments) > *maxComments {
		comments = comments[:*maxComments]
	}

	// 在這裡進行你想要的操作，例如將 comments 存入檔案、發送到其他系統，等等。
	// 這裡只是示例，你可以根據實際需求進行調整。
	for i, comment := range comments {
		fmt.Printf("%d. 名字：%s，留言%s，時間： %s\n", i+1, comment.UserID, comment.Content, comment.IPDatetime)
	}
}
