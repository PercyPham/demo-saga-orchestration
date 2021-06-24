module services.payment_contract

go 1.16

require (
	services.shared v0.0.0
)

replace (
	services.shared v0.0.0 => ../shared
)
