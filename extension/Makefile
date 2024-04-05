# The source CSS file
SRC = html/styles/src.css

# The destination for the processed CSS file
OUT = html/styles/dist.css

.PHONY: watch
watch:
	NODE_ENV=development tailwindcss -i $(SRC) -o $(OUT) --watch

.PHONY: build
build:
	NODE_ENV=production tailwindcss -i $(SRC) -o $(OUT)
