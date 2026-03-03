package main

import (
    "encoding/json"
    "fmt"
    "net/http"
    "time"
)

type Client struct {
    BaseURL  string
    Username string
    Password string
    http     *http.Client
}

func NewClient(baseURL, username, password string) *Client {
    return &Client{
        BaseURL:  baseURL,
        Username: username,
        Password: password,
        http:     &http.Client{Timeout: 15 * time.Second},
    }
}

type Category struct {
    ID   string `json:"category_id"`
    Name string `json:"category_name"`
}

type LiveStream struct {
    ID                 int    `json:"stream_id"`
    Name               string `json:"name"`
    CategoryID         string `json:"category_id"`
    ContainerExtension string `json:"container_extension"`
}

type VODStream struct {
    ID                 int        `json:"stream_id"`
    Name               string     `json:"name"`
    CategoryID         string     `json:"category_id"`
    ContainerExtension string     `json:"container_extension"`
    Rating             FlexString `json:"rating"`   // era string
    Plot               string     `json:"plot"`
}

type Series struct {
    ID         int    `json:"series_id"`
    Name       string `json:"name"`
    CategoryID string `json:"category_id"`
    Plot       string `json:"plot"`
    Rating     string `json:"rating"`
}

type SeriesInfo struct {
    Info    struct{ Name string `json:"name"` } `json:"info"`
    Seasons map[string][]Episode               `json:"episodes"`
}

type Episode struct {
    ID                 string     `json:"id"`
    EpisodeNum         FlexInt    `json:"episode_num"`  // era int
    Title              string     `json:"title"`
    ContainerExtension string     `json:"container_extension"`
    Season             FlexInt    `json:"season"`       // pode ter o mesmo problema
}

type UserInfo struct {
    Username string `json:"username"`
    Status   string `json:"status"`
    ExpDate  string `json:"exp_date"`
}

func (c *Client) apiURL(action string, extras ...string) string {
    url := fmt.Sprintf("%s/player_api.php?username=%s&password=%s&action=%s",
        c.BaseURL, c.Username, c.Password, action)
    for _, e := range extras {
        if e != "" {
            url += "&" + e
        }
    }
    return url
}

func (c *Client) get(url string, v any) error {
    resp, err := c.http.Get(url)
    if err != nil {
        return fmt.Errorf("request: %w", err)
    }
    defer resp.Body.Close()
    if resp.StatusCode != http.StatusOK {
        return fmt.Errorf("HTTP %d", resp.StatusCode)
    }
    return json.NewDecoder(resp.Body).Decode(v)
}

func (c *Client) Authenticate() (*UserInfo, error) {
    var r struct {
        UserInfo UserInfo `json:"user_info"`
    }
    if err := c.get(c.apiURL(""), &r); err != nil {
        return nil, err
    }
    if r.UserInfo.Username == "" {
        return nil, fmt.Errorf("credenciais invalidas")
    }
    return &r.UserInfo, nil
}

func (c *Client) GetLiveCategories() ([]Category, error) {
    var v []Category
    return v, c.get(c.apiURL("get_live_categories"), &v)
}

func (c *Client) GetLiveStreams(catID string) ([]LiveStream, error) {
    var v []LiveStream
    extra := ""
    if catID != "" {
        extra = "category_id=" + catID
    }
    return v, c.get(c.apiURL("get_live_streams", extra), &v)
}

func (c *Client) GetVODCategories() ([]Category, error) {
    var v []Category
    return v, c.get(c.apiURL("get_vod_categories"), &v)
}

func (c *Client) GetVODStreams(catID string) ([]VODStream, error) {
    var v []VODStream
    extra := ""
    if catID != "" {
        extra = "category_id=" + catID
    }
    return v, c.get(c.apiURL("get_vod_streams", extra), &v)
}

func (c *Client) GetSeriesCategories() ([]Category, error) {
    var v []Category
    return v, c.get(c.apiURL("get_series_categories"), &v)
}

func (c *Client) GetSeries(catID string) ([]Series, error) {
    var v []Series
    extra := ""
    if catID != "" {
        extra = "category_id=" + catID
    }
    return v, c.get(c.apiURL("get_series", extra), &v)
}

func (c *Client) GetSeriesInfo(seriesID int) (*SeriesInfo, error) {
    var v SeriesInfo
    extra := fmt.Sprintf("series_id=%d", seriesID)
    return &v, c.get(c.apiURL("get_series_info", extra), &v)
}

func (c *Client) LiveStreamURL(id int, ext string) string {
    if ext == "" { ext = "ts" }
    return fmt.Sprintf("%s/live/%s/%s/%d.%s", c.BaseURL, c.Username, c.Password, id, ext)
}

func (c *Client) VODStreamURL(id int, ext string) string {
    if ext == "" { ext = "mp4" }
    return fmt.Sprintf("%s/movie/%s/%s/%d.%s", c.BaseURL, c.Username, c.Password, id, ext)
}

func (c *Client) SeriesStreamURL(episodeID, ext string) string {
    if ext == "" { ext = "mkv" }
    return fmt.Sprintf("%s/series/%s/%s/%s.%s", c.BaseURL, c.Username, c.Password, episodeID, ext)
}