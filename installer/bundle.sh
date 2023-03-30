#!/bin/bash

# Create the directory structure
mkdir -p MailCount.app/Contents/MacOS
mkdir -p MailCount.app/Contents/Resources

# Copy the command line program to the MacOS directory
cp ../dist/mailcount_darwin_arm64/mailcount MailCount.app/Contents/MacOS/MailCount

# Copy icons to Resources dir
cp icons/icon.icns MailCount.app/Contents/Resources/icon.icns

# Create the Info.plist file in the Resources directory
cat << EOF > MailCount.app/Contents/Info.plist
<?xml version="1.0" encoding="UTF-8"?>
<!DOCTYPE plist PUBLIC "-//Apple//DTD PLIST 1.0//EN"
    "http://www.apple.com/DTDs/PropertyList-1.0.dtd">
<plist version="1.0">
    <dict>
        <key>CFBundleExecutable</key>
        <string>MailCount</string>
        <key>CFBundleIdentifier</key>
        <string>com.example.MailCount</string>
        <key>CFBundleName</key>
        <string>MailCount</string>
        <key>CFBundleVersion</key>
        <string>1.0</string>
        <key>CFBundleIconFile</key>
        <string>icon</string>
    </dict>
</plist>
EOF

# Set the executable bit on the MyApp file
chmod +x MailCount.app/Contents/MacOS/MailCount