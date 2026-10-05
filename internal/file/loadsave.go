package file

import (
	"bytes"
	"errors"
	"os"
	"regexp"
	"strings"

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

func readFileFromURI(uri string) (string, error) {
	protRe := regexp.MustCompile(`^[a-z]+://`)
	if protRe.MatchString(uri) {
		// TODO: probably an online resource. download TOML string.
		panic("Unimplemented")
	} else {
		bs, err := os.ReadFile(uri)
		if err != nil {
			return "", err
		} else {
			return string(bs), nil
		}
	}
}

func extractObjectByLocalID(filestr string, idstr string) (string, string, error) {
	filelines := strings.Split(filestr, "\n")

	idxID := -1
	idxObjHead := -1
	idxObjEnd := -1
	ouri := ""

	for i, fileline := range filelines {
		fileline = strings.TrimSpace(fileline)

		if len(fileline) > 1 && fileline[0:1] == "[" {
			idxObjEnd = i-1
			idxObjHead = i+1
			ouri = ""
			if idxID != -1 {
				break
			}
			continue
		}

		reID := regexp.MustCompile(`^ID\s*=\s*"` + idstr + `"`)
		if matches := reID.FindStringSubmatch(fileline); matches != nil {
			idxID = i
			continue
		}

		reLink := regexp.MustCompile(`^(?:Link|EmbedChart|EmbedEvent)\s*=\s*"([.*])"`)
		if matches := reLink.FindStringSubmatch(fileline); matches != nil {
			ouri = matches[1]
			continue
		}
	}

	if idxID == -1 {
		return "", "", errors.New("Object '" + idstr + "' not found.")
	}
	if idxObjHead == -1 {
		return "", "", errors.New("Cannot find the enclosing object of ID '" + idstr + "'")
	}
	if idxObjEnd == -1 {
		idxObjEnd = len(filelines)
	}

	ostr := strings.Join(filelines[idxObjHead:idxObjEnd], "\n")
	return ostr, ouri, nil
}

func readObjectFromURI(uri string) (string, error) {
	// If `uri` contains ".obj/", the trailing part should be considered 
	// an object ID. Find that id from the preceding path. Otherwise, 
	// simply read the file and return.

	ret := ""
	path, qualid, seped := strings.Cut(uri, "?id=")
	filestr, err := readFileFromURI(path)
	if err != nil {
		return "", err
	}

	if seped {
		ids := strings.Split(qualid, "/")
		for _, idstr := range ids {
			if ostr, ouri, err := extractObjectByLocalID(filestr, idstr); err == nil {
				if len(ouri) != 0 {
					if newOstr, err := readObjectFromURI(ouri); err == nil {
						ostr = newOstr
					} else {
						return "", err
					}
				}
				ret = ostr
			} else {
				return "", err
			}
		}
	}

	return ret, nil
}

func LoadFromURI[T IFileObject, U IDiagnosablePtr[T]](uri string) (U, warning.Warnings, error) {
	bs, err := readObjectFromURI(uri)
	if err != nil {
		warns := warning.Warnings{}
		warns.Add("cannot read a local file.")
		return nil, warns, err
	} else {
		return LoadFromString[T, U](string(bs))
	}
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

// For now, saving non-file to URI is unsupported.
func SaveToURI[T interface{File}, U IDiagnosablePtr[T]](uri string, fobj U) (warning.Warnings, error) {
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
