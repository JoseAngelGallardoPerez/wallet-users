module github.com/Confialink/wallet-users

go 1.13

replace github.com/Confialink/wallet-users/rpc/proto/users => ./rpc/proto/users

require (
	github.com/Confialink/wallet-accounts/rpc/accounts v0.0.0-20210218063536-4e2d21b26af2
	github.com/Confialink/wallet-customization/rpc/proto v0.0.0-20210218070643-6b465594074d
	github.com/Confialink/wallet-files/rpc/files v0.0.0-20210218063711-154e22336663
	github.com/Confialink/wallet-logs/rpc/logs v0.0.0-20210218064020-81b818342efd
	github.com/Confialink/wallet-notifications/rpc/proto/notifications v0.0.0-20210218064438-818cea3b20db
	github.com/Confialink/wallet-permissions/rpc/permissions v0.0.0-20210218064621-7b7ddad868c8
	github.com/Confialink/wallet-pkg-acl v0.0.0-20210218070839-a03813da4b89
	github.com/Confialink/wallet-pkg-discovery/v2 v2.0.0-20210217105157-30e31661c1d1
	github.com/Confialink/wallet-pkg-env_config v0.0.0-20210217112253-9483d21626ce
	github.com/Confialink/wallet-pkg-env_mods v0.0.0-20210217112432-4bda6de1ee2c
	github.com/Confialink/wallet-pkg-errors v1.0.2
	github.com/Confialink/wallet-pkg-list_params v0.0.0-20210217104359-69dfc53fe9ee
	github.com/Confialink/wallet-pkg-model_serializer v0.0.0-20210217111055-c5e1cb1a75c7
	github.com/Confialink/wallet-pkg-utils v0.0.0-20210217112822-e79f6d74cdc1
	github.com/Confialink/wallet-settings/rpc/proto/settings v0.0.0-20210218070334-b4153fc126a0
	github.com/Confialink/wallet-users/rpc/proto/users v0.0.0-00010101000000-000000000000
	github.com/DATA-DOG/go-sqlmock v1.4.1
	github.com/SebastiaanKlippert/go-wkhtmltopdf v1.5.0
	github.com/dgrijalva/jwt-go v3.2.0+incompatible
	github.com/gin-contrib/cors v1.3.1
	github.com/gin-gonic/gin v1.6.3
	github.com/go-playground/validator/v10 v10.2.0
	github.com/google/uuid v1.1.1
	github.com/inconshreveable/log15 v0.0.0-20200109203555-b30bc20e4fd1
	github.com/jasonlvhit/gocron v0.0.0-20200423141508-ab84337f7963
	github.com/jinzhu/gorm v1.9.15
	github.com/nats-io/nats-streaming-server v0.18.0 // indirect
	github.com/nats-io/stan.go v0.7.0
	github.com/ompluscator/dynamic-struct v1.2.0
	github.com/onsi/ginkgo v1.14.0
	github.com/onsi/gomega v1.10.1
	github.com/pkg/errors v0.9.1
	github.com/stretchr/testify v1.6.1
	github.com/twitchtv/twirp v5.12.0+incompatible
	go.uber.org/dig v1.10.0
	golang.org/x/crypto v0.0.0-20200709230013-948cd5f35899
	gopkg.in/go-playground/validator.v8 v8.18.2
)
