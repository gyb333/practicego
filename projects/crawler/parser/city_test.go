package parser

import (
	"testing"
	"io/ioutil"
)

func TestCityParser_Parse(t *testing.T) {
	contents,err:=ioutil.ReadFile("city_test_data.html")
	if err != nil{
		panic(err)
	}

	parseResult := CityParser{}.Parse(contents)
	const resultSize = 20
	expectedUrls := []string{
		"http://album.zhenai.com/u/1280064210",
		"http://album.zhenai.com/u/1312667054",
		"http://album.zhenai.com/u/1321413610",
	}
	if len(parseResult.Requests) != resultSize{
		t.Errorf("result should have %d requests, but had %d",resultSize,len(parseResult.Requests))
	}
	if len(parseResult.Items) != resultSize{
		t.Errorf("result should have %d Items, but had %d",resultSize,len(parseResult.Items))
	}
	for i, url := range expectedUrls {
		if parseResult.Requests[i].Url != url {
			t.Errorf("expected url #%d: %s; but "+"was %s",
				i, url, parseResult.Requests[i].Url)
		}
	}

}
