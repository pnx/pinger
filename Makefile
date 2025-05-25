
GO=go

pinger : main.go
	$(GO) build -o $@ $^
