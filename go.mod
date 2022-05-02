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
	github.com/confio/ics23/go v0.7.0 // indirect
	github.com/cosmos/cosmos-sdk v0.45.3
	github.com/ethereum/go-ethereum v1.10.16
	github.com/getsentry/sentry-go v0.13.0
	github.com/gin-gonic/gin v1.7.7
	github.com/go-ole/go-ole v1.2.6 // indirect
	github.com/go-playground/universal-translator v0.18.0
	github.com/go-playground/validator/v10 v10.10.0
	github.com/go-redis/redis/v8 v8.11.4
	github.com/gofrs/uuid v4.2.0+incompatible
	github.com/google/go-cmp v0.5.7 // indirect
	github.com/google/gofuzz v1.2.0 // indirect
	github.com/google/uuid v1.3.0 // indirect
	github.com/iamolegga/enviper v1.4.0
	github.com/jackc/pgx/v4 v4.15.0
	github.com/jmoiron/sqlx v1.3.3
	github.com/lib/pq v1.10.4
	github.com/mattn/go-sqlite3 v2.0.3+incompatible // indirect
	github.com/natefinch/lumberjack v2.0.0+incompatible
	github.com/onsi/gomega v1.18.1 // indirect
	github.com/rogpeppe/go-internal v1.8.1 // indirect
	github.com/spf13/viper v1.10.1
	github.com/stretchr/testify v1.7.1
	github.com/tendermint/tm-db v0.6.7 // indirect
	github.com/tklauser/go-sysconf v0.3.9 // indirect
	go.uber.org/zap v1.21.0
	golang.org/x/crypto v0.0.0-20220214200702-86341886e292 // indirect
	golang.org/x/sys v0.0.0-20220209214540-3681064d5158 // indirect
	k8s.io/api v0.21.1
	k8s.io/apimachinery v0.21.1
	k8s.io/client-go v11.0.0+incompatible
	sigs.k8s.io/controller-runtime v0.9.0
)
