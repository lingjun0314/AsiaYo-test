package controllers

import (
	"fmt"
	"ginapigateway/models"
	pb "ginapigateway/proto/order"
	"net/http"

	"github.com/gin-gonic/gin"
)

type ValidatoinController struct{}

func (con *ValidatoinController) CheckDataType(ctx *gin.Context) {
	var order models.Order
	//	Check data valid or not
	if err := ctx.ShouldBindJSON(&order); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid JSON",
		})
		return
	}

	//	Call service
	orderService := pb.NewOrderService("order", models.MicroClient)
	res, _ := orderService.CheckAndTransformData(ctx, &pb.CheckAndTransformDataRequest{
		Order: &pb.OrderModule{
			Id:   order.Id,
			Name: order.Name,
			Address: &pb.AddressModule{
				City:     order.Address.City,
				District: order.Address.District,
				Street:   order.Address.Street,
			},
			Price:    order.Price,
			Currency: order.Currency,
		},
	})
	fmt.Println(res.Success)

	ctx.JSON(int(res.StatusCode), gin.H{
		"order":   res.Order,
		"message": res.Message,
		"success": res.Success,
	})
}
