package github

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

type Latest struct {
	ReleaseTag         string
	ReleaseSHA         string
	ReleaseDescription string
}

type latestResponse struct {
	TagName     string `json:"tag_name"`
	Description string `json:"body"`
}

type tagResponse struct {
	Object struct {
		SHA string `json:"sha"`
	} `json:"object"`
}

func get(owner, repo, path string) (*http.Response, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()

	url := fmt.Sprintf("https://api.github.com/repos/%s/%s/%s", owner, repo, path)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed create request: %w", err)
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		_ = resp.Body.Close()
		return nil, fmt.Errorf("failed request to Github API at URL='%s', %w", url, err)
	}

	if resp.StatusCode >= http.StatusMultipleChoices {
		_ = resp.Body.Close()
		return nil, fmt.Errorf("failed request to Github API at URL='%s' with StatusCode='%d'", url, resp.StatusCode)
	}

	return resp, nil
}

func getLatest(owner, repo string) (*latestResponse, error) {
	resp, err := get(owner, repo, "releases/latest")
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	result := &latestResponse{}
	err = json.NewDecoder(resp.Body).Decode(&result)
	if err != nil {
		return nil, fmt.Errorf("failed to decode latest releaes response: %w", err)
	}

	return result, nil
}

func getTag(owner, repo, tag string) (*tagResponse, error) {
	resp, err := get(owner, repo, "git/ref/tags/"+tag)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	result := &tagResponse{}
	err = json.NewDecoder(resp.Body).Decode(&result)
	if err != nil {
		return nil, fmt.Errorf("failed to decode tag response: %w", err)
	}

	return result, nil
}

func LatestVersion(owner, repo string) (*Latest, error) {
	release, err := getLatest(owner, repo)
	if err != nil {
		return nil, err
	}

	tag, err := getTag(owner, repo, release.TagName)
	if err != nil {
		return nil, err
	}

	return &Latest{
		ReleaseTag:         release.TagName,
		ReleaseSHA:         tag.Object.SHA,
		ReleaseDescription: release.Description,
	}, nil
}
