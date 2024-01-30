package utils

import "regexp"

type RegExUtils struct{}

var RegExUtilsInstance RegExUtils

func (ru RegExUtils) FindNamedMatches(regex *regexp.Regexp, str string) map[string]string {
	match := regex.FindStringSubmatch(str)

	results := map[string]string{}

	for i, name := range match {
		results[regex.SubexpNames()[i]] = name
	}

	return results
}

func (rp *RegExUtils) HasMatch(regexs []*regexp.Regexp, str string) bool {
	for _, re := range regexs {
		if re.MatchString(str) {
			return true
		}
	}

	return false
}
