package goutil

import (
	"testing"
)

func Test_main(t *testing.T) {
	t.Run("testCommonMap", testCommonMap)
}

func testCommonMap(t *testing.T) {
	var (
		commonMap *CommonMap
	)
	commonMap = NewCommonMap(1)
	commonMap.AddUniq("1")
	//if len(commonMap.Keys()) != 1 {
	//	t.Error("testCommonMap fail")
	//}
	commonMap.Clear()
	if len(commonMap.Keys()) != 0 {
		t.Error("testCommonMap fail")
	}
	commonMap.AddCount("count", 1)
	commonMap.Add("count")
	if v, ok := commonMap.GetValue("count"); ok {
		if v.(int) != 2 {
			t.Error("testCommonMap fail")
		}
	}
	if !commonMap.Contains("count") {
		t.Error("testCommonMap fail")
	}
	commonMap.Zero()
	if v, ok := commonMap.GetValue("count"); ok {
		if v.(int) != 0 {
			t.Error("testCommonMap fail")
		}
	}
	commonMap.Remove("count")

	if _, ok := commonMap.GetValue("count"); ok {
		t.Error("testCommonMap fail")
	}

}
