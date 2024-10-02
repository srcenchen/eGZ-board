package service

import (
	"eGZ-Board/internal/dao"
	"eGZ-Board/utility"
	"fmt"
)

func WeatherProxy() {
	const key = "a8a6122bea3c4cf1812f4eabec79f61f"
	const location = "119.15,34.81" // 江苏省赣榆高级中学经济开发区校区
	const baseUrl = "https://devapi.qweather.com"
	LoadWeather("now", fmt.Sprintf("%s/v7/weather/now?location=%s&key=%s", baseUrl, location, key))  // 获取当前天气
	LoadWeather("3d", fmt.Sprintf("%s/v7/weather/3d?location=%s&key=%s", baseUrl, location, key))    // 获取当前天气
	LoadWeather("rain", fmt.Sprintf("%s/v7/minutely/5m?location=%s&key=%s", baseUrl, location, key)) // 获取当前天气
	// LoadWeather("now", fmt.Sprintf("%s/v7/minutely/5m?location=%s&key=%s", baseUrl, location, key))  // 获取当前天气
}

func LoadWeather(key string, url string) {
	body := utility.HttpGet(url)
	dao.SetWeather(key, body)
	//g.Log().Printf(gctx.New(), url)
}
