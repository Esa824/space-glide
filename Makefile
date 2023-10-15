
.PHONY: clean build race

build:
	cd cmd/space-glide && go build -gcflags "all=-N -l" -o ../../space-glide .

clean:
	$(RM) space-glide

race:
	cd cmd/space-glide && go build -race -o ../../space-glide .
