package src

import (
	"errors"
	"os"
	"regexp"
	"strings"

	"github.com/fatih/color"
	"github.com/urfave/cli/v2"
)

type Input struct {
	PathFromStdin     string
	RegexArg          string
	RegexFlag         cli.StringSlice
	LineStrings       cli.StringSlice
	RegexFiles        cli.StringSlice
	SearchFromLeft    bool
	ShowBaseDirectory bool
	MatchAll          bool
	ShowColour        bool
}

// Finds the project root
func FindMatchesInPath(c Input) ([]string, error) {

	rgx, err := compileRegex(c)
	if err != nil {
		return nil, err
	}

	c.PathFromStdin = strings.ReplaceAll(c.PathFromStdin, "\n", "")
	pathDirs := strings.Split(c.PathFromStdin, "/")

	var toReturn []string

	pathPos := 0
	for i := len(pathDirs); i > 0; i-- {

		// Invert pathPos depending on flag
		if c.SearchFromLeft {
			pathPos = len(pathDirs) + 1 - i // +1 as i will only reach 1
		} else {
			pathPos = i
		}

		path := strings.Join(pathDirs[:pathPos], "/")

		if path == "" {
			path = "/"
		}

		names, err := os.ReadDir(path)
		if err != nil {
			if strings.Contains(err.Error(), "not a directory") {
				continue
			}
			return nil, err
		}

		for _, name := range names {
			mName := name.Name()

			isMatch := rgx.MatchString(mName)
			if isMatch {

				if c.ShowBaseDirectory {
					toReturn = append(toReturn, path)
				} else {
					if c.ShowColour {
						mName = applyColour(rgx, mName)
					}
					toReturn = append(toReturn, path+"/"+mName)
				}

				if len(toReturn) == 1 && !c.MatchAll {
					return toReturn, nil
				}
			}
		}
	}

	if len(toReturn) == 0 {
		return []string{}, errors.New("no match")
	}

	return toReturn, nil
}

// Compiles all the patterns into one regex pattern
func compileRegex(c Input) (*regexp.Regexp, error) {
	combinedRegexp := c.RegexArg

	for _, p := range c.RegexFiles.Value() {
		b, err := os.ReadFile(p)
		if err != nil {
			return nil, err
		}

		regexps := string(b)

		if len(regexps) == 0 {
			continue
		}
		combinedRegexp += strings.ReplaceAll(regexps, "\n", "|")
	}

	combinedRegexp += strings.Join(c.RegexFlag.Value(), "|")

	lnstr := strings.Join(c.LineStrings.Value(), "$|^")
	if len(lnstr) != 0 {
		lnstr = "^" + lnstr + "$"

		if combinedRegexp != "" {
			combinedRegexp += "|" + lnstr
		} else {
			combinedRegexp = lnstr
		}
	}

	r, err := regexp.Compile(combinedRegexp)
	return r, err
}

// Applies a colour on the matching term
func applyColour(rgx *regexp.Regexp, fileName string) string {
	loc := rgx.FindIndex([]byte(fileName))

	matched := color.GreenString(fileName[loc[0]:loc[1]])

	return fileName[:loc[0]] + matched + fileName[loc[1]:]
}
