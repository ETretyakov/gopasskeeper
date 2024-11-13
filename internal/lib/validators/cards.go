package validators

import (
	"errors"
	"strconv"
	"strings"
)

// CreditCard is a structure to represent a credit card.
type CreditCard struct {
	Number string
	Month  int32
	Year   int32
	CVC    string
	PIN    string
}

// NewCreditCard is a builder function for CreditCard.
func NewCreditCard(
	number string,
	month, year int32,
	cvc, pin string,
) *CreditCard {
	number = strings.ReplaceAll(number, " ", "")
	cvc = strings.ReplaceAll(cvc, " ", "")
	pin = strings.ReplaceAll(pin, " ", "")

	return &CreditCard{
		Number: number,
		Month:  month,
		Year:   year,
		CVC:    cvc,
		PIN:    pin,
	}
}

// Validate will check the credit card's number against the Luhn algorithm
func (c *CreditCard) Validate() error {
	if c.Month < 1 || 12 < c.Month {
		return errors.New("invalid month value")
	}

	if c.Year < 1970 {
		return errors.New("invalid year value")
	}

	if len(c.CVC) < 3 || len(c.CVC) > 4 {
		return errors.New("invalid cvc value")
	}

	if len(c.PIN) < 4 || len(c.PIN) > 6 {
		return errors.New("invalid pin value")
	}

	var sum int
	var alternate bool

	numberLen := len(c.Number)

	if numberLen < 13 || numberLen > 19 {
		return errors.New("invalid number value")
	}

	for i := numberLen - 1; i > -1; i-- {
		mod, _ := strconv.Atoi(string(c.Number[i]))
		if alternate {
			mod *= 2
			if mod > 9 {
				mod = (mod % 10) + 1
			}
		}

		alternate = !alternate

		sum += mod
	}

	if !(sum%10 == 0) {
		return errors.New("invalid number value")
	}

	return nil
}

// Mask is retrieving card number mask.
func (c *CreditCard) Mask() string {
	mask := ""

	for i := 0; i < len(c.Number)-4; i++ {
		mask += "*"

		if i%4 == 3 {
			mask += " "
		}
	}

	return mask + c.Number[len(c.Number)-4:]
}
