package parser

import (
	"../engine"
	"../model"
	"github.com/bitly/go-simplejson"
	"log"
	"regexp"
	"strconv"
)

type ProfileParser struct {
	Name string
	Url string
}


var re = regexp.MustCompile(`<script>window.__INITIAL_STATE__=(.+);\(function`)
func (p ProfileParser) Parse(contents []byte) engine.ParseResult {
	match := re.FindSubmatch(contents)
	result :=engine.ParseResult{}
	if len(match) >= 2 {
		profile, id :=parseJson(match[1])
		profile.Name=p.Name
		result.Items =append(result.Items,engine.Item{
			Url:p.Url,
			Type:"zhanai:profile",
			Id:id,
			Payload:profile,
		})
	}

	return result
}

func (p ProfileParser) Serialize() (name string, args interface{}){
	return "ProfileParser",[]string{p.Url,p.Name}
}
//解析json数据
func parseJson(json []byte) (model.Profile,string) {
	res, err := simplejson.NewJson(json)
	if err != nil {
		log.Println("解析json失败。。")
	}
	infos, err := res.Get("objectInfo").Get("basicInfo").Array()
	//infos是一个切片，里面的类型是interface{}
	var profile model.Profile
	//所以我们遍历这个切片，里面使用断言来判断类型
	for k, v := range infos {
		/*
		 "basicInfo":[
			"未婚",
			"25岁",
			"魔羯座(12.22-01.19)",
			"152cm",
			"42kg",
			"工作地:阿坝茂县",
			"月收入:3-5千",
			"医生",
			"大专"
		],
		 */
		if e, ok := v.(string); ok {
			switch k {
			case 0:
				profile.Marriage = e
			case 1:
				//年龄:47岁，我们可以设置int类型，所以可以通过另一个json字段来获取
				profile.Age = e
			case 2:
				profile.Xingzuo = e
			case 3:
				profile.Height = e
			case 4:
				profile.Weight = e
			case 6:
				profile.Income = e
			case 7:
				profile.Occupation = e
			case 8:
				profile.Education = e
			}
		}

	}
	//id
	id, err := res.Get("objectInfo").Get("memberID").Int()
	return profile,strconv.Itoa(id)
}



