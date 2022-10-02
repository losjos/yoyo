package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"regexp"
	"strings"

	"github.com/canhlinh/hlsdl"
)

var c = &http.Client{}
var token string
var policy string
var signature string
var keypair string
var playback string
var m3u8 string
var source string

func auth() {
	LOGIN_ENDPOINT := "https://beta-api.crunchyroll.com/auth/v1/token"

	u := reader("username_ ")
	p := reader("password_ ")
	payload := strings.NewReader("username=" + strings.TrimSpace(u) + "&password=" + strings.TrimSpace(p) + "&grant_type=password&scope=offline_access")

	r, _ := http.NewRequest("POST", LOGIN_ENDPOINT, payload)
	r.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	r.Header.Add("Authorization", "Basic aHJobzlxM2F3dnNrMjJ1LXRzNWE6cHROOURteXRBU2Z6QjZvbXVsSzh6cUxzYTczVE1TY1k=")
	res, _ := c.Do(r)
	defer res.Body.Close()

	body, _ := ioutil.ReadAll(res.Body)
	re := regexp.MustCompile(`access_token".*?"(.*?)"`)
	match := re.FindStringSubmatch(string(body))
	token = match[1]
}

func params() {
	INDEX_ENDPOINT := "https://beta.crunchyroll.com/index/v2"

	r, _ := http.NewRequest("GET", INDEX_ENDPOINT, nil)
	r.Header.Add("Authorization", "Bearer "+token)
	res, _ := c.Do(r)
	defer res.Body.Close()

	body, _ := ioutil.ReadAll(res.Body)
	re := regexp.MustCompile(`cms_beta.*?"policy":"(.*?)","signature":"(.*?)","key_pair_id":"(.*?)"`)
	match := re.FindStringSubmatch(string(body))
	policy = match[1]
	signature = match[2]
	keypair = match[3]
}

func episode() {
	e := reader("episode_ ")
	EPISODE_ENDPOINT := "https://beta.crunchyroll.com/cms/v2/DE/M3/crunchyroll/objects/"
	u := strings.TrimSpace(EPISODE_ENDPOINT) + strings.TrimSpace(e) + "?Signature=" + strings.TrimSpace(signature) + "&Policy=" + strings.TrimSpace(policy) + "&Key-Pair-Id=" + strings.TrimSpace(keypair)

	r, _ := http.NewRequest("GET", u, nil)
	res, _ := c.Do(r)
	defer res.Body.Close()

	body, _ := ioutil.ReadAll(res.Body)
	re := regexp.MustCompile(`playback":"(.*?)"`)
	match := re.FindStringSubmatch(string(body))
	t := strings.Replace(match[1], "*", "", -1)
	playback = strings.Replace(t, `\u0026`, "&", -1)
}

func stream() {
	r, _ := http.NewRequest("GET", playback, nil)
	res, _ := c.Do(r)
	defer res.Body.Close()

	body, _ := ioutil.ReadAll(res.Body)
	re := regexp.MustCompile(`adaptive_hls.*?hardsub_locale":"en-US","url":"(.*?)"`)
	match := re.FindStringSubmatch(string(body))
	m3u8 = match[1]
}

func quality() {
	r, _ := http.NewRequest("GET", m3u8, nil)
	res, _ := c.Do(r)
	defer res.Body.Close()

	body, _ := ioutil.ReadAll(res.Body)
	re := regexp.MustCompile(`RESOLUTION=(.*?),`)
	match := re.FindAllStringSubmatch(string(body), -1)
	for i, _ := range match {
		fmt.Println(i, match[i][1])
	}
	fmt.Print("resolution_ ")
	var i int
	_, _ = fmt.Scanf("%d", &i)
	q := match[i][1]
	extract(q)
}

func extract(q string) {
	r, _ := http.NewRequest("GET", m3u8, nil)
	res, _ := c.Do(r)
	defer res.Body.Close()

	body, _ := ioutil.ReadAll(res.Body)
	re := regexp.MustCompile(q + `.*?".*?"((?s).*?)#`)
	match := re.FindStringSubmatch(string(body))
	source = match[1]
}

func download() {
	hlsDL := hlsdl.New(strings.TrimSpace(source), nil, "download", 64, true)
	filepath, _ := hlsDL.Download()
	fmt.Println(filepath)
}
