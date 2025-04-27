package main

import (
	"encoding/json"
	"log"
	"io"
	"net/http"
	"fmt"
	"strings"
)

const (
	owner = "eliottcassidy2000"
)

type GitHubRelease struct {
	TagName string `json:"tag_name"`
}

func getLatestVersion(repo string) (string, error) {
	url := fmt.Sprintf("https://api.github.com/repos/%s/%s/releases/latest", owner, repo)
	resp, err := http.Get(url)
	log.Println("vvvvURL:", url)
	log.Println("vvvvResponse:", resp)
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
	parts := strings.SplitN(path, "/", 4)
	log.Println("Parts:", parts)
	//https://github.com/eliottcassidy2000/stream-forward/releases/download/0.0.0/stream-forward_0.0.0_linux_${attr.cpu.arch}.tar.gz
	if len(parts) != 4 {
		http.Error(w, "Request must contain 4 parts: /version/repo/arch/os", http.StatusBadRequest)
		return
	}
	version := parts[0]
	repo := parts[1]
	arch := parts[2]
	os := parts[3]
	if version == "latest" {
		latest, err := getLatestVersion(repo)
		if err != nil {
			http.Error(w, "Failed to fetch latest version", 500)
			log.Println("getLatestVersion error:", err)
			return
		}
		version = latest
	}
	//  https://github.com/eliottcassidy2000/stream-forward/releases/download//stream-forward__${attr.os.name}_${attr.cpu.arch}.tar.gz
	//  https://github.com/eliottcassidy2000/monad-forwarder/releases/download/0.0.3/monad-forwarder_0.0.3_linux_arm64.tar.gz
	url := fmt.Sprintf("https://github.com/%s/%s/releases/download/%s/%s_%s_%s_%s.tar.gz", owner, repo, version, repo, version, os, arch)
	//url := log.Sprintf("https://github.com/%s/%s/releases/download/%s", owner, repo, version)
	log.Println("URL:", url)
	resp, err := http.Get(url)
	if err != nil {
		http.Error(w, "Failed to fetch file", 502)
		return
	}
	log.Println("Response:", resp)
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
	log.Println("Listening on :8085")
	log.Fatal(http.ListenAndServe(":8085", nil))
}
