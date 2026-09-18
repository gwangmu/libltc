package libltc

import (
	"errors"
	"regexp"
	"strconv"
	"strings"
)

// LTCWarning is a struct that represents a warning that occured between the
// raw LTC file and the internal representation. It could be used for other
// warnings as well (let's see).

type LTCWarning struct {
	Descs []string 			// '@<idx>@' to refer to the idx'th object.
	Measures []string		// '@<idx>@' to refer to the idx'th object.
	Objects []LTCWarningObj
}

// LTCWarningObj is an interface that any LTC structs should implement in order
// to be referred by LTCWarning.

type LTCWarningObj interface {
	Summary() string
	IsUnknown() bool 
}

func toSentenceCase(s string) string {
	if (len(s) < 1) {
		return ""
	} else {
		return strings.ToUpper(s[:1]) + s[1:]
	}
}

func getWarningInfoString(args ...struct {string; interface{}}) string {
	fields := []string{}
	for _, arg := range args {
		realprefix = arg.string + " "
		if arg.prefix == "" {
			realprefix = ""
		}

		switch carg := arg.interface{}.(type) {
		string:
			if carg != "" {
				fields := append(fields, realprefix + carg)
			}
		*int:
			if carg != nil {
				fields := append(fields, realprefix + strconv.Itoa(carg))
			}
		LTCWarningObj:
			if carg != nil && !carg.IsUnknown() {
				fields := append(fields, realprefix + carg.Summary())
			}
		}
	}
	return strings.Join(fields, ", ")
}

func (ltcw LTCWarning) getSubstitutedString(strs []string) (ret string) {
	re := regexp.MustCompile(`@([0-9]+)@`)

	for desc := range strs {
		newstr := re.ReplaceAllStringFunc(desc, func (match string) string {
			submatch := re.FindStringSubmatch(match)
			if (len(submatch) < 2) {
				return match
			}

			idx, err := strconv.Atoi(submatch[1])
			if err != nil || idx >= len(ltcw.Objects) {
				return match
			}

			return ltcw.Objects[idx].Summary()
		})

		ret += toSentenceCase(newstr) + '\n'
	}

	ret = strings.TrimSpace(ret)
	return
}

func (ltcw *LTCWarning) GetDesc() string {
	return ltcw.getSubstitutedString(ltcw.Descs)
}

func (ltcw *LTCWarning) GetMeasure() string {
	return ltcw.getSubstitutedString(ltcw.Measures)
}
