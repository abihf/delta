package delta

import "os"

func GetDefaultTransformer() Transformer {
	switch os.Getenv("DELTA_MODE") {
	case "alb":
		return AlbTransformer{}

	case "apigwv2":
		return ApiGatewayV2Transformer{}

	case "lambdaurl":
		return LambdaURLTransformer{}

	default:
		return ApiGatewayV1Transformer{}
	}
}
