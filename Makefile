
.PHONY: clean build

build:
	cd cmd/space-glide && go build -o ../../space-glide .

clean:
	$(RM) space-glide
