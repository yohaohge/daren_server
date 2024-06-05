package store

import (
	"strconv"
	"time"
)

type VideoInfo struct {
	Id    string `json:"id" db:"id"`
	Name  string `json:"name" db:"name"`
	Data  string `json:"data" db:"data"`
	Cover string `json:"cover" db:"cover"`
	Total int    `json:"total" db:"total"`
	Desc  string `json:"desc" db:"desc"`
	Label string `json:"label" db:"label"`
}

type VideoEpisodeInfo struct {
	Num      int    `json:"num""`
	SubCover string `json:"sub_cover""`
	PlayUrl  string `json:"play_url"`
}

type VideoDetailInfo struct {
	*VideoInfo
	VData map[int]VideoEpisodeInfo `json:"v_data"`
}

func GetVideoDetail(id int) *VideoDetailInfo {
	key := "VideoDetail_" + strconv.Itoa(id)
	val, b := C.Get(key)
	if b && val != nil {
		vd, ok := val.(*VideoDetailInfo)
		if ok {
			return vd
		}
	}
	item := MC.getVideoDetail(id)
	if item != nil {
		C.Add(key, item, time.Hour*2)
	}
	return item
}
