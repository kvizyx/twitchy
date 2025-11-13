package htmlbuffer

const (
	openingBracket = '<'
	closingBracket = '>'
)

type HTMLBuffer struct {
	data       []byte
	cursor     int
	currentTag tag
	closingTag tag
}

func New(data []byte) HTMLBuffer {
	return HTMLBuffer{
		data:   data,
		cursor: 0,
	}
}

// MarkClosingTag marks current tag as closing tag so buffer will be not allowed to seek after cursor will pass current
// tag closer.
func (hb *HTMLBuffer) MarkClosingTag() {
	hb.closingTag = newTag(hb.currentTag.name)
}

// SeekToIdentifiedTag seeks cursor to the first occurrence of tag with given identifier (id or class) and returns whether
// it was found or not.
func (hb *HTMLBuffer) SeekToIdentifiedTag(tag, identifier string) bool {
	return hb.seek(tag, identifier)
}

// SeekToTag seeks cursor to the first occurrence of given tag and returns whether it was found or not.
func (hb *HTMLBuffer) SeekToTag(tag string) bool {
	return hb.seek(tag, "")
}

// ReadTagValue reads and returns value of tag that cursor is currently pointing on without seeking cursor.
// Empty string and false value will be returned if cursor is not pointing to any tag or its void tag.
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

func (hb *HTMLBuffer) seek(tag string, identifier string) bool {
	data := hb.data[hb.cursor:]

	for index, value := range data {
		if !isOpeningBracket(value) {
			continue
		}

		if rune(data[index+1]) == '/' {
			if !hb.isClosingTag(data, index) {
				continue
			}

			return false
		}

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

		var (
			tagClosingBracketIndex int
			pass                   bool
		)

		if len(identifier) == 0 {
			tagClosingBracketIndex, pass = hb.findClosingBracket(tagEndCursor)
		} else {
			tagClosingBracketIndex, pass = hb.isIdentifiedWith(tagEndCursor, identifier)
		}

		if !pass {
			continue
		}

		hb.currentTag = newTag(tagName)
		hb.cursor = hb.cursor + tagEnd + tagClosingBracketIndex
		return true
	}

	return false
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

func (hb *HTMLBuffer) isClosingTag(data []byte, openingBracketIndex int) bool {
	if hb.closingTag.isZero() {
		return false
	}

	var (
		tagStart      = openingBracketIndex + 2
		closingTagEnd = tagStart + len(hb.closingTag.name) + 1
	)

	for index, value := range string(data[tagStart:closingTagEnd]) {
		if rune(value) != closingBracket {
			continue
		}

		if string(data[tagStart:tagStart+index]) == hb.closingTag.name {
			return true
		}
	}

	return false
}

func (hb *HTMLBuffer) isIdentifiedWith(tagEndIndex int, identifier string) (int, bool) {
	data := hb.data[tagEndIndex:]

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
