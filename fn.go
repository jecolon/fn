// Package fn has utility functions for fixing filenames according to the ideas
// presented in http://www.dwheeler.com/essays/fixing-unix-linux-filenames.html
package fn

import (
	"regexp"
	"strings"
	"unicode"
	"unicode/utf8"

	ud "github.com/fiam/gounidecode/unidecode"
)

// SPECIAL characters that cause problems in various contexts.
const Special = "[]{}()<>;:,\\|/*?¿¡\"'~!@#$%^&=+"

// Maximum length of a filename
const (
	MaxLenShell = 255
	MaxLenURL   = 128
)

// FixForShell rewrites a string so that it qualifies as a "friendly" filename for shells.
// This tries to adhere to http://www.dwheeler.com/essays/fixing-unix-linux-filenames.html
func FixForShell(source string) (result string) {
	result = stripControl(source)
	result = stripSpecial(result)
	result = replaceSpaces(result, '_')
	// Remove duplicate, leading, or trailing dash, underscore, or space
	result = trim(result, "-_ ")
	result = truncate(result, MaxLenShell)
	// Handle all invalid characters
	if result == "" {
		result = "FN_NO_NAME"
	}
	return
}

// FixForURL rewrites a string so that it qualifies as a "friendly" filename for URLs.
// This tries to adhere to http://www.dwheeler.com/essays/fixing-unix-linux-filenames.html
func FixForURL(source string) (result string) {
	result = stripControl(source)
	result = stripSpecial(result)
	result = replaceSpaces(result, '-')
	result = ud.Unidecode(result)
	// Remove duplicate, leading, or trailing dash, underscore, or space
	result = trim(result, "-_ ")
	result = truncate(result, MaxLenURL)
	// Handle all invalid characters
	if result == "" {
		result = "FN-NO-NAME"
	}
	return
}

func stripControl(source string) (result string) {
	var runes []rune
	for _, r := range source {
		if unicode.IsControl(r) {
			continue
		}
		runes = append(runes, r)
	}
	result = string(runes)
	return
}

func stripSpecial(source string) (result string) {
	var runes []rune
	for _, r := range source {
		if strings.ContainsRune(Special, r) {
			continue
		}
		runes = append(runes, r)
	}
	result = string(runes)
	return
}

func replaceSpaces(source string, replacement rune) (result string) {
	var runes []rune
	for _, r := range source {
		if r == ' ' {
			runes = append(runes, replacement)
			continue
		}
		runes = append(runes, r)
	}
	result = string(runes)
	return
}

func trim(source, set string) (result string) {
	// Setup
	parts := strings.Split(source, ".")
	setRegExes := make(map[string]*regexp.Regexp)
	for _, r := range set {
		character := string(r)
		setRegExes[character] = regexp.MustCompile("[" + character + "]{2,}")
	}
	beginEndRegEx := regexp.MustCompile("(^[" + set + "]|[" + set + "]$)")

	// Process
	for i := range parts {
		for c, re := range setRegExes {
			parts[i] = re.ReplaceAllString(parts[i], c)
		}
		for beginEndRegEx.MatchString(parts[i]) {
			parts[i] = beginEndRegEx.ReplaceAllString(parts[i], "")
		}
	}
	result = strings.Join(parts, ".")
	return
}

func truncate(source string, length int) (result string) {
	result = source
	if utf8.RuneCountInString(source) > length {
		var runes []rune
		i := 0
		for _, r := range source {
			if i > length {
				break
			}
			runes = append(runes, r)
			i++
		}
		result = string(runes)
	}
	return
}
