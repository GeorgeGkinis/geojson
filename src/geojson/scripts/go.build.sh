CGO_ENABLED=0 GOOS=linux GOARCH=amd64
go build -a -tags netgo -ldflags '-w -extldflags "-static"' -o ../service-A/service-A ../service-A/service-A.go
go build -a -tags netgo -ldflags '-w -extldflags "-static"' -o ../service-B/service-B ../service-B/service-B.go
go build -a -tags netgo -ldflags '-w -extldflags "-static"' -o ../service-C/service-C ../service-C/service-C.go