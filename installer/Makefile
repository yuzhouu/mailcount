all: build_darwin_arm64

icns:
	sh icns.sh

build_darwin_arm64: icns
	@echo "build_darwin_arm64"
	sh bundle.sh && sh -c "appdmg ./appdmg_arm.json ./MailCount_ARM.dmg"