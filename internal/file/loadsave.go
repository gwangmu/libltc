package file

import (
	"bytes"
	"os"
	"regexp"

	"github.com/BurntSushi/toml"
	"github.com/gwangmu/libltc/warning"
)

func LoadFromString[T IFileObject, U IDiagnosablePtr[T]](ltcstr string) (U, warning.Warnings, error) {
	var out U = new(T)
	warns := warning.Warnings{}

	if meta, err := toml.Decode(ltcstr, out); err != nil {
		warns.Add("cannot decode TOML format.")
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
	var ltcstr string

	protRe := regexp.MustCompile(`^[a-z]+://`)
	if protRe.MatchString(uri) {
		// TODO: probably an online resource. download TOML string.
		// TODO: send requests after '?'
		// ltcstr = ...
		panic("Unimplemented")
	} else {
		bs, err := os.ReadFile(uri)
		if err != nil {
			warns := warning.Warnings{}
			warns.Add("cannot read a local file.")
			return nil, warns, err
		} else {
			ltcstr = string(bs)
		}
	}

	return LoadFromString[T, U](ltcstr)
}

func SaveToString[T IFileObject, U IDiagnosablePtr[T]](fobj U) (string, warning.Warnings, error) {
	warns := warning.Warnings{}

	moreWarns := fobj.Diagnose()
	warns.Concat(moreWarns)

	out := new(bytes.Buffer)
	err := toml.NewEncoder(out).Encode(fobj)
	if err != nil {
		warns.Add("cannot encode TOML format.")
		return "", warns, err
	}

	return out.String(), warns, nil
}

func SaveToURI[T IFileObject, U IDiagnosablePtr[T]](uri string, fobj U) (warning.Warnings, error) {
	ltcstr, warns, err := SaveToString[T, U](fobj)
	if err != nil {
		warns.Add("unsaved due to errors.")
		return warns, err
	}

	protRe := regexp.MustCompile(`^[a-z]+://`)
	if protRe.MatchString(uri) {
		// TODO: probably an online resource. upload TOML string.
		// TODO: send requests after '?'
		// ltcstr = ...
		panic("Unimplemented")
	} else {
		err = os.WriteFile(uri, []byte(ltcstr), 0600)
		if err != nil {
			warns.Add("cannot write to a local file at '" + uri + "'.")
			return warns, err
		}
	}

	return warns, err
}
