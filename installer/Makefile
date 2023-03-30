all: build_darwin_arm64 build_darwin_amd64

goreleaser:
	cd .. && goreleaser release --snapshot --clean

icns:
	sh icns.sh

build_darwin_arm64: icns goreleaser
	@echo "build_darwin_arm64"
	sh bundle.sh darwin_arm64 darwin_arm64 && sh -c "appdmg ./darwin_arm64/appdmg.json ../dist/MailCount_ARM_64.dmg"

build_darwin_amd64: icns goreleaser
	@echo "build_darwin_amd64"
	sh bundle.sh darwin_amd64_v1 darwin_amd64 && sh -c "appdmg ./darwin_amd64/appdmg.json ../dist/MailCount_AMD_64.dmg"