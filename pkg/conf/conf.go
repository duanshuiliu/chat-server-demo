package conf

import (
	"errors"
	"io/ioutil"
	"strings"

	"chat/pkg/conf/parsers"
)

var Instance *Conf

var DefaultConfigPath string = "./config"

var ErrEmptyInstance = errors.New("not found instance")
var ErrEmptyFile     = errors.New("not found config file")

type Conf struct {
	ConfMap map[string]ConfParser
}

type ConfParser interface {
	Load(string) error
	Reload() error
	All() interface{}
	String(string) (string, error)
	Int(string) (int, error)
}

func Register(confPath string) (err error) {
	if confPath == "" {
		confPath = DefaultConfigPath
	}

	conf_files, err := ioutil.ReadDir(confPath)

	if err != nil {
		return
	}

	Instance = &Conf{
		ConfMap: make(map[string]ConfParser),
	}

	s := make(map[string]string)

	for _, file := range conf_files {
		if !file.IsDir() {
			filename           := file.Name()
			filepath           := confPath+"/"+filename
			filenameTrimSuffix := strings.Split(filename, ".")[0] 

			s[filenameTrimSuffix] = filepath
		}
	}

	for name := range s {
		// TODO 工厂模式：这里switch判定是哪个parser
		parser := &parsers.IniParser{}
		
		if loadErr := parser.Load(s[name]); loadErr != nil {
			err = loadErr
			return
		}

		Instance.ConfMap[name] = parser
	}

	return
}

func New(filename string) (parser ConfParser, err error) {
	if Instance == nil {
		err = ErrEmptyInstance
		return
	}

	if _, ok := Instance.ConfMap[filename]; !ok {
		err = ErrEmptyFile
		return
	}

	parser = Instance.ConfMap[filename]
	return
}

func Reload(filename string) {
	// 
}
