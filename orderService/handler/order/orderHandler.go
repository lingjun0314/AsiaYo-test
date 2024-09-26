package order

import (
	"context"
	pb "orderService/proto/order"
	"regexp"
	"strconv"
	"strings"
)

type Order struct{}

func (con *Order) CheckAndTransformData(ctx context.Context, req *pb.CheckAndTransformDataRequest, res *pb.CheckAndTransformDataResponse) error {
	//	Validate name
	if !isEnglish(req.Order.Name) {
		SetResponseFailure(res,"Name contains non-English characters")
		return nil
	}
	if !isCapitalized(req.Order.Name) {
		SetResponseFailure(res, "Name is not capitalized")
		return nil
	}

	//	Validate price
	price, err := strconv.Atoi(req.Order.Price)
	if err != nil {
		SetResponseFailure(res, "Price not a number")
		return nil
	}
	if price > 2000 {
		SetResponseFailure(res, "Price is over 2000")
		return nil
	}

	//	Validate currency
	if req.Order.Currency != "TWD" && req.Order.Currency != "USD" {
		SetResponseFailure(res, "Currency format is wrong")
		return nil
	}

	//	Change USD to TWD
	if req.Order.Currency == "USD" {
		price *= 31
		priceString := strconv.Itoa(price)
		req.Order.Price = priceString
		req.Order.Currency = "TWD"
	}

	res.Order = req.Order
	res.Message = "Transform success"
	res.StatusCode = 200
	res.Success = true

	return nil
}

func SetResponseFailure(res *pb.CheckAndTransformDataResponse, message string) {
	res.Message = message
	res.StatusCode = 400
	res.Order = nil
	res.Success = false
}

func isEnglish(name string) bool {
	return regexp.MustCompile(`^[a-zA-Z\s]+$`).MatchString(name)
}

func isCapitalized(name string) bool {
	words := strings.Fields(name)
	for _, word := range words {
		if word[0] < 'A' || word[0] > 'Z' {
			return false
		}
	}
	return true
}
