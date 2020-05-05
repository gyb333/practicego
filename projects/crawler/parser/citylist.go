package parser

import (
	"../engine"
	"regexp"
)


const cityListRe = `<a href="(http://www.zhenai.com/zhenghun/([0-9a-z]+))" [^>]*>([^<]+)</a>`
var expr=`<a (target="_blank" )?href="(http://www.zhenai.com/zhenghun/[0-9a-z]+)" data-v-[0-9a-z]{8}>([^<]+)</a>`
type CityListParser struct {

}

func (p CityListParser) Parse (contents []byte) engine.ParseResult  {
	re := regexp.MustCompile(cityListRe)
	all := re.FindAllSubmatch(contents, -1)

	result := engine.ParseResult{}
	for _, c := range all {
		url :=string(c[1])
		//cityID:=string(c[2])
		//result.Items = append(result.Items,
		//	engine.Item{
		//		Url:url,
		//		Type :"zhenai:citylist",
		//		Id :cityID,
		//		Payload :string(c[3]),
		//	},
		//	) //城市名字
		result.Requests = append(result.Requests, engine.Request{
			Url:     url    ,
			Parser:  CityParser{},
		})
	}
	return result
}


func (p CityListParser) Serialize() (name string, args interface{}){
	return "CityListParser",nil
}


