package mack

import (
	"strconv"
	"strings"
)

const (
	OPENBRACE      = "{"
	CLOSEBRACE     = "}"
	COMMENT        = "/*"
	ENDCOMMENT     = "*/"
	LINECOMMENT    = "#"
	ENDLINECOMMENT = "\n"
	QUOTE          = "\""
	//ESCAPEDQUOTE          = "\\\""
	CHARQUOTE = "'"
	//SEMI        = ";"
	OPENBRACKET  = "["
	CLOSEBRACKET = "]"
	OPENPAREN    = "("
	CLOSEPAREN   = ")"
	COLON        = ":"
	COMMA        = ","
)

func getCloseString(openString string) string {
	closeStrings := map[string]string{
		"class":     CLOSEBRACKET,
		OPENBRACKET: CLOSEBRACKET,
		OPENBRACE:   CLOSEBRACE,
		COMMENT:     ENDCOMMENT,
		LINECOMMENT: ENDLINECOMMENT,
		QUOTE:       QUOTE,
		CHARQUOTE:   CHARQUOTE,
		OPENPAREN:   CLOSEPAREN,
	}
	return closeStrings[openString]
}

func inQuote(blockStack []string) bool {
	if len(blockStack) == 0 {
		return false
	}
	return blockStack[len(blockStack)-1] == QUOTE || blockStack[len(blockStack)-1] == CHARQUOTE
}

func parseToMap(contents string) map[string]string {
	var blockStack []string
	var blocks []string
	wordBuffer := ""
	var prefix string

	valueMap := make(map[string]string)

	for _, char := range contents {
		current := string(char)
		if inQuote(blockStack) {
			// the ONLY thing I'm looking for in these cases is the end of the comment
			wordBuffer += string(char)
			if len(blockStack) > 0 && strings.HasSuffix(wordBuffer, getCloseString(blockStack[len(blockStack)-1])) {
				blockStack = blockStack[:len(blockStack)-1]
				// keep the quote - as if it's part of the value
				// wordBuffer = wordBuffer[0 : len(wordBuffer)-1]
				// now we continue towards the end of field character
				continue
			} else {
				continue
			}
		}

		if current == OPENBRACE {
			if wordBuffer != "" {
				blocks = append(blocks, wordBuffer)
				wordBuffer = ""
			}
			blockStack = append(blockStack, OPENBRACE)
			continue
		}

		if current == QUOTE {
			blockStack = append(blockStack, QUOTE)
			// let's keep the quote to indicate strings from non-strings
			wordBuffer += string(char)
			continue
		}

		if current == COLON {
			prefix = strings.TrimSpace(strings.ReplaceAll(wordBuffer, "\n", " "))
			wordBuffer = ""
			continue
		}

		if current == COMMA {
			if strings.HasSuffix(wordBuffer, QUOTE) && !strings.HasPrefix(wordBuffer, QUOTE) {
				// trim off anything up to the first quote
				wordBuffer = wordBuffer[strings.Index(wordBuffer, QUOTE):]
			}
			valueMap[prefix] = strings.TrimSpace(wordBuffer)
			wordBuffer = ""
			continue
		}

		if current == CLOSEBRACE {
			wordBuffer = strings.TrimSpace(wordBuffer)
			if strings.HasSuffix(wordBuffer, QUOTE) && !strings.HasPrefix(wordBuffer, QUOTE) {
				// trim off anything up to the first quote
				wordBuffer = wordBuffer[strings.Index(wordBuffer, QUOTE):]
			}
			valueMap[prefix] = strings.TrimSpace(wordBuffer)
			wordBuffer = ""
			continue
		}
		wordBuffer += string(char)
	}

	return valueMap
}

func quoteWrap(str string) string {
	return strings.Join([]string{QUOTE, str, QUOTE}, "")
}

func finalValue(value string) string {
	// it's already a string, let it be
	if strings.HasPrefix(value, QUOTE) && strings.HasSuffix(value, QUOTE) {
		return value
	}
	// json-safe enums
	switch value {
	case "true":
		return value
	case "false":
		return value
	default:
		break
	}
	// numbers
	if _, err := strconv.ParseFloat(value, 64); err == nil {
		return value
	}
	if _, err := strconv.ParseInt(value, 10, 64); err == nil {
		return value
	}
	// anything else probably needs quotes
	return quoteWrap(value)
}

func stringMapToJson(source map[string]string) string {
	var outlist []string
	for key, value := range source {
		outlist = append(outlist, strings.Join([]string{quoteWrap(key), finalValue(value)}, ":"))
		source[key] = strings.ReplaceAll(value, "\n", " ")
	}
	joined := strings.Join(outlist, ",")
	return strings.Join([]string{OPENBRACE, joined, CLOSEBRACE}, "")
}

// ResponseToJson converts a string response to a JSON string, making object conversion easier
func ResponseToJson(response string) string {
	return stringMapToJson(parseToMap(response))
}
