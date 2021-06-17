module services.kitchen

go 1.16

require (
	github.com/gin-gonic/gin v1.7.2
	github.com/go-playground/validator/v10 v10.6.1 // indirect
	github.com/golang/protobuf v1.5.2 // indirect
	github.com/ilyakaznacheev/cleanenv v1.2.5
	github.com/json-iterator/go v1.1.11 // indirect
	github.com/leodido/go-urn v1.2.1 // indirect
	github.com/lib/pq v1.6.0 // indirect
	github.com/mattn/go-isatty v0.0.13 // indirect
	github.com/percypham/saga-go v0.0.0
	github.com/streadway/amqp v1.0.0
	github.com/ugorji/go v1.2.6 // indirect
	golang.org/x/crypto v0.0.0-20210513164829-c07d793c2f9a // indirect
	golang.org/x/sys v0.0.0-20210615035016-665e8c7367d1 // indirect
	golang.org/x/text v0.3.6 // indirect
	gopkg.in/yaml.v2 v2.4.0 // indirect
	gorm.io/driver/postgres v1.1.0
	gorm.io/gorm v1.21.11
	services.shared v0.0.0
)

replace (
	github.com/percypham/saga-go v0.0.0 => ../../../saga-go
	services.shared v0.0.0 => ../shared
)
