package main

import (
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
	"google.golang.org/api/googleapi/transport"
	"google.golang.org/api/youtube/v3"
)

type VideoInfo struct {
	Title        string
	Description  string
	ChannelTitle string
	ViewCount    int64
	PublishedAt  string
}

func getVideoInfo(videoID string) (*VideoInfo, error) {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	apiKey := os.Getenv("YOUTUBE_API_KEY")

	client := &http.Client{
		Transport: &transport.APIKey{Key: apiKey},
	}

	service, err := youtube.New(client)
	if err != nil {
		log.Fatalf("Error creating YouTube client: %v", err)
		return nil, err
	}

	call := service.Videos.List("snippet,statistics").Id(videoID)
	response, err := call.Do()
	if err != nil {
		log.Fatalf("Error making API call: %v", err)
		return nil, err
	}

	if len(response.Items) == 0 {
		log.Fatalf("Video not found")
		return nil, err
	}

	video := response.Items[0]
	info := &VideoInfo{
		Title:        video.Snippet.Title,
		Description:  video.Snippet.Description,
		ChannelTitle: video.Snippet.ChannelTitle,
		ViewCount:    video.Statistics.ViewCount,
		PublishedAt:  video.Snippet.PublishedAt,
	}

	return info, nil
}

func YouTubePage(w http.ResponseWriter, r *http.Request) {
	videoID := r.URL.Query().Get("v")
	if videoID == "" {
		http.Error(w, "Missing video ID", http.StatusBadRequest)
		return
	}

	info, err := getVideoInfo(videoID)
	if err != nil {
		http.Error(w, "Error retrieving video information", http.StatusInternalServerError)
		return
	}

	http.ServeFile(w, r, "index.html")
}

func main() {
	http.HandleFunc("/", YouTubePage)
	log.Fatal(http.ListenAndServe(":8085", nil))
}
