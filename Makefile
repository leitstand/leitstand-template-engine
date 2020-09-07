ui: build-swagger build-ui

build-swagger:
	@echo "[run] build-swagger"
	./generateSwagger.sh

build-ui:
	@echo "[run] build-ui"
	./generateStatik.sh

build-docker-image:
	docker build --force-rm -t leitstand-template-engine .

run-docker-container:
	docker stop leistand-template-engine || true
	docker rm leitstand-template-engine || true
	docker run -d --restart unless-stopped --name leitstand-template-engine \
		-p 8082:80 leitstand-template-engine
	# To mount the templates folder from the host
    # Just add
	#    --mount type=bind,source=$(pwd)/templates,target=/etc/rtbrick/leitstand-template-engine/templates \

linters:
	@echo "[run] linters"
	@# We do this instead of a simple `go fmt ...` because (at least in the
	@# begining) it's better too see the changes than blindly run it.
	@echo "gofmt -d -e ./cmd/ ./pkg/"; \
		fmt_out=`gofmt -d -e ./cmd/ ./pkg/` || exit 1; \
		[ -z "$$fmt_out" ] || { \
			echo "$$fmt_out"; \
			echo "#"; \
			echo "# If you want a quick fix just run: go fmt ."; \
			echo "#"; \
			exit 1; \
		};
	@which golint > /dev/null || { \
		echo "#"; \
		echo "# Either you don't have golint installed or it's not accessible."; \
		echo "#"; \
		echo "# Make sure you have \$$GOPATH set up correctly and that \$$GOPATH/bin is included in your \$$PATH,"; \
		echo "# see https://golang.org/doc/code.html#GOPATH & https://github.com/golang/go/wiki/GOPATH ."; \
		echo "#"; \
		echo "# After that run: go get -u golang.org/x/lint/golint"; \
		echo "# see https://github.com/golang/lint ."; \
		echo "#"; \
		exit 1; \
	};
	golint -set_exit_status ./cmd/... ./pkg/...
	go vet ./cmd/... ./pkg/...

build_for_major_platforms:
	env GOOS=windows GOARCH=386 go build -o ./bin/windows_386/template-engine-test.exe ./cmd/template-engine-test/
	env GOOS=windows GOARCH=amd64 go build -o ./bin/windows_amd64/template-engine-test.exe ./cmd/template-engine-test/
	env GOOS=linux GOARCH=amd64 go build -o ./bin/linux_amd64/template-engine-test ./cmd/template-engine-test/
	env GOOS=linux GOARCH=386 go build -o ./bin/linux_386/template-engine-test ./cmd/template-engine-test/
	env GOOS=darwin GOARCH=amd64 go build -o ./bin/darwin_amd64/template-engine-test ./cmd/template-engine-test/
