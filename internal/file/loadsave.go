package file

import (
	"regexp"

	"github.com/BurntSushi/toml"
	"github.com/gwangmu/libltc/warning"
)

func LoadFromString[T IFileObject, U IDiagnosablePtr[T]](tomlstr string) (U, warning.Warnings, error) {
	var out U = new(T)
	warns := warning.Warnings{}

	if meta, err := toml.Decode(tomlstr, out); err != nil {
		warns.Add("cannot read TOML format.")
		return nil, warns, err
	} else {
		for _, ukey := range meta.Undecoded() {
			warns.Add("cannot decode a field '%s' in the LTC file.", ukey)
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

func SaveToString[T IFileObject, U IDiagnosablePtr[T]](fobj U) (string, warning.Warnings, error) {
	// TODO
	panic("Unimplemented")
}

func SaveToURI[T IFileObject, U IDiagnosablePtr[T]](uri string, fobj U) (warning.Warnings, error) {
	// TODO
	panic("Unimplemented")
}
