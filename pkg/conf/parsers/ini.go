package parsers

import (
	"errors"
	"github.com/go-ini/ini"
	"strings"
)

var ErrEmptyReader = errors.New("not found parser reader")

type IniParser struct {
	reader *ini.File	
}

func (this *IniParser) Load(config_file string) (err error) {
	conf, err := ini.Load(config_file)

	if err != nil {
		this.reader = nil
		return
	}

	this.reader = conf
	return 
}

func (this *IniParser) Reload() (err error) {
	return this.reader.Reload()
}

func (this *IniParser) All() interface{} {
	sections := this.reader.Sections()

	result := make(map[string]map[string]string)

	for _, section := range sections {
		result[section.Name()] = section.KeysHash()
	}
	
	return result
}

func (this *IniParser) String(key string) (value string, err error) {
	if this.reader == nil {
		err = ErrEmptyReader
		return
	}

	sectionStr, keyStr := this.GetSecionAndKey(key)

	iniSection := this.reader.Section(sectionStr)
	if iniSection == nil {
		err = errors.New("not found section: "+sectionStr)
		return
	}

	value = iniSection.Key(keyStr).String()	
	return
}

func (this *IniParser) Int(key string) (value int, err error) {
	if this.reader == nil {
		err = ErrEmptyReader
		return
	}

	sectionStr, keyStr := this.GetSecionAndKey(key)

	iniSection := this.reader.Section(sectionStr)
	if iniSection == nil {
		err = errors.New("not found section: "+sectionStr)
		return
	}

	value, err = iniSection.Key(keyStr).Int()
	return
}

func (this *IniParser) GetSecionAndKey(key string) (sectionStr, keyStr string) {
	s := strings.Split(key, "::")

	switch (len(s)) {
		case 1:
			keyStr = s[0]
		case 2:
			sectionStr = s[0]
			keyStr     = s[1]
	}

	return
}