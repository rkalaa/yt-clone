package main

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
)

type VideoMessage struct {
	Name string `json:"name"`
}

type Video struct {
	ID     string `json:"id"`
	UID    string `json:"uid"`
	Status string `json:"status"`
}

func setupDirectories() {
	// Create necessary directories
	dirs := []string{"raw_videos", "processed_videos"}
	for _, dir := range dirs {
		if _, err := os.Stat(dir); os.IsNotExist(err) {
			os.Mkdir(dir, os.ModePerm)
		}
	}
}

func isVideoNew(videoID string) bool {
	// Check if the video is new in the database
	// This is a placeholder for actual implementation
	return true
}

func setVideo(videoID string, video Video) error {
	// Store video metadata in Firestore or another DB
	// Placeholder for actual implementation
	return nil
}

func downloadRawVideo(filename string) error {
	// Placeholder for downloading video from cloud storage
	fmt.Println("Downloading raw video:", filename)
	return nil
}

func convertVideo(inputFile, outputFile string) error {
	// Placeholder for video conversion logic
	fmt.Println("Converting video:", inputFile, "to", outputFile)
	return nil
}

func uploadProcessedVideo(filename string) error {
	// Placeholder for uploading video to cloud storage
	fmt.Println("Uploading processed video:", filename)
	return nil
}

func deleteRawVideo(filename string) error {
	// Placeholder for deleting raw video
	fmt.Println("Deleting raw video:", filename)
	return nil
}

func deleteProcessedVideo(filename string) error {
	// Placeholder for deleting processed video
	fmt.Println("Deleting processed video:", filename)
	return nil
}

func handleProcessVideo(w http.ResponseWriter, r *http.Request) {
	// Parse Pub/Sub message
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Bad Request: failed to read request body", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	var reqBody map[string]interface{}
	if err := json.Unmarshal(body, &reqBody); err != nil {
		http.Error(w, "Bad Request: invalid JSON payload", http.StatusBadRequest)
		return
	}

	messageData, ok := reqBody["message"].(map[string]interface{})["data"].(string)
	if !ok {
		http.Error(w, "Bad Request: missing message data", http.StatusBadRequest)
		return
	}

	messageBytes, err := base64.StdEncoding.DecodeString(messageData)
	if err != nil {
		http.Error(w, "Bad Request: failed to decode base64 message", http.StatusBadRequest)
		return
	}

	var videoMessage VideoMessage
	if err := json.Unmarshal(messageBytes, &videoMessage); err != nil || videoMessage.Name == "" {
		http.Error(w, "Bad Request: invalid message payload", http.StatusBadRequest)
		return
	}

	inputFileName := videoMessage.Name
	outputFileName := "processed-" + inputFileName
	videoID := strings.Split(inputFileName, ".")[0]

	if !isVideoNew(videoID) {
		http.Error(w, "Bad Request: video already processing or processed", http.StatusBadRequest)
		return
	}

	if err := setVideo(videoID, Video{
		ID:     videoID,
		UID:    strings.Split(videoID, "-")[0],
		Status: "processing",
	}); err != nil {
		http.Error(w, "Internal Server Error: failed to update video metadata", http.StatusInternalServerError)
		return
	}

	// Process video
	if err := downloadRawVideo(inputFileName); err != nil {
		http.Error(w, "Internal Server Error: failed to download raw video", http.StatusInternalServerError)
		return
	}

	if err := convertVideo(inputFileName, outputFileName); err != nil {
		deleteRawVideo(inputFileName)
		deleteProcessedVideo(outputFileName)
		http.Error(w, "Internal Server Error: video processing failed", http.StatusInternalServerError)
		return
	}

	if err := uploadProcessedVideo(outputFileName); err != nil {
		http.Error(w, "Internal Server Error: failed to upload processed video", http.StatusInternalServerError)
		return
	}

	setVideo(videoID, Video{
		Status:   "processed",
		ID:       videoID,
		UID:      strings.Split(videoID, "-")[0],
		Filename: outputFileName,
	})

	deleteRawVideo(inputFileName)
	deleteProcessedVideo(outputFileName)

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Processing finished successfully"))
}

func main() {
	setupDirectories()

	http.HandleFunc("/process-video", handleProcessVideo)

	port := os.Getenv("PORT")
	if port == "" {
		port = "3000"
	}

	fmt.Printf("Video processing service listening on http://localhost:%s\n", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}