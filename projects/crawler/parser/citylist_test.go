package parser

import (
	"testing"
	"io/ioutil"
)

func TestCityListParser_Parse(t *testing.T){
	//contents, err := fetcher.Fetch("http://www.zhenai.com/zhenghun")
	contents, err := ioutil.ReadFile("citylist_test_data.html")
	if err != nil{
		panic(err)
	}

	parseResult := CityListParser{}.Parse(contents)

	const resultSize = 470
	expectedUrls := []string{
		"http://www.zhenai.com/zhenghun/aba",
		"http://www.zhenai.com/zhenghun/akesu",
		"http://www.zhenai.com/zhenghun/alashanmeng",
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
