package order_test

import (
	"context"
	"github.com/stretchr/testify/assert"
	o "orderService/handler/order"
	"orderService/proto/order"
	"strconv"
	"testing"
)

func TestCheckAndTransformData(t *testing.T) {
	tests := []struct {
		name     string
		req      *order.CheckAndTransformDataRequest
		expected *order.CheckAndTransformDataResponse
	}{
		{
			name: "Name contains non-English characters",
			req: &order.CheckAndTransformDataRequest{
				Order: &order.OrderModule{
					Id:   "A0001",
					Name: "This Name Contians Non-English Characters! 123",
					Address: &order.AddressModule{
						City:     "taipei-city",
						District: "da-an-district",
						Street:   "fuxing-shouth-road",
					},
					Price:    "1000",
					Currency: "TWD",
				},
			},
			expected: &order.CheckAndTransformDataResponse{
				Message:    "Name contains non-English characters",
				StatusCode: 400,
				Success:    false,
			},
		},
		{
			name: "Name is not capitalized",
			req: &order.CheckAndTransformDataRequest{
				Order: &order.OrderModule{
					Id:   "A0001",
					Name: "Name not Capitalized",
					Address: &order.AddressModule{
						City:     "taipei-city",
						District: "da-an-district",
						Street:   "fuxing-shouth-road",
					},
					Price:    "1000",
					Currency: "TWD",
				},
			},
			expected: &order.CheckAndTransformDataResponse{
				Message:    "Name is not capitalized",
				StatusCode: 400,
				Success:    false,
			},
		},
		{
			name: "Price not a number",
			req: &order.CheckAndTransformDataRequest{
				Order: &order.OrderModule{
					Id:   "A0001",
					Name: "Melody Holiday Inn",
					Address: &order.AddressModule{
						City:     "taipei-city",
						District: "da-an-district",
						Street:   "fuxing-shouth-road",
					},
					Price:    "abc",
					Currency: "TWD",
				},
			},
			expected: &order.CheckAndTransformDataResponse{
				Message:    "Price is not a number",
				StatusCode: 400,
				Success:    false,
			},
		},
		{
			name: "Price is over 2000",
			req: &order.CheckAndTransformDataRequest{
				Order: &order.OrderModule{
					Id:   "A0001",
					Name: "Melody Holiday Inn",
					Address: &order.AddressModule{
						City:     "taipei-city",
						District: "da-an-district",
						Street:   "fuxing-shouth-road",
					},
					Price:    "3000",
					Currency: "TWD",
				},
			},
			expected: &order.CheckAndTransformDataResponse{
				Message:    "Price is over 2000",
				StatusCode: 400,
				Success:    false,
			},
		},
		{
			name: "Currency format is wrong",
			req: &order.CheckAndTransformDataRequest{
				Order: &order.OrderModule{
					Id:   "A0001",
					Name: "Melody Holiday Inn",
					Address: &order.AddressModule{
						City:     "taipei-city",
						District: "da-an-district",
						Street:   "fuxing-shouth-road",
					},
					Price:    "100",
					Currency: "JPY",
				},
			},
			expected: &order.CheckAndTransformDataResponse{
				Message:    "Currency format is wrong",
				StatusCode: 400,
				Success:    false,
			},
		},
		{
			name: "Currency is USD and converted to TWD",
			req: &order.CheckAndTransformDataRequest{
				Order: &order.OrderModule{
					Id:   "A0001",
					Name: "Melody Holiday Inn",
					Address: &order.AddressModule{
						City:     "taipei-city",
						District: "da-an-district",
						Street:   "fuxing-shouth-road",
					},
					Price:    "100",
					Currency: "USD",
				},
			},
			expected: &order.CheckAndTransformDataResponse{
				Order: &order.OrderModule{
					Id:   "A0001",
					Name: "Melody Holiday Inn",
					Address: &order.AddressModule{
						City:     "taipei-city",
						District: "da-an-district",
						Street:   "fuxing-shouth-road",
					},
					Price:    strconv.Itoa(100 * 31),
					Currency: "TWD",
				},
				Message:    "Transform success",
				StatusCode: 200,
				Success:    true,
			},
		},
		{
			name: "All vlidations pass",
			req: &order.CheckAndTransformDataRequest{
				Order: &order.OrderModule{
					Id:   "A0001",
					Name: "Melody Holiday Inn",
					Address: &order.AddressModule{
						City:     "taipei-city",
						District: "da-an-district",
						Street:   "fuxing-shouth-road",
					},
					Price:    "1000",
					Currency: "TWD",
				},
			},
			expected: &order.CheckAndTransformDataResponse{
				Order: &order.OrderModule{
					Id:   "A0001",
					Name: "Melody Holiday Inn",
					Address: &order.AddressModule{
						City:     "taipei-city",
						District: "da-an-district",
						Street:   "fuxing-shouth-road",
					},
					Price:    "1000",
					Currency: "TWD",
				},
				Message:    "Transform success",
				StatusCode: 200,
				Success:    true,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			res := &order.CheckAndTransformDataResponse{}
			nameHandler := &o.NameHandler{}
			priceHandler := &o.PriceHandler{}
			responseHandler := &o.ResponseHandler{}
			currencyHandler := &o.CurrencyHandler{}

			order := o.NewOrder(nameHandler, priceHandler, currencyHandler, responseHandler)
			err := order.CheckAndTransformData(context.Background(), tt.req, res)
			assert.NoError(t, err)
			assert.Equal(t, tt.expected, res)
		})
	}
}
