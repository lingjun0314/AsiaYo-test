package order

import "strconv"

type PriceHandler struct{}

func (p *PriceHandler) PriceIsNumber(price string) bool {
	_, err := strconv.Atoi(price)
	return err == nil
}

func (p *PriceHandler) PriceOverTwoThousands(price string) bool {
	priceInt, _ := strconv.Atoi(price)

	return priceInt > 2000
}

func (p *PriceHandler) GetPriceInt(price string) int {
	priceInt, _ := strconv.Atoi(price)
	return priceInt
}

func (p *PriceHandler) PriceUSDToTWD(price string) string {
	priceInt, _ := strconv.Atoi(price)
	priceInt *= 31
	return strconv.Itoa(priceInt)
}
