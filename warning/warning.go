package warning 

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

// Warning is a struct that represents a warning that occured between the
// raw LTC file and the internal representation. It could be used for other
// warnings as well (let's see).

type Warning struct {
	desc string 		// '@<idx>@' to refer to the idx'th object.
	objs []IWarningObj
}

type Warnings []*Warning

//-- struct Warning: method (getters and setters)

func (ltcw *Warning) GetDesc() string {
	return ltcw.getSubstitutedString(ltcw.desc)
}

func (ltcw *Warning) getSubstitutedString(orgstr string) string {
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

func (ws *Warnings) Concat(ws2 Warnings) {
	*ws = append(*ws, ws2...)
}

func (ws *Warnings) AddDirect(wobj *Warning) {
	*ws = append(*ws, wobj)
}

func (ws *Warnings) Add(fmtstr string, args ...any) {
	*ws = append(*ws, Create(fmtstr, args...))
}

//-- struct Warning: method (creation)

func Create(fmtstr string, args ...any) *Warning {
	fmtargs := []interface{}{}
	objs := []IWarningObj{}

	for _, arg := range args {
		switch carg := arg.(type) {
		case IWarningObj:
			objs = append(objs, carg)
		default:
			fmtargs = append(fmtargs, carg)
		}
	}

	return &Warning{
		desc: fmt.Sprintf(fmtstr, fmtargs...),
		objs: objs,
	}
}

// IWarningObj is an interface that any LTC structs should implement in order
// to be referred by Warning.

type IWarningObj interface {
	Summary() string
	IsUnknown() bool 
}

type IDiagnosable interface {
	DiagnoseLocal() []Warning
	DiagnoseNonLocal() []Warning
}

type SummaryElement struct {
	Prefix string
	Object any
}

//-- IWarningObj helper methods

func toSentenceCase(s string) string {
	re := regexp.MustCompile(`(?:^\s*|[\.\?\!]\s*)([a-z])`)
	if mas := re.FindAllStringSubmatchIndex(s, -1); mas != nil {
		for _, ma := range mas {
			s = s[:ma[2]] + string(strings.ToUpper(s[ma[2]:ma[2]+1])) + s[ma[2]+1:]
		}
	}
	return s
}

func BuildSummaryString(args ...SummaryElement) string {
	fields := []string{}
	for _, arg := range args {
		realprefix := arg.Prefix + " "
		if arg.Prefix == "" {
			realprefix = ""
		}

		switch carg := arg.Object.(type) {
		case string:
			if carg != "" {
				fields = append(fields, realprefix + carg)
			}
		case *string:
			if carg != nil {
				fields = append(fields, realprefix + *carg)
			}
		case *int:
			if carg != nil {
				fields = append(fields, realprefix + strconv.Itoa(*carg))
			}
		case IWarningObj:
			if carg != nil && !carg.IsUnknown() {
				fields = append(fields, realprefix + carg.Summary())
			}
		}
	}
	return strings.Join(fields, ", ")
}
