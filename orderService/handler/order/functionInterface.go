package order

import pb "orderService/proto/order"

type NameLogic interface {
	IsEnglish(string) bool
	IsCapitalized(string) bool
}

type PriceLogic interface {
	PriceIsNumber(string) bool
	PriceOverTwoThousands(string) bool
	PriceUSDToTWD(string) string
	GetPriceInt(string) int
}

type CurrencyLogic interface {
	IsCurrencyFormatValid(string) bool
	TransformUSDToTWD(currency, price string) (string, string)
}

type ResponseLogic interface {
	SetResponseFailure(res *pb.CheckAndTransformDataResponse, message string)
	SetResponseSuccess(req *pb.CheckAndTransformDataRequest, res *pb.CheckAndTransformDataResponse, message string)
}

func NewOrder(name NameLogic, price PriceLogic, currency CurrencyLogic, response ResponseLogic) *Order {
	return &Order{
		name:     name,
		price:    price,
		currency: currency,
		response: response,
	}
}

func NewCurrency(price PriceLogic) *CurrencyHandler {
	return &CurrencyHandler{
		price: price,
	}
}
