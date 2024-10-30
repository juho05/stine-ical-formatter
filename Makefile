OUT_DIR=bin
BIN_NAME=stine-ical-formatter

.PHONY: build
build: tailwind-build
	go build -o ${OUT_DIR}/${BIN_NAME} .

.PHONY: init
init:
	go mod download
	npm install
	make tailwind-build

.PHONY: tailwind-watch
tailwind-watch:
	npx tailwindcss build -o web/static/css/tailwind.css --watch

.PHONY: tailwind-build
tailwind-build:
	npx tailwindcss build -o web/static/css/tailwind.css --minify

.PHONY: go-watch
go-watch:
	go run -mod=mod github.com/bokwoon95/wgo run -file 'web/static/.*' -file 'web/templates/.*' .

.PHONY: watch
watch: tailwind-watch go-watch

.PHONY: run
run: tailwind-build
	go run .

.PHONY: clean
clean:
	go clean
	rm -r ${OUT_DIR}