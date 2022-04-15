module github.com/Fillonj/clamav-http/clamav-http

go 1.13

require (
	github.com/Fillonj/clamav-http/clamav-http/server v0.0.0-00010101000000-000000000000 // indirect
	github.com/dutchcoders/go-clamd v0.0.0-20170520113014-b970184f4d9e
	github.com/konsorten/go-windows-terminal-sequences v1.0.1 // indirect
	github.com/sirupsen/logrus v1.8.1
	github.com/stretchr/objx v0.1.1 // indirect
)

replace github.com/Fillonj/clamav-http/clamav-http/server => ./server
