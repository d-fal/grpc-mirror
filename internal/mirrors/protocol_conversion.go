package mirrors

import (
	topupService "grpc-mirror/pkg/protobufs/destinationpb"
	pb "grpc-mirror/pkg/protobufs/mirrorpb"
)

// UpliftTopupRequest this func would convert
func UpliftTopupRequest(msg *topupService.TopupRequestType) *pb.TopupRequestType {
	return &pb.TopupRequestType{
		Username: msg.Username,
		Password: msg.Password,
		Msisdn:   msg.Msisdn,
		Pin:      msg.Pin,
		MobNo:    msg.MobNo,
		Amount:   msg.Amount,
		Desc:     msg.Desc,
		Type:     msg.Type,
		AddData:  msg.AddData,
	}
}

// ConvertTopupResponse this func would convert
func ConvertTopupResponse(response *pb.TopupResponseType) *topupService.TopupResponseType {
	return &topupService.TopupResponseType{
		RespCode: response.RespCode,
		Ref:      response.Ref,
		Serial:   response.Serial,
	}
}
