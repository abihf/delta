module github.com/abihf/delta/v2

go 1.18

retract (
	v2.0.1 // should be tagged as alpha :bow:
	v2.0.0 // api gateway v1 error :bow:
)

require (
	github.com/aws/aws-lambda-go v1.34.1
	github.com/json-iterator/go v1.1.12
)

require (
	github.com/modern-go/concurrent v0.0.0-20180228061459-e0a39a4cb421 // indirect
	github.com/modern-go/reflect2 v1.0.2 // indirect
)
