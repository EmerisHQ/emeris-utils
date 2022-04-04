module github.com/emerishq/emeris-utils

go 1.16

replace (
	github.com/gogo/protobuf => github.com/regen-network/protobuf v1.3.3-alpha.regen.1
	github.com/jmoiron/sqlx => github.com/abraithwaite/sqlx v1.3.2-0.20210331022513-df9bf9884350
	google.golang.org/grpc => google.golang.org/grpc v1.33.2
	k8s.io/client-go => k8s.io/client-go v0.21.1
)

require (
	github.com/alicebob/miniredis/v2 v2.18.0
	github.com/allinbits/starport-operator v0.0.1-alpha.45
	github.com/cockroachdb/cockroach-go/v2 v2.2.8
	github.com/cosmos/cosmos-sdk v0.42.8
	github.com/cosmos/gaia/v5 v5.0.4
	github.com/ethereum/go-ethereum v1.10.15
	github.com/gin-gonic/gin v1.7.7
	github.com/go-playground/universal-translator v0.18.0
	github.com/go-playground/validator/v10 v10.10.1
	github.com/go-redis/redis/v8 v8.11.4
	github.com/gofrs/uuid v4.2.0+incompatible
	github.com/iamolegga/enviper v1.4.0
	github.com/jackc/pgx/v4 v4.15.0
	github.com/jmoiron/sqlx v1.3.3
	github.com/lib/pq v1.10.4
	github.com/mattn/go-sqlite3 v2.0.3+incompatible // indirect
	github.com/natefinch/lumberjack v2.0.0+incompatible
	github.com/spf13/viper v1.10.1
	github.com/stretchr/testify v1.7.1-0.20210427113832-6241f9ab9942
	go.uber.org/zap v1.21.0
	k8s.io/api v0.21.1
	k8s.io/apimachinery v0.21.1
	k8s.io/client-go v11.0.0+incompatible
	sigs.k8s.io/controller-runtime v0.9.0
)
