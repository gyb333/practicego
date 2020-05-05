package worker

import (
	"../config"
	"../engine"
	"../parser"
	"errors"
	"fmt"
	"log"
)

type SerializedParser struct {
	Name string //functionName
	Args interface{}
}

type Request struct {
	Url    string
	Parser SerializedParser
}

type ParseResult struct {
	Items    []engine.Item
	Requests []Request
}

func SerializeRequest(r engine.Request) Request {
	name, args := r.Parser.Serialize()
		return Request{
		Url: r.Url,
		Parser: SerializedParser{
			Name: name,
			Args: args,
		},
	}
}

func DeserializeRequest(r Request) (engine.Request, error) {

	parser, err := deserializeParser(r.Parser)

	if err != nil {
		return engine.Request{}, err
	}

	return engine.Request{
		Url:    r.Url,
		Parser: parser,
	}, nil
}

func SerializeParseResult(r engine.ParseResult) ParseResult {
	result := ParseResult{
		Items: r.Items,
	}
	for _, req := range r.Requests {
		result.Requests = append(result.Requests, SerializeRequest(req))
	}

	return result
}

func DeserializeResult(r ParseResult) engine.ParseResult {
	result := engine.ParseResult{
		Items: r.Items,
	}
	for _, req := range r.Requests {
		engineReq, err := DeserializeRequest(req)
		if err != nil {
			log.Printf("error desrializing request : %v", err)
			continue
		}
		result.Requests = append(result.Requests, engineReq)
	}
	return result
}

func deserializeParser(p SerializedParser) (engine.Parser, error) {
	switch p.Name {
	case config.CityListParser:
		return parser.CityListParser{}, nil
	case config.CityParser:
		return parser.CityParser{}, nil
	case config.NilParser:
		return parser.NilParser{}, nil
	case config.ProfileParser:
		if args, ok := p.Args.([]interface {}); ok {
			url,_:=args[0].(string)
			name,_:=args[1].(string)
			return parser.ProfileParser{
				Url:  url,
				Name: name,
			}, nil
		} else {
			return nil, fmt.Errorf("invalid arg : %v", p.Args)
		}
	default:
		return nil, errors.New("unknown parser name")
	}

}
