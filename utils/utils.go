package utils

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
)

type RedisSetting struct {
	Host string `json:"host"`
	Port int    `json:"port"`
}

type CrontabSetting struct {
	RedisPersistenceFrequency string `json:"redis_persistence_frequency"` // crontab定时持久化到数据库的频率
	RedisCheckCacheFrequency  string `json:"redis_check_cache_frequency"` // crontab定时检测清除没有更新的cache频率
	RedisClearCacheTimeout    string `json:"redis_clear_cache_timeout"`   // 清除超过多久没有更新的cache
}

type Settings struct {
	RedisSettings   RedisSetting   `json:"redis"`
	CrontabSettings CrontabSetting `json:"crontab"`
}

func GetSettings() Settings {
	data, err := ioutil.ReadFile("settings.json")
	if err != nil {
		fmt.Println("error happened when reading settings.json")
		panic(err)
	}

	var settings Settings
	_ = json.Unmarshal(data, &settings)
	return settings
}
