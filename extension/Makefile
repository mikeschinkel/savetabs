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

.PHONY: install
install:
	npm install tailwindcss daisyui --save-dev

build:

.PHONY: install
integrity:
ifndef FILE
	$(error FILE is undefined. Please specify FILE like so: make integrity FILE=./path/to/file.js)
endif
	@echo && \
	sha384="$$(shasum -a 384 "$(FILE)" | awk '{print $1}' | xxd -r -p | base64)" && \
	echo "Add the following to your HTML file, adjusting your <script src> as necessary:" && \
	echo && \
	echo "   <script src=\"$(FILE)\" integrity=\"sha384-$${sha384}\" crossorigin=\"anonymous\"></script>" && \
	echo