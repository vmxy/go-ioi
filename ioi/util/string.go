package util

import (
	"regexp"
	"strconv"
	"strings"
)

type String string
type Strings []string

var Reg_Trim *regexp.Regexp = regexp.MustCompile(`^[^\w]+|[^\w]+$`)

func (str String) Trim() String {
	v := Reg_Trim.ReplaceAllString(string(str), "")
	return String(v)
}

func (str String) String() string {
	return string(str)
}

func (str String) Split(sep string) Strings {
	vs := strings.Split(str.String(), sep)
	return vs
}
func (str String) Join(args ...string) String {
	v := strings.Join(args, str.String())
	return String(v)
}
func (str String) Search(reg string) []string {
	r := regexp.MustCompile(reg)
	list := r.FindAllString(str.String(), -1)
	return list
}
func (str String) Replace(reg string, val string) String {
	r := regexp.MustCompile(reg)
	v := r.ReplaceAllString(str.String(), val)
	return String(v)
}
func (str String) Match(reg string) bool {
	r := regexp.MustCompile(reg)
	return r.MatchString(str.String())
}
func (str String) Includes(v string) bool {
	return strings.Contains(str.String(), v)
}
func (str String) ToLower() String {
	v := strings.ToLower(str.String())
	return String(v)
}
func (str String) ToUpper() String {
	v := strings.ToUpper(str.String())
	return String(v)
}
func (str String) ToInt() int {
	v, err := strconv.Atoi(str.String())
	if err != nil {
		return -1
	}
	return v
}
func (str String) ToUInt8() uint8 {
	v, err := strconv.Atoi(str.String())
	if err != nil {
		return 0
	}
	return uint8(v)
}

func (str String) Slice(start int, end int) String {
	if end < 0 {
		end = len(str)
	}
	if end < start {
		end = start
	}
	v := str.String()[start:end]
	return String(v)
}
func (str String) Size() int {
	return len(str)
}
func (list Strings) Map(handle func(val String, index int) string) Strings {
	as := make(Strings, len(list))
	for i, v := range list {
		val := String(v)
		as[i] = handle(val, i)
	}
	return as
}
func (list Strings) String() string {
	return string(String("").Join(list...))
}
func (list Strings) Get(index int) String {
	if index >= len(list) {
		return ""
	}
	return String(list[index])
}
func (list Strings) GetV(index int, defaultValue string) String {
	if index >= len(list) {
		return String(defaultValue)
	}
	v := list.Get(index)
	if v == "" {
		return String(defaultValue)
	}
	return v
}
func (list *Strings) Append(vals ...string) int {
	*list = append(*list, vals...)
	return len(*list)
}

func (list Strings) Join(sep string) String {
	v := strings.Join(list, sep)
	return String(v)
}

func (list Strings) Slice(start int, end int) Strings {
	if end < 0 {
		end = len(list)
	}
	if end < start {
		end = start
	}
	v := list[start:end]
	return v
}
