package order

import pb "orderService/proto/order"

type ResponseHandler struct{}

func (r *ResponseHandler) SetResponseFailure(res *pb.CheckAndTransformDataResponse, message string) {
	res.Message = message
	res.StatusCode = 400
	res.Order = nil
	res.Success = false
}

func (r *ResponseHandler) SetResponseSuccess(req *pb.CheckAndTransformDataRequest, res *pb.CheckAndTransformDataResponse, message string) {
	res.Order = req.Order
	res.Message = "Transform success"
	res.StatusCode = 200
	res.Success = true
}
