package playback

import (
	"os"

	"github.com/erwijet/spotipub/internal/logging"
	"gopkg.in/yaml.v3"
)

type ShameList struct {
	Albums  []string `yaml:albums`
	Artists []string `yaml:artists`
	Items   []string `yaml:items`
}

func NewShameList() ShameList {
	return ShameList{
		Albums:  make([]string, 0),
		Artists: make([]string, 0),
		Items:   make([]string, 0),
	}
}

func (s *ShameList) Load(name string) *ShameList {
	log := logging.GetLogger("ShameList")
	file, err := os.ReadFile(name)
	if err != nil {
		log.Printf("error loading shamelist '%s', #%v", name, err)
	}

	err = yaml.Unmarshal(file, s)
	if err != nil {
		log.Fatalf("Unmarshal error: %v", err)
	}

	return s
}
