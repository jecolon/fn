// Package fn has utility functions for fixing filenames according to the ideas
// presented in http://www.dwheeler.com/essays/fixing-unix-linux-filenames.html
package fn

import (
	"regexp"
	"strings"
	"unicode"

	ud "github.com/fiam/gounidecode/unidecode"
)

// SPECIAL characters that cause problems in various contexts.
const SPECIAL = "[]{}()<>;:\\|/*?¿¡\"'~!@#$%^&=+"

// FixForShell rewrites a string so that it qualifies as a "friendly" filename for shells.
// This tries to adhere to http://www.dwheeler.com/essays/fixing-unix-linux-filenames.html
func FixForShell(source string) (result string) {
	result = stripControl(source)
	result = stripSpecial(result)
	result = replaceSpaces(result, "_")
	// Remove duplicate, leading, or trailing dash, underscore, or space
	result = trim(result, "-_ ")
	return
}

// FixForURL rewrites a string so that it qualifies as a "friendly" filename for URLs.
// This tries to adhere to http://www.dwheeler.com/essays/fixing-unix-linux-filenames.html
func FixForURL(source string) (result string) {
	result = stripControl(source)
	result = stripSpecial(result)
	result = replaceSpaces(result, "_")
	result = ud.Unidecode(result)
	// Remove duplicate, leading, or trailing dash, underscore, or space
	result = trim(result, "-_ ")
	return
}

func stripControl(source string) (result string) {
	for _, r := range source {
		if unicode.IsControl(r) {
			continue
		}
		result += string(r)
	}
	return
}

func stripSpecial(source string) (result string) {
	for _, r := range source {
		// rune to string
		s := string(r)
		if strings.Contains(SPECIAL, s) {
			continue
		}
		result += string(s)
	}
	return
}

func replaceSpaces(source, replacement string) (result string) {
	for _, r := range source {
		if r == ' ' {
			result += replacement
			continue
		}
		result += string(r)
	}
	return
}

func trim(source, set string) (result string) {
	parts := strings.Split(source, ".")
	for i := range parts {
		for _, r := range set {
			character := string(r)
			parts[i] = regexp.MustCompile("["+character+"]{2,}").ReplaceAllString(parts[i], character)
		}
		re := regexp.MustCompile("(^[" + set + "]|[" + set + "]$)")
		for re.MatchString(parts[i]) {
			parts[i] = re.ReplaceAllString(parts[i], "")
		}
	}
	result = strings.Join(parts, ".")
	return
}
