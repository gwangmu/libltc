package file

import (
	"bytes"
	"errors"
	"io"
	"os"
	"regexp"
	"strings"
	"net/http"

	"github.com/BurntSushi/toml"
	"libltc/warning"
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

func readFileFromURI(uri string) (ret string, local bool, err error) {
	protRe := regexp.MustCompile(`^https?://`)
	if protRe.MatchString(uri) {
		resp, err := http.Get(uri)
		if err != nil {
			return "", false, err
		} else {
			defer resp.Body.Close()
			bs, err := io.ReadAll(resp.Body)
			if err != nil {
				return "", false, err
			} else {
				return string(bs), false, nil
			}
		}
	} else {
		bs, err := os.ReadFile(uri)
		if err != nil {
			return "", true, err
		} else {
			return string(bs), true, nil
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
		//fmt.Println(strconv.Itoa(i) + " : " + fileline)

		if len(fileline) > 1 && fileline[0:1] == "[" {
			idxObjEnd = i-1
			if idxID != -1 {
				break
			}
			idxObjHead = i+1
			ouri = ""
			continue
		}

		reID := regexp.MustCompile(`^ID\s*=\s*"([a-zA-Z0-9]*)"`)
		if matches := reID.FindStringSubmatch(fileline); matches != nil {
			if matches[1] == idstr {
				idxID = i
				continue
			}
		}

		reLink := regexp.MustCompile(`^(?:Link|EmbedChart|EmbedEvent)\s*=\s*"([.*])"`)
		if matches := reLink.FindStringSubmatch(fileline); matches != nil {
			ouri = matches[1]
			continue
		}
	}

	if idxObjEnd == -1 {
		idxObjEnd = len(filelines)
	}

	if idxID == -1 {
		return "", "", errors.New("Object '" + idstr + "' not found.")
	}
	if idxObjHead == -1 || idxObjHead > idxObjEnd {
		//fmt.Println(idxID, idxObjHead, idxObjEnd)
		return "", "", errors.New("Cannot find the enclosing object of ID '" + idstr + "'")
	}

	ostr := strings.Join(filelines[idxObjHead:idxObjEnd], "\n")
	return ostr, ouri, nil
}

func readObjectFromURI(uri string) (ret string, lpath string, err error) {
	// If `uri` contains ".obj/", the trailing part should be considered 
	// an object ID. Find that id from the preceding path. Otherwise, 
	// simply read the file and return.

	path, qualid, seped := strings.Cut(uri, "?id=")
	filestr, local, err := readFileFromURI(path)
	if err != nil {
		return "", path, err
	} else {
		ret = filestr
		if local {
			lpath = path
		}
	}

	if seped {
		ids := strings.Split(qualid, "/")
		for _, idstr := range ids {
			if ostr, ouri, err := extractObjectByLocalID(filestr, idstr); err == nil {
				if len(ouri) != 0 {
					if newOstr, nlpath, err := readObjectFromURI(ouri); err == nil {
						lpath = nlpath
						ostr = newOstr
					} else {
						return "", path, err
					}
				}
				ret = ostr
			} else {
				return "", path, err
			}
		}
	}

	return ret, path, nil
}

func LoadFromURI[T IFileObject, U IDiagnosablePtr[T]](uri string) (U, warning.Warnings, error) {
	bs, lpath, err := readObjectFromURI(uri)
	if err != nil {
		warns := warning.Warnings{}
		warns.Add("cannot read a local file.")
		return nil, warns, err
	} else {
		var cwdpath string
		if lpath != "" {
			if maybeCwdpath, err := os.Getwd(); err != nil {
				cwdpath = maybeCwdpath
			} 
			os.Chdir(lpath)
		}
		ret, warns, err := LoadFromString[T, U](string(bs))
		if lpath != "" {
			os.Chdir(cwdpath)
		}
		return ret, warns, err
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

	protRe := regexp.MustCompile(`^https?://`)
	if protRe.MatchString(uri) {
		// TODO: probably an online resource. upload TOML string.
		// TODO: send requests after '?'
		warns.Add("saving to an online URI is not supported yet.")
		return warns, errors.New("attempted to save to an online URI")
	} else {
		err = os.WriteFile(uri, []byte(ltcstr), 0600)
		if err != nil {
			warns.Add("cannot write to a local file at '" + uri + "'.")
			return warns, err
		}
	}

	return warns, err
}
