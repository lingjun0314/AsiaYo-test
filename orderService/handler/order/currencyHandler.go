package order

type CurrencyHandler struct {
	price PriceLogic
}

func (c *CurrencyHandler) IsCurrencyFormatValid(currency string) bool {
	return currency == "TWD" || currency == "USD"
}

func (c *CurrencyHandler) TransformUSDToTWD(currency, price string) (string, string) {
	newPrice := price
	newCurrency := currency
	if currency == "USD" {
		newPrice = c.price.PriceUSDToTWD(price)
		newCurrency = "TWD"
	}
	return newCurrency, newPrice
}
