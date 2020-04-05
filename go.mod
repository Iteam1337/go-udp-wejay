module github.com/Iteam1337/go-udp-wejay

go 1.13

require (
	bou.ke/monkey v1.0.2
	github.com/Iteam1337/go-protobuf-wejay v0.0.0-20200404185428-a12951023ef6
	github.com/Iteam1337/go-udp-wejay/mocks v0.0.0-00010101000000-000000000000
	github.com/ankjevel/spotify v0.0.0-20200403101354-52db1de203c1
	github.com/golang/protobuf v1.3.5
	github.com/joho/godotenv v1.3.0
	golang.org/x/net v0.0.0-20191126235420-ef20fe5d7933 // indirect
	golang.org/x/oauth2 v0.0.0-20200107190931-bf48bf16ab8d
)

replace github.com/Iteam1337/go-udp-wejay/mocks => ./mocks
