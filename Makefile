VERSION := 0.9.0
REPO := registry.wklive.net:5000/dg/pass-go
TAG := ${REPO}:${VERSION}

LOCALE_DIR := locales
MESSAGES := $(LOCALE_DIR)/messages.pot
TEMPLATE_DIR := templates
TEMPLATE_FILES := $(shell find $(TEMPLATE_DIR) -name '*.html')
TEMPLATE_TRANSLATIONS := $(TEMPLATE_FILES:.html=.py)
LOCALE_FILES := $(shell find $(LOCALE_DIR) -name '*.po')


.PHONY: build clean docker run push translate 

# Convert golang templates into a format pybabel can understand. 
# So very, very sorry about this.
%.py : %.html
	sed -n "s/^.*{{[ ]*gettext \(\".*\"\)[ ]*}}.*$$/gettext(\1)/p" $< > $@

$(MESSAGES) : $(TEMPLATE_TRANSLATIONS) venv
	$(VENV)/pybabel extract -F babel.cfg --no-location --no-wrap --omit-header -o $(MESSAGES) $(TEMPLATE_DIR)
	rm templates/*.py

$(LOCALE_FILES): $(MESSAGES) venv
	$(VENV)/pybabel update --omit-header -i $(MESSAGES) -d $(LOCALE_DIR)
	$(VENV)/pybabel compile -f -d $(LOCALE_DIR)

translate: $(MESSAGES) $(LOCALE_FILES)


clean:
	rm -rf .venv
	
build:
	rm -f pass-go pkged.go
	# Cross-compile for linux
	pkger
	GOOS=linux GOARCH=amd64 go build -a -tags netgo -ldflags '-s -w -extldflags "-static"' .

docker: build
	docker build --progress=plain --no-cache --squash -t ${TAG}  .

run: docker
	docker run --rm -it -p 5000:5000 ${TAG}

include Makefile.venv