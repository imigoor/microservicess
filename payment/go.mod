module github.com/huseyinbabal/microservices/payment

go 1.25.0

require (
	github.com/ruandg/microservices-proto/golang/payment v0.0.0-00010101000000-000000000000
	github.com/sirupsen/logrus v1.9.0
	google.golang.org/grpc v1.81.1
	gorm.io/driver/mysql v1.5.2
	gorm.io/gorm v1.25.5
)

require (
	github.com/go-sql-driver/mysql v1.7.0 // indirect
	github.com/jinzhu/inflection v1.0.0 // indirect
	github.com/jinzhu/now v1.1.5 // indirect
	golang.org/x/net v0.51.0 // indirect
	golang.org/x/sys v0.42.0 // indirect
	golang.org/x/text v0.34.0 // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20260226221140-a57be14db171 // indirect
	google.golang.org/protobuf v1.36.11 // indirect
)

replace github.com/ruandg/microservices-proto/golang/payment => ../../microservices-proto/golang/payment
