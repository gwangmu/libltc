package libltc 

import (
	"errors"
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

// LTCWarningObj is an interface that any LTC structs should implement in order
// to be referred by LTCWarning.

type LTCWarningObj interface {
	Summary() string
	IsUnknown() bool 
}

//-- LTCWarningObj helper methods

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
		case string:
			if carg != "" {
				fields := append(fields, realprefix + carg)
			}
		case *string:
			if carg != nil {
				fields := append(fields, realprefix + *carg)
			}
		case *int:
			if carg != nil {
				fields := append(fields, realprefix + strconv.Itoa(*carg))
			}
		case LTCWarningObj:
			if carg != nil && !carg.IsUnknown() {
				fields := append(fields, realprefix + carg.Summary())
			}
		}
	}
	return strings.Join(fields, ", ")
}

// LTCWarning is a struct that represents a warning that occured between the
// raw LTC file and the internal representation. It could be used for other
// warnings as well (let's see).

type LTCWarning struct {
	desc string 		// '@<idx>@' to refer to the idx'th object.
	objs []LTCWarningObj
}

//-- struct LTCWarning: method (getters and setters)

func (ltcw *LTCWarning) GetDesc() string {
	return ltcw.getSubstitutedString(ltcw.desc)
}

func (ltcw *LTCWarning) getSubstitutedString(orgstr string) string {
	re := regexp.MustCompile(`@([0-9]+)@`)

	newstr := re.ReplaceAllStringFunc(orgstr, func (match string) string {
		submatch := re.FindStringSubmatch(match)
		if (len(submatch) < 2) {
			return match
		}

		idx, err := strconv.Atoi(submatch[1])
		if err != nil || idx >= len(ltcw.objs) {
			return match
		}

		return ltcw.objs[idx].Summary()
	})

	return strings.TrimSpace(newstr)
}

//-- struct LTCWarning: method (creation)

func CreateWarning(fmtstr string, args ...interface{}) LTCWarning {
	fmtargs := []interface{}{}
	objs := []LTCWarningObj{}

	for _, arg := range args {
		switch carg := arg.interface{}.(type) {
		case LTCWarningObj:
			objs = append(objs, carg)
		default:
			fmtargs = append(fmtargs, carg)
		}
	}

	return LTCWarning{
		desc: fmt.Sprintf(fmtstr, fmtargs...),
		objs: objs,
	}
}
