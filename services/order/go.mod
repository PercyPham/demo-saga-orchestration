module services.order

go 1.16

require (
	github.com/gin-gonic/gin v1.7.2
	github.com/stretchr/testify v1.7.0
	services.shared v0.0.0
)

replace services.shared v0.0.0 => ../shared
