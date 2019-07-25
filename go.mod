module github.com/busgo/forest

go 1.12

require (
	github.com/coreos/bbolt v1.3.3 // indirect
	github.com/coreos/etcd v3.3.13+incompatible
	github.com/coreos/go-semver v0.3.0 // indirect
	github.com/coreos/go-systemd v0.0.0-20190620071333-e64a0ec8b42a // indirect
	github.com/coreos/pkg v0.0.0-20180928190104-399ea9e2e55f // indirect
	github.com/dgrijalva/jwt-go v3.2.0+incompatible // indirect
	github.com/golang/groupcache v0.0.0-20190702054246-869f871628b6 // indirect
	github.com/gorilla/websocket v1.4.0 // indirect
	github.com/grpc-ecosystem/go-grpc-middleware v1.0.0 // indirect
	github.com/grpc-ecosystem/go-grpc-prometheus v1.2.0 // indirect
	github.com/grpc-ecosystem/grpc-gateway v1.9.4 // indirect
	github.com/jonboulle/clockwork v0.1.0 // indirect
	github.com/labstack/echo v3.3.10+incompatible
	github.com/labstack/gommon v0.2.9
	github.com/prometheus/client_golang v1.0.0 // indirect
	github.com/prometheus/common v0.4.1
	github.com/robfig/cron v1.2.0
	github.com/soheilhy/cmux v0.1.4 // indirect
	github.com/tmc/grpc-websocket-proxy v0.0.0-20190109142713-0ad062ec5ee5 // indirect
	github.com/xiang90/probing v0.0.0-20190116061207-43a291ad63a2 // indirect
	go.etcd.io/bbolt v1.3.3 // indirect
	go.etcd.io/etcd v3.3.13+incompatible
	go.uber.org/atomic v1.4.0 // indirect
	go.uber.org/multierr v1.1.0 // indirect
	go.uber.org/zap v1.10.0 // indirect
)

replace (
	cloud.google.com/go => github.com/googleapis/google-cloud-go v0.41.0

	golang.org/x/crypto => github.com/golang/crypto v0.0.0-20190701094942-4def268fd1a4

	golang.org/x/exp => github.com/golang/exp v0.0.0-20190627132806-fd42eb6b336f

	golang.org/x/image => github.com/golang/image v0.0.0-20190703141733-d6a02ce849c9

	golang.org/x/lint => github.com/golang/lint v0.0.0-20190409202823-959b441ac422

	golang.org/x/mobile => github.com/golang/mobile v0.0.0-20190607214518-6fa95d984e88

	golang.org/x/net => github.com/golang/net v0.0.0-20190628185345-da137c7871d7

	golang.org/x/oauth2 => github.com/golang/oauth2 v0.0.0-20190604053449-0f29369cfe45

	golang.org/x/sync => github.com/golang/sync v0.0.0-20190423024810-112230192c58

	golang.org/x/sys => github.com/golang/sys v0.0.0-20190626221950-04f50cda93cb

	golang.org/x/text => github.com/golang/text v0.3.2

	golang.org/x/time => github.com/golang/time v0.0.0-20190308202827-9d24e82272b4

	golang.org/x/tools => github.com/golang/tools v0.0.0-20190709211700-7b25e351ac0e

	google.golang.org/api => github.com/googleapis/google-api-go-client v0.7.0

	google.golang.org/appengine => github.com/golang/appengine v1.6.1

	google.golang.org/genproto => github.com/googleapis/go-genproto v0.0.0-20190708153700-3bdd9d9f5532

	google.golang.org/grpc => github.com/grpc/grpc-go v1.22.0

)
