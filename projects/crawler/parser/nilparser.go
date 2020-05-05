package parser

import "../engine"

type NilParser struct {}

func (NilParser) Parse([] byte) engine.ParseResult {
	return engine.ParseResult{}
}
func (NilParser) Serialize() (name string, args interface{}){
	return "NilParser",nil
}