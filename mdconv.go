package mdconv

import (
	"fmt"
	"html"
	"strings"
)

// ConvertParams is a struct with available conversion options
type ConvertParams struct {
	OpenBold    string
	CloseBold   string
	OpenItalic  string
	CloseItalic string
	OpenCode    string
	CloseCode   string
	LinkFormat  string
}

const (
	italicSep1    = "_"
	italicSep2    = "*"
	codeSep       = "`"
	linkFirstSym  = "["
	linkMediumSym = "]("
	linkLastSym   = ")"

	linkParam  = ":link:"
	aliasParam = ":alias:"
)

// ToText converts mdtext to simple text
func ToText(mdtext string) string {
	if len(mdtext) == 0 {
		return ""
	}

	return reformat(mdtext, &ConvertParams{
		OpenBold:    `"`,
		CloseBold:   `"`,
		OpenItalic:  "",
		CloseItalic: "",
		OpenCode:    "",
		CloseCode:   "",
		LinkFormat:  fmt.Sprintf(`%s (%s)`, aliasParam, linkParam),
	})
}

// ToText converts mdtext to HTML
func ToHTML(mdtext string) string {
	if len(mdtext) == 0 {
		return ""
	}
	mdtext = html.EscapeString(mdtext)

	return reformat(mdtext, &ConvertParams{
		OpenBold:    "<strong>",
		CloseBold:   "</strong>",
		OpenItalic:  "<em>",
		CloseItalic: "</em>",
		OpenCode:    `<font face="Courier New">`,
		CloseCode:   "</font>",
		LinkFormat:  fmt.Sprintf(`<a href="%s" target="_blank">%s</a>`, linkParam, aliasParam),
	})
}

// ToText converts mdtext to Slack message
func ToSlack(mdtext string) string {
	if len(mdtext) == 0 {
		return ""
	}
	mdtext = html.EscapeString(mdtext)

	return reformat(mdtext, &ConvertParams{
		OpenBold:    "*",
		CloseBold:   "*",
		OpenItalic:  "_",
		CloseItalic: "_",
		OpenCode:    "`",
		CloseCode:   "`",
		LinkFormat:  fmt.Sprintf(`<%s|%s>`, linkParam, aliasParam),
	})
}

// ToText converts mdtext to Hangouts message
func ToHangouts(mdtext string) string {
	return ToSlack(mdtext)
}

// ToText converts mdtext by parameters
func Convert(text string, params *ConvertParams) (string, error) {
	if len(text) == 0 || params == nil {
		return "", nil
	}

	if isCorrect := checkLinkFormat(params.LinkFormat); !isCorrect {
		return "", fmt.Errorf("wrong link format")
	}

	return reformat(text, params), nil
}

func reformat(text string, params *ConvertParams) string {
	var charLen int
	var openSym, closeSym string

	for i, sym := range text {
		char := string(sym)
		if char == linkFirstSym {
			charLen = 1
			return text[:i] + reformatWithLink(text[i+charLen:], params)
		} else if char == italicSep1 || char == italicSep2 {
			charLen = 1
			openSym = params.OpenItalic
			closeSym = params.CloseItalic
			if len(text) > i+1 && string(text[i+1]) == char {
				charLen = 2
				openSym = params.OpenBold
				closeSym = params.CloseBold
				char += char
			}
		} else if char == codeSep {
			charLen = 1
			openSym = params.OpenCode
			closeSym = params.CloseCode
		} else {
			continue
		}

		closeIndex := strings.Index(text[i+charLen:], char)
		if closeIndex != -1 {
			substr := text[i+charLen : i+charLen+closeIndex]
			convSubstr := reformat(substr, params)

			var convEndSubstr string
			if len(text) > i+charLen+closeIndex+charLen {
				convEndSubstr = reformat(text[i+charLen+closeIndex+charLen:], params)
			}

			text = text[:i] + openSym + convSubstr + closeSym + convEndSubstr
			break
		}
	}

	return text
}

func reformatWithLink(text string, params *ConvertParams) string {
	closeIndex := strings.Index(text, linkLastSym)
	if closeIndex == -1 {
		return linkFirstSym + reformat(text, params)
	}
	mediumIndex := strings.Index(text[:closeIndex], linkMediumSym)
	if mediumIndex == -1 {
		return linkFirstSym + reformat(text, params)
	}

	link := strings.Replace(params.LinkFormat, aliasParam, text[:mediumIndex], 1)
	link = strings.Replace(link, linkParam, text[mediumIndex+2:closeIndex], 1)

	return link + reformat(text[closeIndex+len(linkLastSym):], params)
}

func checkLinkFormat(format string) bool {
	return strings.Index(format, ":link:") != -1 && strings.Index(format, ":alias:") != -1
}
