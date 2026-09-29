package file

import (
	"regexp"

	"github.com/BurntSushi/toml"
	"github.com/gwangmu/libltc/warning"
)

type IFileObject interface {
	File | Setting | Subject | Event | Annex | Import | Name | Time
}

type IDiagnosable interface {
	Diagnose() warning.Warnings
}

type IDiagnosablePtr[T any] interface {
	*T
	IDiagnosable
}

func LoadFromString[T IFileObject, U IDiagnosablePtr[T]](tomlstr string) (U, warning.Warnings, error) {
	var out U = new(T)
	warns := warning.Warnings{}

	if meta, err := toml.Decode(tomlstr, out); err != nil {
		warns.Add("cannot read TOML.")
		return nil, warns, err
	} else {
		for _, ukey := range meta.Undecoded() {
			warns.Add("cannot decode a field '%s'.", ukey)
		}
	}

	moreWarns := out.Diagnose() 
	warns.Concat(moreWarns)

	return out, warns, nil
}

func LoadFromURI[T IFileObject, U IDiagnosablePtr[T]](uri string) (U, warning.Warnings, error) {
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

	return LoadFromString[T, U](tomlstr)
}
