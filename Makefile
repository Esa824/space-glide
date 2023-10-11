
.PHONY: clean build race run

build:
	cd cmd/space-glide && go build -o ../../space-glide .

clean:
	$(RM) space-glide

race:
	cd cmd/space-glide && go build -race -o ../../space-glide .

run:
	cd cmd/space-glide && go build -o ../../space-glide .
	./space-glide

