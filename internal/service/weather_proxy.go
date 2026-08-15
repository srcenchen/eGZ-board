package service

import (
	"eGZ-Board/internal/dao"
	"eGZ-Board/utility"
	"encoding/json"
	"fmt"
	"net/url"
	"os"
	"strings"
)

var weatherRefresh = make(chan struct{}, 1)

func WeatherProxy() error {
	key := strings.TrimSpace(os.Getenv("QWEATHER_KEY"))
	if key == "" {
		return nil
	}
	select {
	case weatherRefresh <- struct{}{}:
		defer func() { <-weatherRefresh }()
	default:
		return nil
	}

	location := strings.TrimSpace(os.Getenv("QWEATHER_LOCATION"))
	if location == "" {
		location = "119.15,34.81"
	}
	baseURL := strings.TrimRight(strings.TrimSpace(os.Getenv("QWEATHER_BASE_URL")), "/")
	if baseURL == "" {
		baseURL = "https://devapi.qweather.com"
	}
	parsedBaseURL, err := url.ParseRequestURI(baseURL)
	if err != nil {
		return fmt.Errorf("invalid QWEATHER_BASE_URL: %w", err)
	}
	if (parsedBaseURL.Scheme != "http" && parsedBaseURL.Scheme != "https") || parsedBaseURL.Host == "" {
		return fmt.Errorf("invalid QWEATHER_BASE_URL: absolute HTTP(S) URL required")
	}

	queries := []struct {
		key  string
		path string
	}{
		{key: "now", path: "/v7/weather/now"},
		{key: "3d", path: "/v7/weather/3d"},
		{key: "rain", path: "/v7/minutely/5m"},
	}
	for _, query := range queries {
		values := url.Values{"location": {location}, "key": {key}}
		if err := LoadWeather(query.key, baseURL+query.path+"?"+values.Encode()); err != nil {
			return err
		}
	}
	return nil
}

func LoadWeather(key, endpoint string) error {
	body, err := utility.HttpGet(endpoint)
	if err != nil {
		return fmt.Errorf("load weather %s: %w", key, err)
	}
	var response struct {
		Code string `json:"code"`
	}
	if err := json.Unmarshal([]byte(body), &response); err != nil {
		return fmt.Errorf("decode weather %s response: %w", key, err)
	}
	if response.Code != "200" {
		return fmt.Errorf("weather %s returned provider code %q", key, response.Code)
	}
	if err := dao.SetWeather(key, body); err != nil {
		return fmt.Errorf("store weather %s: %w", key, err)
	}
	return nil
}
