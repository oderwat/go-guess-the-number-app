build:
	GOARCH=wasm GOOS=js go build -o web/app.wasm
	go build

run: build
	sleep 2 && open http://127.0.0.1:8000 &
	./go-guess-the-number-app
