package amara

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/url"
	"time"
)

type Video struct {
	ID                       string      `json:"id"`
	VideoType                string      `json:"video_type"`
	PrimaryAudioLanguageCode string      `json:"primary_audio_language_code"`
	OriginalLanguage         string      `json:"original_language"`
	Title                    string      `json:"title"`
	Description              string      `json:"description"`
	Duration                 int         `json:"duration"`
	Thumbnail                string      `json:"thumbnail"`
	Created                  time.Time   `json:"created"`
	Team                     interface{} `json:"team"`
	TeamType                 interface{} `json:"team_type"`
	Project                  interface{} `json:"project"`
	AllUrls                  []string    `json:"all_urls"`
	Metadata                 struct {
		SpeakerName string `json:"speaker-name"`
		Location    string `json:"location"`
	} `json:"metadata"`
	Languages []struct {
		Code         string `json:"code"`
		Name         string `json:"name"`
		Published    bool   `json:"published"`
		Dir          string `json:"dir"`
		SubtitlesURI string `json:"subtitles_uri"`
		ResourceURI  string `json:"resource_uri"`
	} `json:"languages"`
	ActivityURI          string `json:"activity_uri"`
	UrlsURI              string `json:"urls_uri"`
	SubtitleLanguagesURI string `json:"subtitle_languages_uri"`
	ResourceURI          string `json:"resource_uri"`
}

type Subtitles struct {
	VersionNumber int    `json:"version_number"`
	SubFormat     string `json:"sub_format"`
	Subtitles     string `json:"subtitles"`
	Author        struct {
		Username string `json:"username"`
		ID       string `json:"id"`
		URI      string `json:"uri"`
	} `json:"author"`
	Language struct {
		Code string `json:"code"`
		Name string `json:"name"`
		Dir  string `json:"dir"`
	} `json:"language"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Metadata    struct {
		SpeakerName string `json:"speaker-name"`
		Location    string `json:"location"`
	} `json:"metadata"`
	VideoTitle       string `json:"video_title"`
	VideoDescription string `json:"video_description"`
	ActionsURI       string `json:"actions_uri"`
	NotesURI         string `json:"notes_uri"`
	ResourceURI      string `json:"resource_uri"`
	SiteURI          string `json:"site_uri"`
	Video            string `json:"video"`
	VersionNo        int    `json:"version_no"`
}

func (c *Client) GetVideo(id string) (*Video, error) {
	data, err := c.doRequest(ReqParams{"GET", fmt.Sprintf("%s/videos/%s/", c.endpoint, id), nil})
	if err != nil {
		return nil, err
	}
	video := Video{}
	if err = json.Unmarshal(data, &video); err != nil {
		return nil, err
	}
	return &video, nil
}

func (c *Client) CreateVideo(params url.Values) (*Video, error) {
	data, err := c.doRequest(ReqParams{
		"POST",
		fmt.Sprintf("%s/", c.endpoint),
		bytes.NewBufferString(params.Encode()),
	})
	if err != nil {
		return nil, err
	}
	video := Video{}
	if err = json.Unmarshal(data, &video); err != nil {
		return nil, err
	}
	return &video, nil
}

func (c *Client) CreateSubtitles(videoID, langCode, format string, params url.Values) (*Subtitles, error) {
	if params == nil {
		return nil, errors.New("Please provide the request body parameters")
	}

	params.Set("sub_format", format)
	data, err := c.doRequest(ReqParams{
		"POST",
		fmt.Sprintf("%s/videos/%s/languages/%s/subtitles/", c.endpoint, videoID, langCode),
		bytes.NewBufferString(params.Encode()),
	})
	if err != nil {
		return nil, err
	}
	fmt.Println(string(data))
	subtitle := Subtitles{}
	if err = json.Unmarshal(data, &subtitle); err != nil {
		return nil, err
	}
	return &subtitle, nil
}

func (c *Client) GetSubtitles(videoID, langCode string) (*Subtitles, error) {
	data, err := c.doRequest(ReqParams{
		"GET",
		fmt.Sprintf("%s/videos/%s/languages/%s/subtitles/?sub_format=vtt", c.endpoint, videoID, langCode),
		nil,
	})
	if err != nil {
		return nil, err
	}
	fmt.Println(string(data))
	subtitle := Subtitles{}
	if err = json.Unmarshal(data, &subtitle); err != nil {
		return nil, err
	}
	return &subtitle, nil
}

//func (c *Client) List() ([]*Video, error) {
//	data, err := c.doRequest(ReqParams{"GET", c.endpoint})
//	if err != nil {
//		return nil, err
//	}
//
//}
