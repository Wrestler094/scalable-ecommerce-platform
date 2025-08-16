module github.com/Wrestler094/scalable-ecommerce-platform/notification-service

go 1.23.0

require (
	github.com/Wrestler094/scalable-ecommerce-platform/pkg v0.0.0
	github.com/caarlos0/env/v11 v11.3.1
	github.com/go-chi/chi/v5 v5.2.2
	github.com/prometheus/client_golang v1.23.0
	github.com/segmentio/kafka-go v0.4.48
)

replace github.com/Wrestler094/scalable-ecommerce-platform/pkg => ../pkg

require (
	github.com/beorn7/perks v1.0.1 // indirect
	github.com/cespare/xxhash/v2 v2.3.0 // indirect
	github.com/google/uuid v1.6.0 // indirect
	github.com/klauspost/compress v1.18.0 // indirect
	github.com/munnerz/goautoneg v0.0.0-20191010083416-a7dc8b61c822 // indirect
	github.com/pierrec/lz4/v4 v4.1.22 // indirect
	github.com/prometheus/client_model v0.6.2 // indirect
	github.com/prometheus/common v0.65.0 // indirect
	github.com/prometheus/procfs v0.17.0 // indirect
	golang.org/x/sys v0.34.0 // indirect
	google.golang.org/protobuf v1.36.7 // indirect
)
