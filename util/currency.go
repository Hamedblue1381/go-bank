package util

const (
	EUR  = "EUR"
	USD  = "USD"
	RIAL = "RIAL"
	CAD  = "CAD"
)

func IsSuppotortedCurrnecy(currency string) bool {
	switch currency {
	case USD, EUR, CAD, RIAL:
		return true
	}
	return false
}
