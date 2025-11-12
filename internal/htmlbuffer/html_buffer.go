package htmlbuffer

const (
	openingBracket = '<'
	closingBracket = '>'
)

type HTMLBuffer struct {
	data       []byte
	cursor     int
	currentTag tag
}

func New(data []byte) HTMLBuffer {
	return HTMLBuffer{
		data:   data,
		cursor: 0,
	}
}

func (hb *HTMLBuffer) SeekToIdentifiedTag(tag, identifier string) bool {
	data := hb.data[hb.cursor:]

	for index, value := range data {
		if isOpeningBracket(value) && value != '/' {
			var (
				tagStart = index + 1
				tagEnd   = tagStart + len(tag)
			)

			if len(data) < tagEnd {
				return false
			}

			tagName := string(data[tagStart:tagEnd])
			if tagName != tag {
				continue
			}

			tagEndCursor := hb.cursor + tagEnd
			tagClosingBracketIndex, isIdentified := hb.isIdentifiedWith(tagEndCursor, identifier)
			if !isIdentified {
				continue
			}

			hb.currentTag = newTag(tagName)
			hb.cursor = hb.cursor + tagEnd + tagClosingBracketIndex
			return true
		}
	}

	return false
}

func (hb *HTMLBuffer) SeekToTag(tag string) bool {
	data := hb.data[hb.cursor:]

	for index, value := range data {
		if isOpeningBracket(value) && value != '/' {
			var (
				tagStart = index + 1
				tagEnd   = tagStart + len(tag)
			)

			if len(data) < tagEnd {
				return false
			}

			tagName := string(data[tagStart:tagEnd])
			if tagName != tag {
				continue
			}

			tagEndCursor := hb.cursor + tagEnd
			tagClosingBracketIndex, isFound := hb.findClosingBracket(tagEndCursor)
			if !isFound {
				continue
			}

			hb.currentTag = newTag(tagName)
			hb.cursor = hb.cursor + tagEnd + tagClosingBracketIndex
			return true
		}
	}

	return false
}

func (hb *HTMLBuffer) ReadTagValue() (string, bool) {
	if hb.currentTag.isZero() || hb.currentTag.isVoid {
		return "", false
	}

	data := hb.data[hb.cursor+1:]

	for index, value := range data {
		if isOpeningBracket(value) {
			if len(data) < index+1+len(hb.currentTag.name) {
				return "", false
			}

			if data[index+1] != '/' {
				continue
			}

			tagName := data[index+2 : index+2+len(hb.currentTag.name)]
			if string(tagName) != hb.currentTag.name {
				continue
			}

			tagValue := data[:index]
			return string(tagValue), true
		}
	}

	return "", false
}

func (hb *HTMLBuffer) findClosingBracket(tagEndCursor int) (int, bool) {
	data := hb.data[tagEndCursor:]

	for index, value := range data {
		if value == closingBracket {
			return index, true
		}
	}

	return 0, false
}

func (hb *HTMLBuffer) isIdentifiedWith(tagEndCursor int, identifier string) (int, bool) {
	data := hb.data[tagEndCursor:]

	var isIdentified bool

	for index, value := range data {
		switch value {
		case 'i':
			if containsIdentifierByKey(data, index, identifier, "id") {
				isIdentified = true
			}
		case 'c':
			if containsIdentifierByKey(data, index, identifier, "class") {
				isIdentified = true
			}
		case closingBracket:
			return index, isIdentified
		}
	}

	return 0, false
}

func containsIdentifierByKey(data []byte, index int, identifier, key string) bool {
	idStart := index + len(key) + 1
	if string(data[index:idStart]) != (key + "=") {
		return false
	}

	var (
		identifierStart = idStart + 1
		identifierValue = string(data[identifierStart : identifierStart+len(identifier)])
	)

	return identifierValue == identifier
}

func isOpeningBracket(value byte) bool {
	return rune(value) == openingBracket
}
