package file

import (
	//"fmt"

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

func splitURIIntoPathAndIDs(uri string) (string, []string) {
	re := regexp.MustCompile(`\?id=(.*?)$`)
	if matches := re.FindStringSubmatchIndex(uri); matches != nil {
		path := uri[:matches[0]]
		ids := strings.Split(uri[matches[2]:matches[3]], "/")
		return path, ids
	} else {
		return uri, []string{}
	}
}

func getLocalPath(uri string) string {
	protRe := regexp.MustCompile(`^https?://`)
	if !protRe.MatchString(uri) {
		path, _ := splitURIIntoPathAndIDs(uri)
		return path
	} else {
		return ""
	}
}

func readFileFromURI(uri string) (ret string, err error) {
	if getLocalPath(uri) == "" {
		resp, err := http.Get(uri)
		if err != nil {
			return "", err
		} else {
			defer resp.Body.Close()
			bs, err := io.ReadAll(resp.Body)
			if err != nil {
				return "", err
			} else {
				return string(bs), nil
			}
		}
	} else {
		bs, err := os.ReadFile(uri)
		if err != nil {
			return "", err
		} else {
			return string(bs), nil
		}
	}
}

func extractObjectByLocalID(filestr string, idstr string) (ostr string, olink string, oembed string, err error) {
	filelines := strings.Split(filestr, "\n")

	idxID := -1
	idxObjHead := -1
	idxObjEnd := -1

	for i, fileline := range filelines {
		fileline = strings.TrimSpace(fileline)
		//fmt.Printf("%d : %s\n", i, fileline)

		if len(fileline) > 1 && fileline[0:1] == "[" {
			idxObjEnd = i-1
			if idxID != -1 {
				break
			}
			idxObjHead = i+1
			olink = ""
			oembed = ""
			continue
		}

		reID := regexp.MustCompile(`^ID\s*=\s*"([a-zA-Z0-9]*)"`)
		if matches := reID.FindStringSubmatch(fileline); matches != nil {
			if matches[1] == idstr {
				idxID = i
				continue
			}
		}

		// FIXME: register 'Link' only when it's in an import object.
		reLink := regexp.MustCompile(`^(?:Link)\s*=\s*"(.*?)"`)
		if matches := reLink.FindStringSubmatch(fileline); matches != nil {
			olink = matches[1]
			continue
		}

		// FIXME: register 'Embed*' only when it's in an event object.
		reEmbed := regexp.MustCompile(`^(?:EmbedChart|EmbedEvent)\s*=\s*"(.*?)"`)
		if matches := reEmbed.FindStringSubmatch(fileline); matches != nil {
			oembed = matches[1]
			continue
		}
	}

	if idxObjEnd == -1 || idxObjEnd < idxObjHead {
		idxObjEnd = len(filelines)
	}
	//fmt.Printf("idxID: %d, idxObjHead: %d, idxObjEnd: %d\n", idxID, idxObjHead, idxObjEnd)

	if idxID == -1 {
		err = errors.New("Object '" + idstr + "' not found.")
		return
	}
	if idxObjHead == -1 || idxObjHead > idxObjEnd {
		err = errors.New("Cannot find the enclosing object of ID '" + idstr + "'")
		return
	}

	ostr = strings.Join(filelines[idxObjHead:idxObjEnd], "\n")
	return
}

func readObjectFromURI(uri string) (ret string, err error) {
	// If `uri` contains ".obj/", the trailing part should be considered 
	// an object ID. Find that id from the preceding path. Otherwise, 
	// simply read the file and return.

	path, ids := splitURIIntoPathAndIDs(uri)

	if len(ids) == 0 {
		filestr, err := readFileFromURI(path)
		if err != nil {
			return "", err
		} else {
			ret = filestr
		}
	} else {
		// Sequentially load objects by local IDs.
		curPath := path
		for _, idstr := range ids {
			filestr, err := readFileFromURI(curPath)
			if err != nil {
				return "", err
			} else {
				ret = filestr
			}

			if ostr, olink, oembed, err := extractObjectByLocalID(filestr, idstr); err == nil {
				if len(oembed) != 0 {
					// If a link was detected, load it also.
					if newOstr, err := readObjectFromURI(oembed); err == nil {
						ostr = newOstr
					} else {
						return "", err
					}
				}
				ret = ostr
				curPath = olink
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
		warns.Add("cannot read an object from URI.")
		return nil, warns, err
	} else {
		var cwdpath string
		lpath := getLocalPath(uri)
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
