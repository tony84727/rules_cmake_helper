package main

import (
    "flag"
    "fmt"
    "io/ioutil"
    "log"
    "regexp"
    "strings"
)

var (
	output   string
	notCMake bool
)

func ReplaceSubgroup(src string, pattern *regexp.Regexp, replacer func(submatches []string) []string) (replaced string) {
	matches := pattern.FindAllStringSubmatchIndex(src, -1)

	lastEnd := 0
	for _, match := range matches {
		// pop first submatch pair, which is match of whole expression
		match = match[2:]
		submatches := make([]string, len(match)/2)
		for p := 0; p < len(match)/2; p++ {
			start := match[p*2]
			end := match[p*2+1]
			if start < 0 {
				submatches[p] = ""
				continue
			}
			submatches[p] = src[start:end]
		}
		replaces := replacer(submatches)
		if len(replaces) > len(match)/2 {
			replaces = replaces[:len(match)/2]
		}
		for i, r := range replaces {
		    if len(submatches[i]) <= 0 {
		        if len(r) > 0 {
                    replaced += r
                }
		        continue
            }
			start := match[i*2]
			end := match[i*2+1]
			replaced += src[lastEnd:start]
			replaced += r
			lastEnd = end
		}
	}
	replaced += src[lastEnd:]

	return replaced
}

type Replacer struct {
	variables map[string]string
	notCMake  bool
}

func (r Replacer) compositePattern() *regexp.Regexp {
	keyword := "cmakedefine"
	if r.notCMake {
		keyword = "undef"
	}
	pattern := fmt.Sprintf(`(?m)^\s*(#%s)\s+(\S+)(\s+\S+)?$|(@\S+@)`, keyword)
	return regexp.MustCompile(pattern)
}

func (r Replacer) Replace(origin string) (replaced string) {
	pattern := r.compositePattern()

	return ReplaceSubgroup(origin, pattern, func(submatches []string) []string {
		defineVariable := submatches[1]
		defineValue := strings.TrimSpace(submatches[2])
		placeholderVariable := submatches[3]

		var defineReplacement []string
		if len(defineVariable) > 0 {
			toReplace, exists := r.variables[defineVariable]
			if exists {
			    r := toReplace
			    if len(defineValue) > 0 {
			        r = defineValue
                }
				defineReplacement = []string{"#define", defineVariable, " " + r}
			} else {
			    r := " */"
			    if len(defineValue) > 0 {
			        r = " " + defineValue + r
                }
				defineReplacement = []string{"/* #undef", defineVariable,  r}
			}
		} else {
			defineReplacement = make([]string, 3)
		}
		if len(placeholderVariable) > 0 {
			// remove @ prefix and suffix
			actual := placeholderVariable[1 : len(placeholderVariable)-1]
			toReplace, exists := r.variables[actual]

			if exists {
				return append(defineReplacement, toReplace)
			}
			// untouched
			return append(defineReplacement, placeholderVariable)
		}
		return append(defineReplacement, "")
	})
}

func parseVariables(source string) map[string]string {
	variables := make(map[string]string)
	lines := strings.Split(source, "\n")
	for _, l := range lines {
		segments := strings.Split(strings.TrimSpace(l), "=")
		if len(segments) < 2 {
			variables[segments[0]] = "1"
		} else {
			variables[segments[0]] = segments[1]
		}
	}
	return variables
}

func main() {
    flag.BoolVar(&notCMake, "nocmake", false, "not cmake template")
    flag.StringVar(&output, "output", "", "output file")
    flag.Parse()
	if flag.NArg() < 2 {
		log.Fatal("Usage: configure [-output output] [-nocmake] <variable map> <source>")
	}

	variableFilePath := flag.Arg(0)
	sourceFilePath := flag.Arg(1)
	variableSource, err := ioutil.ReadFile(variableFilePath)
	if err != nil {
		log.Fatal(err)
	}
	variables := parseVariables(string(variableSource))
	replacer := Replacer{variables: variables, notCMake: notCMake}
	source, err := ioutil.ReadFile(sourceFilePath)
	if err != nil {
		log.Fatal(err)
	}
	replaced := replacer.Replace(string(source))
	if len(output) <= 0 {
		fmt.Print(replaced)
	} else {
		if err := ioutil.WriteFile(output, []byte(replaced), 0755); err != nil {
			log.Fatal(err)
		}
	}
}
