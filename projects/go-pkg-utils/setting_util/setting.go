package setting_util

import (
	"github.com/go-ini/ini"
	"log"
)

type Setting struct {
	*ini.File
}

func NewSetting(path string) *Setting{
	cfg, err := ini.Load(path)
	if err != nil {
		log.Fatalf("Fail to parse '%s': %v", path,err)
	}
	return &Setting{cfg}
}

// mapTo map section
func (cfg *Setting) MapTo(section string, v interface{}) {
	err := cfg.Section(section).MapTo(v)
	if err != nil {
		log.Fatalf("Cfg.MapTo %s err: %v", section, err)
	}
}


