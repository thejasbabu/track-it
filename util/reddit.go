package util

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math/rand"
	"net/http"
	"time"
)

const (
	QuoteSubredditURL = "https://www.reddit.com/r/quotes/"
)

type RedditClient struct {
	SubReddit string
}

type ChildData struct {
	SubRedditName string `json:"subreddit_name_prefixed"`
	Title         string `json:"title"`
}

type Children struct {
	Kind string    `json:"kind"`
	Data ChildData `json:"data"`
}

type RedditData struct {
	Dist     int        `json:"dist"`
	Children []Children `json:"children"`
}

type RedditTopResponse struct {
	Kind string     `json:"kind"`
	Data RedditData `json:"data"`
}

func (r RedditClient) TopTitle() (string, error) {
	req, err := http.NewRequest("GET", fmt.Sprintf("https://www.reddit.com/r/%s/top.json?count=10", r.SubReddit), nil)
	if err != nil {
		return "", err
	}
	client := http.Client{
		Timeout: time.Duration(5 * time.Second),
	}
	headers := http.Header{}
	headers.Add("user-agent", "localhost")
	req.Header = headers
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("non successful response code: %d", resp.StatusCode)
	}
	var redditResp RedditTopResponse
	bytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	err = json.Unmarshal(bytes, &redditResp)
	if err != nil {
		return "", err
	}
	rand.Seed(time.Now().UnixNano())
	maxQuotes := len(redditResp.Data.Children)
	i := randomInt(0, maxQuotes)
	return redditResp.Data.Children[i].Data.Title, nil
}

func randomInt(min, max int) int {
	return min + rand.Intn(max-min)
}
