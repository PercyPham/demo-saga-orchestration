module services.order

go 1.16

require (
	github.com/gin-gonic/gin v1.7.2
	github.com/ilyakaznacheev/cleanenv v1.2.5
	github.com/stretchr/testify v1.7.0
	gorm.io/driver/postgres v1.1.0
	gorm.io/gorm v1.21.10
	services.shared v0.0.0
)

replace services.shared v0.0.0 => ../shared
