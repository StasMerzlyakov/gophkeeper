package creditcard_test

import (
	"testing"

	"github.com/StasMerzlyakov/gophkeeper/internal/domain/creditcard"
	"github.com/stretchr/testify/assert"
)

func TestCards(t *testing.T) {
	assert := assert.New(t)
	card := creditcard.Card{
		Number: "2200400128400690", ExpiryMonth: 11, ExpiryYear: 2019, CVV: "123",
	}
	val := card.Validate()
	assert.Contains(val.Errors, "creditcard is expired")

	card = creditcard.Card{
		Type: "Something", Number: "5019717010103742", ExpiryMonth: 11, ExpiryYear: 2019, CVV: "1234",
	}
	val = card.Validate()
	assert.Contains(val.Errors, "given card type doesn't match determined card type")

	card = creditcard.Card{
		Type: "Something", Number: "5019717010103742", ExpiryMonth: 111, ExpiryYear: 2019, CVV: "1234",
	}
	val = card.Validate()
	assert.Contains(val.Errors, "month '111' is not a valid month")

	card = creditcard.Card{
		Type: "Something", Number: "5019717010103742", ExpiryMonth: 11, ExpiryYear: 1899, CVV: "1234",
	}
	val = card.Validate()
	assert.Contains(val.Errors, "year '1899' is not a valid year")

	card = creditcard.Card{
		Type: "Dankort", Number: "5019717010103742", ExpiryMonth: 11, ExpiryYear: 1899, CVV: "1234",
	}
	val = card.Validate()
	assert.Contains(val.Errors, "year '1899' is not a valid year")

	card = creditcard.Card{
		Number: "5019717010103742", ExpiryMonth: 11, ExpiryYear: 2020, CVV: "1234",
	}
	val = card.Validate()
	assert.Equal(val.Card.Type, "Dankort")

	card = creditcard.Card{
		Number: "0000000000", ExpiryMonth: 11, ExpiryYear: 2020, CVV: "1234",
	}
	val = card.Validate()
	assert.Contains(val.Errors, "unknown creditcard type")

	card = creditcard.Card{
		Number: "378282246310005", ExpiryMonth: 11, ExpiryYear: 2020, CVV: "1234",
	}
	val = card.Validate()
	assert.Equal(val.Card.Type, "American Express")

	card = creditcard.Card{
		Number: "655021246310005", ExpiryMonth: 11, ExpiryYear: 2020, CVV: "1234",
	}
	val = card.Validate()
	assert.Equal(val.Card.Type, "Elo")

	card = creditcard.Card{
		Number: "604201246310005", ExpiryMonth: 11, ExpiryYear: 2020, CVV: "1234",
	}
	val = card.Validate()
	assert.Equal(val.Card.Type, "Cabal")

	card = creditcard.Card{
		Number: "384140246310005", ExpiryMonth: 11, ExpiryYear: 2020, CVV: "1234",
	}
	val = card.Validate()
	assert.Equal(val.Card.Type, "Hipercard")

	card = creditcard.Card{
		Number: "560221246310005", ExpiryMonth: 11, ExpiryYear: 2020, CVV: "1234",
	}
	val = card.Validate()
	assert.Equal(val.Card.Type, "Bankcard")

	card = creditcard.Card{
		Number: "620221246310005", ExpiryMonth: 11, ExpiryYear: 2020, CVV: "1234",
	}
	val = card.Validate()
	assert.Equal(val.Card.Type, "China UnionPay")

	card = creditcard.Card{
		Number: "300221246310005", ExpiryMonth: 11, ExpiryYear: 2020, CVV: "1234",
	}
	val = card.Validate()
	assert.Equal(val.Card.Type, "Diners Club Carte Blanche")

	card = creditcard.Card{
		Number: "201421246310005", ExpiryMonth: 11, ExpiryYear: 2020, CVV: "1234",
	}
	val = card.Validate()
	assert.Equal(val.Card.Type, "Diners Club Enroute")

	card = creditcard.Card{
		Number: "39022124631000", ExpiryMonth: 11, ExpiryYear: 2020, CVV: "1234",
	}
	val = card.Validate()
	assert.Equal(val.Card.Type, "Diners Club International")

	card = creditcard.Card{
		Number: "601121246310005", ExpiryMonth: 11, ExpiryYear: 2020, CVV: "1234",
	}
	val = card.Validate()
	assert.Equal(val.Card.Type, "Discover")

	card = creditcard.Card{
		Number: "63612124631000500", ExpiryMonth: 11, ExpiryYear: 2020, CVV: "1234",
	}
	val = card.Validate()
	assert.Equal(val.Card.Type, "InterPayment")

	card = creditcard.Card{
		Number: "6371212463100050", ExpiryMonth: 11, ExpiryYear: 2020, CVV: "1234",
	}
	val = card.Validate()
	assert.Equal(val.Card.Type, "InstaPayment")

	card = creditcard.Card{
		Number: "501821246310005", ExpiryMonth: 11, ExpiryYear: 2020, CVV: "1234",
	}
	val = card.Validate()
	assert.Equal(val.Card.Type, "Maestro")

	card = creditcard.Card{
		Number: "511821246310005", ExpiryMonth: 11, ExpiryYear: 2020, CVV: "1234",
	}
	val = card.Validate()
	assert.Equal(val.Card.Type, "Mastercard")

	card = creditcard.Card{
		Number: "351821246310005", ExpiryMonth: 11, ExpiryYear: 2020, CVV: "1234",
	}
	val = card.Validate()
	assert.Equal(val.Card.Type, "JCB")

	card = creditcard.Card{
		Number: "508821246310005", ExpiryMonth: 11, ExpiryYear: 2020, CVV: "1234",
	}
	val = card.Validate()
	assert.Equal(val.Card.Type, "Aura")

	card = creditcard.Card{
		Number: "402621246310005", ExpiryMonth: 11, ExpiryYear: 2020, CVV: "1234",
	}
	val = card.Validate()
	assert.Equal(val.Card.Type, "Visa Electron")

	card = creditcard.Card{
		Number: "409921246310005", ExpiryMonth: 11, ExpiryYear: 2020, CVV: "1234",
	}
	val = card.Validate()
	assert.Equal(val.Card.Type, "Visa")

	card = creditcard.Card{
		Number: "0000000000", ExpiryMonth: 11, ExpiryYear: 2020, CVV: "1234",
	}
	luhn := card.ValidateLuhn()
	assert.Equal(luhn, false)
}
