package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"os"
	"strings"
)

func OriginalFeedBody() string {
	res, err := http.Get(os.Getenv("FEED_URL"))
	if err != nil {
		log.Fatal(err)
	}
	body, err := ioutil.ReadAll(res.Body)
	res.Body.Close()
	if err != nil {
		log.Fatal(err)
	}
	return string(body)
}

func UpdateFileUrls(body string, newHost string) string {
	return strings.Replace(body, os.Getenv("REPLACE_URL"), newHost, -1)
}

func Feed(w http.ResponseWriter, r *http.Request) {
	updatedBody := UpdateFileUrls(OriginalFeedBody(), r.Host)
	w.Header().Set("Content-Type", "application/rss+xml")
	fmt.Fprint(w, updatedBody)
}

func File(w http.ResponseWriter, r *http.Request) {
	cookieJar, err := cookiejar.New(nil)
	if err != nil {
		log.Fatal(err)
	}
	client := &http.Client{
		CheckRedirect: nil,
		Jar:           cookieJar,
	}

	fileUrl := fmt.Sprintf("%s/%s", os.Getenv("FILE_URL"), r.URL.Path[1:])
	v := url.Values{}
	v.Set("amember_login", os.Getenv("USER"))
	v.Set("amember_pass", os.Getenv("PASS"))
	client.PostForm(os.Getenv("LOGIN_URL"), v)

	fileRequest, err := http.NewRequest("GET", fileUrl, nil)
	if err != nil {
		log.Fatal(err)
	}
	audioRes, err := client.Do(fileRequest)
	if err != nil {
		log.Fatal(err)
	}
	reader := bufio.NewReader(audioRes.Body)
	w.Header().Set("Content-Length", audioRes.Header.Get("Content-Length"))
	w.Header().Set("Accept-Ranges", audioRes.Header.Get("Accept-Ranges"))
	w.Header().Set("Content-Type", audioRes.Header.Get("Content-Type"))
	for {
		line, err := reader.ReadBytes('\n')
		w.Write(line)
		if err != nil {
			break
		}
	}
	audioRes.Body.Close()
}

func Log(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("%s %s %s", r.RemoteAddr, r.Method, r.URL)
		handler.ServeHTTP(w, r)
	})
}

func main() {
	http.HandleFunc("/feed.xml", Feed)
	http.HandleFunc("/", File)
	port := fmt.Sprintf(":%s", os.Getenv("PORT"))
	log.Printf("Starting on port %s ....", port)
	http.ListenAndServe(port, Log(http.DefaultServeMux))
}
