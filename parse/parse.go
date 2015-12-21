package parse

import (
	"bufio"
	"fmt"
	"io"
	"strings"

	"golang.org/x/text/encoding/charmap"
)

type ComprehensiveRules map[Title]*Rule

type Rule struct {
	Title
	Body     string   `json:",omitempty"`
	Notes    []string `json:",omitempty"`
	Examples []string `json:",omitempty"`
	Parent   *Rule    `json:"-"`
	Children []*Rule  `json:",omitempty"`
}

func NewRule(title string) *Rule {
	return &Rule{
		Title: Title(title),
	}
}

func (r *Rule) Text() string {
	if r.Title == "HEAD" {
		return ""
	}

	lines := make([]string, 1, 1+len(r.Notes)+len(r.Body))
	for _, note := range r.Notes {
		lines = append(lines, notePrefix+note)
	}
	for _, example := range r.Examples {
		lines = append(lines, examplePrefix+example)
	}
	return fmt.Sprintf("%v %v", r.Title, r.Body)
}

func (r *Rule) CompleteText() string {
	lines := make([]string, 1, len(r.Children)+1)
	lines[0] = r.Text()
	for _, child := range r.Children {
		lines = append(lines, child.CompleteText())
	}
	return strings.Join(lines, "\n\n")
}

func (r *Rule) addChild(child *Rule) {
	if child == nil || r == nil || child.Parent != nil {
		panic(fmt.Sprintf("cannot add child %+v to %+v", child, r))
	}
	r.Children = append(r.Children, child)
	child.Parent = r
}

type Title string

func (t Title) Depth() int {
	str := string(t)
	// Synthetic Head
	if str == "HEAD" {
		return -1
	}
	// 100.1a
	if !strings.HasSuffix(str, ".") {
		return 3
	}
	dotSplit := strings.Split(str, ".")
	// 100.1.
	if len(dotSplit) == 3 {
		return 2
	}
	// 100.
	if len(dotSplit[0]) == 3 {
		return 1
	}
	// 1.
	return 0
}

func (t Title) Name() string {
	str := string(t)
	switch t.Depth() {
	// Synthetic Head
	case -1:
		return "HEAD"
	// 1.
	case 0:
		fallthrough
	// 100.
	case 1:
		return strings.Split(str, ".")[0]
	// 100.1.
	case 2:
		return strings.Split(str, ".")[1]
	// 100.1a
	case 3:
		return str[len(str)-1:]
	default:
		panic(fmt.Sprintf("unknown depth %v", t.Depth()))
	}
}

func (t Title) ParentTitle() Title {
	str := string(t)
	switch t.Depth() {
	// Synthetic Head
	case -1:
		return ""
	// 1.
	case 0:
		return "HEAD"
	// 100.
	case 1:
		return Title(str[:1] + ".")
	// 100.1.
	case 2:
		return Title(strings.Split(str, ".")[0] + ".")
	// 100.1a
	case 3:
		return Title(str[:len(str)-1] + ".")
	default:
		panic(fmt.Sprintf("unknown depth %v", t.Depth()))
	}
}

const (
	crModeHeader = "HEADER"
	crModeRules  = "RULES"
)

const (
	notePrefix    = "     "
	examplePrefix = "Example:"
)

func ParseCR(r io.Reader) ComprehensiveRules {
	cr := make(ComprehensiveRules)

	cr["HEAD"] = NewRule("HEAD")
	scanner := bufio.NewScanner(charmap.Windows1252.NewDecoder().Reader(r))
	mode := crModeHeader
	var lastRule *Rule
	for scanner.Scan() {
		line := scanner.Text()
		if mode == crModeRules && line == "Glossary" {
			break
		}
		if mode == crModeHeader && line == "Credits" {
			mode = crModeRules
			continue
		}

		if mode == crModeRules && len(line) > 0 {
			titleSplit := strings.SplitN(line, " ", 2)
			if strings.HasPrefix(line, notePrefix) {
				lastRule.Notes = append(lastRule.Notes, strings.TrimSpace(line))
			} else if titleSplit[0] == examplePrefix {
				lastRule.Examples = append(lastRule.Examples, titleSplit[1])
			} else {
				rule := NewRule(titleSplit[0])
				rule.Body = titleSplit[1]

				lastRule = rule
				cr[rule.Title] = rule
				cr[rule.ParentTitle()].addChild(rule)
			}
		}
	}

	return cr
}
