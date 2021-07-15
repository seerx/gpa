package utils

const (
	CapitalLetterMin   = 'A'
	CapitalLetterMax   = 'Z'
	LowercaseLetterMin = 'a'
	LowercaseLetterMax = 'z'
	NumberMin          = '0'
	NumberMax          = '9'
	UnderLine          = '_'
)

func IsLower(ch byte) bool {
	return ch >= LowercaseLetterMin && ch <= LowercaseLetterMax
}

func IsCapital(ch byte) bool {
	return ch >= CapitalLetterMin && ch <= CapitalLetterMax
}

func IsValidSQLFieldName(name string) bool {
	for n, ch := range name {
		if ch >= LowercaseLetterMin && ch <= LowercaseLetterMax {
			continue
		}
		if ch >= CapitalLetterMin && ch <= CapitalLetterMax {
			continue
		}

		if n > 0 {
			if ch == UnderLine {
				continue
			}
			if ch >= NumberMin && ch <= NumberMax {
				continue
			}
		}

		return false
	}
	return true
}
