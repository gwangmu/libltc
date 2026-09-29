package file

import (
	"regexp"

	"github.com/BurntSushi/toml"
	"github.com/gwangmu/libltc/warning"
)

type LTCFileObject interface {
	File | Setting | Subject | Event | Annex | Import | Name | Time
}

func LoadFromString[T LTCFileObject](tomlstr string) (*T, warning.Warnings, error) {
	warns := warning.Warnings{}

	var out T
	if _, err := toml.Decode(tomlstr, &out); err != nil {
		warns.Add("cannot read TOML.")
		return nil, warns, err
	}

	//warns := inspectFile(&ltcFile, meta) 
	// TODO: Auto-fix some issues and report via Warning.

	return &out, warns, nil
}

func LoadFromURI[T LTCFileObject](uri string) (*T, warning.Warnings, error) {
	var tomlstr string

	protRe := regexp.MustCompile(`^[a-z]+://`)
	if protRe.MatchString(uri) {
		// TODO: probably an online resource. download TOML string.
		// TODO: send requests after '?'
		// tomlstr = ...
		panic("Unimplemented")
	} else {
		// TODO: assume this is a local file. read TOML file.
		// tomlstr = ...
		panic("Unimplemented")
	}

	return LoadFromString[T](tomlstr)
}
