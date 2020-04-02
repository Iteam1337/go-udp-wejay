module github.com/Iteam1337/go-udp-wejay

go 1.13

require (
	bou.ke/monkey v1.0.2
	github.com/Iteam1337/go-protobuf-wejay v0.0.0-20200402122254-9edb1cf8393a
	github.com/Iteam1337/go-udp-wejay/mocks v0.0.0-00010101000000-000000000000
	github.com/cortesi/modd v0.0.0-20191202231957-98a770274f90 // indirect
	github.com/golang/protobuf v1.3.5
	github.com/joho/godotenv v1.3.0
	github.com/zmb3/spotify v0.0.0-20200112163645-71a4c67d18db
	golang.org/x/oauth2 v0.0.0-20200107190931-bf48bf16ab8d
)

replace github.com/Iteam1337/go-udp-wejay/mocks => ./mocks
