package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/joho/godotenv"
	"google.golang.org/api/googleapi/transport"
	"google.golang.org/api/youtube/v3"
)

type VideoInfo struct {
	Title        string
	Description  string
	PublishedAt  time.Time
	ChannelTitle string
	LikeCount    int64
	ViewCount    int64
	CommentCount int64
}

type PageData struct {
	Title        string
	Id           string
	ChannelTitle string
	LikeCount    string
	ViewCount    string
	PublishedAt  string
	CommentCount string
}

func init() {
	if err := godotenv.Load("env.env"); err != nil {
		log.Fatal("Error loading .env file")
	}
}

func formatNumber(number int64) string {
	strNumber := strconv.FormatInt(number, 10)

	formattedNumber := ""
	for i := len(strNumber) - 1; i >= 0; i-- {
		formattedNumber = string(strNumber[i]) + formattedNumber
		if (len(strNumber)-i)%3 == 0 && i > 0 {
			formattedNumber = "," + formattedNumber
		}
	}

	return formattedNumber
}

func getVideoInfo(videoID string) (*VideoInfo, error) {
	developerKey := os.Getenv("YOUTUBE_API_KEY")
	client := &http.Client{
		Transport: &transport.APIKey{Key: developerKey},
	}

	service, err := youtube.New(client)
	if err != nil {
		return nil, err
	}

	call := service.Videos.List([]string{"snippet", "statistics"}).
		Id(videoID)

	response, err := call.Do()
	if err != nil {
		return nil, err
	}

	if len(response.Items) == 0 {
		return nil, fmt.Errorf("video not found")
	}

	video := response.Items[0]
	publishedAt, _ := time.Parse(time.RFC3339, video.Snippet.PublishedAt)

	info := &VideoInfo{
		Title:        video.Snippet.Title,
		Description:  video.Snippet.Description,
		PublishedAt:  publishedAt,
		ChannelTitle: video.Snippet.ChannelTitle,
		LikeCount:    int64(video.Statistics.LikeCount),
		ViewCount:    int64(video.Statistics.ViewCount),
		CommentCount: int64(video.Statistics.CommentCount),
	}

	return info, nil
}

func YouTubePage(w http.ResponseWriter, r *http.Request) {
	videoID := r.URL.Query().Get("v")
	if videoID == "" {
		RenderErrorPage(w, fmt.Errorf("missing video ID in the URL"))
		return
	}

	videoInfo, err := getVideoInfo(videoID)
	if err != nil {
		RenderErrorPage(w, err)
		return
	}

	data := PageData{
		Title:        videoInfo.Title,
		Id:           videoID,
		ChannelTitle: videoInfo.ChannelTitle,
		LikeCount:    formatNumber(videoInfo.LikeCount),
		ViewCount:    formatNumber(videoInfo.ViewCount),
		PublishedAt:  videoInfo.PublishedAt.Format("2006年01月02日"),
		CommentCount: formatNumber(videoInfo.CommentCount),
	}

	renderPage(w, data)
}

func renderPage(w http.ResponseWriter, data PageData) {
	tmpl, err := template.ParseFiles("index.html")
	if err != nil {
		RenderErrorPage(w, err)
		return
	}

	err = tmpl.Execute(w, data)
	if err != nil {
		RenderErrorPage(w, err)
	}
}

func RenderErrorPage(w http.ResponseWriter, err error) {
	tmpl, parseErr := template.ParseFiles("error.html")
	if parseErr != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusInternalServerError)

	errorData := struct {
		ErrorMessage string
	}{
		ErrorMessage: err.Error(),
	}

	execErr := tmpl.Execute(w, errorData)
	if execErr != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}
}
func main() {
	http.HandleFunc("/", YouTubePage)
	log.Fatal(http.ListenAndServe(":8085", nil))
}
