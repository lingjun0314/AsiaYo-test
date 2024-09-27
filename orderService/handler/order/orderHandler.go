package order

import (
	"context"
	pb "orderService/proto/order"
)

type Order struct {
	name     NameLogic
	price    PriceLogic
	currency CurrencyLogic
	response ResponseLogic
}

func (con *Order) CheckAndTransformData(ctx context.Context, req *pb.CheckAndTransformDataRequest, res *pb.CheckAndTransformDataResponse) error {
	order := NewOrder(&NameHandler{}, &PriceHandler{}, &CurrencyHandler{}, &ResponseHandler{})
	//	Validate name
	if !order.name.IsEnglish(req.Order.Name) {
		order.response.SetResponseFailure(res, "Name contains non-English characters")
		return nil
	}
	if !order.name.IsCapitalized(req.Order.Name) {
		order.response.SetResponseFailure(res, "Name is not capitalized")
		return nil
	}

	//	Validate price
	if !order.price.PriceIsNumber(req.Order.Price) {
		order.response.SetResponseFailure(res, "Price is not a number")
		return nil
	}

	if order.price.PriceOverTwoThousands(req.Order.Price) {
		order.response.SetResponseFailure(res, "Price is over 2000")
		return nil
	}

	//	Validate currency
	currency := NewCurrency(&PriceHandler{})

	if !currency.IsCurrencyFormatValid(req.Order.Currency) {
		order.response.SetResponseFailure(res, "Currency format is wrong")
		return nil
	}

	//	Change USD to TWD
	req.Order.Currency, req.Order.Price = currency.TransformUSDToTWD(req.Order.Currency, req.Order.Price)

	order.response.SetResponseSuccess(req, res, "Transform success!")

	return nil
}
