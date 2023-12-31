root	:=		$(shell dirname $(realpath $(lastword $(MAKEFILE_LIST))))

.PHONY: clean build deploy publish

clean:
	rm -rfv bin
	$(MAKE) -C "${root}/management" clean
	$(MAKE) -C "${root}/management/api" clean

build:
	mkdir -p bin
	scripts/create_template.sh
	GOOS=linux GOARCH=arm64 go build -ldflags="-s -w" -o bin/bootstrap
	$(MAKE) -C "${root}/management" build
	$(MAKE) -C "${root}/management/api" build

deploy:
	sam package --output-template-file "${root}"/packaged.yml --s3-bucket "${bucket}"
	sam deploy --stack-name "${stack}" --capabilities CAPABILITY_IAM --template-file "${root}/packaged.yml"

publish:
	sam package --output-template-file "${root}"/packaged.yml --s3-bucket "${bucket}"
	sam publish --template "${root}/packaged.yml" --region "${AWS_DEFAULT_REGION}"
