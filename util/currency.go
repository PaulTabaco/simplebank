package util

// TODO: - Look for golang way check if string in list

const (
	USD = "USD"
	EUR = "EUR"
	CAD = "CAD"
	UAH = "UAH"
)

func IsSupportedCurrency(currency string) bool {
	switch currency {
	case USD, EUR, CAD, UAH:
		return true
	}
	return false
}
