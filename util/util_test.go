package util

import "testing"

func TestIsBlockMediaName(t *testing.T) {
	res := IsBlockMediaName("新京报")
	if !res {
		t.Error("新京报 not ok")
	}
	res = IsBlockMediaName("中国新闻网")
	if !res {
		t.Error("中国新闻网 not ok")
	}
	res = IsBlockMediaName("中央纪委")
	if !res {
		t.Error("中央纪委国家监委 not ok")
	}
	res = IsBlockMediaName("国家监委")
	if !res {
		t.Error("国家监委 not ok")
	}
	res = IsBlockMediaName("中国善力量")
	if !res {
		t.Error("中国善力量 not ok")
	}
}
