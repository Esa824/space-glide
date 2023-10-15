
.PHONY: clean build race

build:
	cd cmd/space-glide && go build -o ../../space-glide .

clean:
	$(RM) space-glide

race:
	cd cmd/space-glide && go build -race -o ../../space-glide .
