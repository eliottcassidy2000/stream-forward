package main

import (
	"encoding/json"
	"fmt"
	//"go/version"
	"io"
	"log"
	"net/http"
	"strings"
)

const (
	owner = "eliottcassidy2000"
	//repo  = "myapp"
)

type GitHubRelease struct {
	TagName string `json:"tag_name"`
}

func getLatestVersion(repo string) (string, error) {
	url := fmt.Sprintf("https://api.github.com/repos/%s/%s/releases/latest", owner, repo)
	resp, err := http.Get(url)
	fmt.Println("vvvvURL:", url)
	fmt.Println("vvvvResponse:", resp)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	var release GitHubRelease
	if err := json.NewDecoder(resp.Body).Decode(&release); err != nil {
		return "", err
	}

	return release.TagName, nil
}

func proxyRelease(w http.ResponseWriter, r *http.Request) {
	path := strings.TrimPrefix(r.URL.Path, "/")
	fmt.Println("Path:", path)
	if strings.HasPrefix(path, "latest/") {
		repo := strings.TrimPrefix(path, "latest/")
		latest, err := getLatestVersion(repo)
		if err != nil {
			http.Error(w, "Failed to fetch latest version", 500)
			log.Println("getLatestVersion error:", err)
			return
		}
		path = strings.Replace(path, "latest", latest, 1)
	}
	parts := strings.SplitN(path, "/", 2)
	repo := parts[1]
	version := parts[0]

	url := fmt.Sprintf("https://github.com/%s/%s/releases/download/%s/%s_%s_linux_amd64.tar.gz", owner, repo, version, repo, version)
	//url := fmt.Sprintf("https://github.com/%s/%s/releases/download/%s", owner, repo, version)
	fmt.Println("URL:", url)
	resp, err := http.Get(url)
	if err != nil {
		http.Error(w, "Failed to fetch file", 502)
		return
	}
	fmt.Println("Response:", resp)
	defer resp.Body.Close()

	// Copy headers
	for k, v := range resp.Header {
		w.Header()[k] = v
	}
	w.WriteHeader(resp.StatusCode)
	io.Copy(w, resp.Body)
}

func main() {
	http.HandleFunc("/", proxyRelease)
	log.Println("Listening on :8081")
	log.Fatal(http.ListenAndServe(":8081", nil))
}
