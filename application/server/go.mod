module application

go 1.14

replace github.com/go-kit/kit => github.com/go-kit/kit v0.8.0

replace github.com/ugorji/go => github.com/ugorji/go/codec v1.1.7

replace github.com/gin-gonic/gin => github.com/gin-gonic/gin v1.7.1

replace google.golang.org/grpc => google.golang.org/grpc v1.29.1

require (
	github.com/gin-contrib/cors v1.5.0
	github.com/gin-gonic/gin v1.9.1
	github.com/google/uuid v1.6.0
	github.com/hyperledger/fabric-sdk-go v1.0.0
	github.com/ipfs/go-ipfs-api v0.7.0
	github.com/robfig/cron/v3 v3.0.1
	golang.org/x/crypto v0.14.0
	gorm.io/driver/postgres v1.5.6
	gorm.io/gorm v1.25.7
)
