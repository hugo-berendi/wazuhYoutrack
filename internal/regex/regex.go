package regex

import "regexp"

func UseRegex(s string, reg string) bool {
	re := regexp.MustCompile("(?i)" + reg)
	return re.MatchString(s)
}

func GetString(s string, reg string) string {
	re := regexp.MustCompile(reg)
	return re.FindString(s)
}

func GetStringSubmatch(s string, reg string) []string {
	re := regexp.MustCompile(reg)
	return re.FindStringSubmatch(s)
}
