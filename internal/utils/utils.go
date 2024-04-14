package utils

type declensionSet struct {
	nominative string
	genitive   string
}

var declensionMonths = map[int32]declensionSet{
	1:  {nominative: "январь", genitive: "января"},
	2:  {nominative: "февраль", genitive: "февраля"},
	3:  {nominative: "март", genitive: "марта"},
	4:  {nominative: "апрель", genitive: "апреля"},
	5:  {nominative: "май", genitive: "мая"},
	6:  {nominative: "июнь", genitive: "июня"},
	7:  {nominative: "июль", genitive: "июля"},
	8:  {nominative: "август", genitive: "августа"},
	9:  {nominative: "сентябрь", genitive: "сентября"},
	10: {nominative: "октябрь", genitive: "октября"},
	11: {nominative: "ноябрь", genitive: "ноября"},
	12: {nominative: "декабрь", genitive: "декабря"},
}

func DeclensionGenitiveMonth(monthNumber int32) string {
	if monthNumber < 1 || monthNumber > 12 {
		monthNumber = 1
	}

	return declensionMonths[monthNumber].genitive
}
