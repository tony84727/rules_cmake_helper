package main

import (
	"testing"
)

func TestReplacerCMake(t *testing.T) {
    cases := []struct {
        variables map[string]string
        template string
        expected string
        notCMake bool
    }{
        {
            variables: map[string]string{
                "HAVE_GOOD_TIME":     "1",
                "SAME_PREFIX":        `"one"`,
                "SAME_PREFIX_PREFIX": `"two"`,
                "VERSION":            "1.1.1",
            },
            template: `#cmakedefine HAVE_GOOD_TIME
#cmakedefine HAVE_BAD_TIME
#cmakedefine SAME_PREFIX
#cmakedefine SAME_PREFIX_PREFIX

std::string get_version() {
    return "@VERSION@";
}`,
            expected: `#define HAVE_GOOD_TIME 1
/* #undef HAVE_BAD_TIME */
#define SAME_PREFIX "one"
#define SAME_PREFIX_PREFIX "two"

std::string get_version() {
    return "1.1.1";
}`,
        },
        {
            variables: map[string]string{
                "HAVE_PREDEFINED": "1",
            },
            template: `#cmakedefine HAVE_PREDEFINED 10
#cmakedefine SHOULD_NOT_DEFINED 1`,
            expected: `#define HAVE_PREDEFINED 10
/* #undef SHOULD_NOT_DEFINED 1 */`,
        },
        {
            variables: map[string]string{
                "INLINE_KEYWORD": "__inline",
            },
            template: `#ifndef __cplusplus
#define inline @INLINE_KEYWORD@
#endif`,
            expected: `#ifndef __cplusplus
#define inline __inline
#endif`,
            notCMake: true,
        },
    }

    for _, c := range cases {
        replacer := Replacer{
            variables: c.variables,
            notCMake: c.notCMake,
        }
        out := replacer.Replace(c.template)

        if c.expected != out {
            t.Errorf("expected:\n%s\ngot:\n%s\ntemplate:\n%s\nvariables:%v", c.expected, out,c.template, c.variables)
        }
    }
}
