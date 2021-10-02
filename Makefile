web/app.wasm: main.go
	GOARCH=wasm GOOS=js go build -o web/app.wasm

app: main.go
	go build -o app

run: app web/app.wasm
	./app

watch:
	sleep 2 && open http://127.0.0.1:8000 &
	ulimit -n 1000 # increasing watchable files
	# I use clearing the scrollback buffer here so it is easier to
	# relate to errors for the last compiled version
	reflex -s -r '.*\.go' -- sh -c "printf '\033[2J\033[3J\033[1;1H' && make run"
