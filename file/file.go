package file

import (
	"io"
	"os"
	"path/filepath"
	"strings"
)

const (
	// Mon Jan 2 15:04:05 -0700 MST 2006 represented as yyyyMMddhhmm
	touchTimeFormat = "200601021504"

	// OS allow 255 character for filenames = 252 + 3 (.md)
	maxNameChars = 252
)

// Save a new file in a given dir with the following content.
// Creates a directory if necessary.
func Save(dir, name string, content io.Reader) error {
	if len(name) == 0 {
		return nil
	}

	if err := os.MkdirAll(dir, os.ModePerm); err != nil {
		return err
	}

	output, err := os.Create(filepath.Join(dir, name))
	if err != nil {
		return err
	}

	if _, err = io.Copy(output, content); err != nil {
		_ = output.Close()
		return err
	}

	return output.Close()
}

// BaseName normalizes a given string to use it as a safe filename
func BaseName(s string) string {

	// Remove any trailing space
	s = strings.Trim(s, " ")

	// Replace inappropriate characters
	s = strings.Replace(s, "<", "_", -1)  // < (less than) -> (u+3c less-than sign)
	s = strings.Replace(s, ">", ">", -1)  // > (greater than) -> (u+3e greater-than sign)
	s = strings.Replace(s, ":", "∶", -1)  // : (colon) -> (U+2236 RATIO)
	s = strings.Replace(s, "\"", "“", -1) // " (double quote) -> (U+201C english leftdoublequotemark)
	s = strings.Replace(s, "/", "∕", -1)  // / (forward slash) -> (DIVISION SLASH U+2215)
	s = strings.Replace(s, "\\", "⧵", -1) // \ (backslash) -> (U+29F5 Reverse solidus operator)
	s = strings.Replace(s, "|", "∣", -1)  // | (pipe) -> (U+2223 divides)
	s = strings.Replace(s, "?", "？", -1)  // ? (question mark) -> (U+FF1F FULLWIDTH QUESTION MARK)
	s = strings.Replace(s, "*", "＊", -1)  // * (asterisk) -> ＊ (Full Width Asterisk U+FF0A)

	// Check file name length in bytes
	if len(s) < maxNameChars {
		return s
	}

	// Trim filename to the max allowed number of bytes
	var sb strings.Builder
	var i = 0
	for index, c := range s {
		if index >= maxNameBytes || i >= maxNameChars {
			return sb.String()
		}
		sb.WriteRune(c)
		i++
	}

	return s
}
