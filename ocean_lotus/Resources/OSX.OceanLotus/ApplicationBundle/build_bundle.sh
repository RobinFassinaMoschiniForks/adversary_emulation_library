#!/usr/bin/env bash

# ---------------------------------------------------------------------------
# build_bundle.sh - Modified appify.sh script from https://gist.github.com/oubiwann/453744744da1141ccc542ff75b47e0cf
#                   and builds an Application Bundle from necessary components

# Copyright 2023 The MITRE Corporation. Approved for public release. Document number CT0005.
# Licensed under the Apache License, Version 2.0 (the "License"); you may not use this file except in compliance with the License. You may obtain a copy of the License at

# http://www.apache.org/licenses/LICENSE-2.0

# Unless required by applicable law or agreed to in writing, software distributed under the License is distributed on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied. See the License for the specific language governing permissions and limitations under the License.

# This project makes use of ATT&CK®
# ATT&CK Terms of Use - https://attack.mitre.org/resources/terms-of-use/

# Usage: ./build_bundle.sh -s [SCRIPT FILE] -i [ICON FILE] -d [DOCUMENT FILE] -p [LAUNCHD FILE] -n [BUNDLE NAME]

# Resources:
#   https://gist.github.com/oubiwann/453744744da1141ccc542ff75b47e0cf

# Revision History:

# ---------------------------------------------------------------------------

VERSION=4.0.1
SCRIPT=`basename "$0"`
APPNAME="My App"
APPICONS="/System/Library/CoreServices/CoreTypes.bundle/Contents/Resources/GenericApplicationIcon.icns"
OSX_VERSION=`sw_vers -productVersion`
PWD=`pwd`

function usage {
	cat <<EOF
$SCRIPT v${VERSION} for for Mac OS X - https://gist.github.com/oubiwann/453744744da1141ccc542ff75b47e0cf
Usage:
  $SCRIPT [options]
Options:
  -h, --help 		Prints this help message, then exits
  -s, --script		Name of the script to 'appify' (required)
  -n, --name 		Name of the application (default "$APPNAME")
  -i, --icons		Name of the icons file to use when creating the app
                        (defaults to $APPICONS)
  -d, --document	Name of the document file to add as a resource
  -p, --persistence	Name of the Launchd plist file to add as PkgInfo
  -v, --version		Prints the version of this script, then exits
Description:
  Creates the simplest possible Mac app from a shell script.
  Appify has one required parameter, the script to appify:
    $SCRIPT --script my-app-script.sh
  Note that you cannot rename appified apps. If you want to give your app
  a custom name, use the '--name' option
    $SCRIPT --script my-app-script.sh --name "Sweet"
Copyright:
  Copyright (c) Thomas Aylott <http://subtlegradient.com/>
  Modified by Mathias Bynens <http://mathiasbynens.be/>
  Modified by Andrew Dvorak <http://OhReally.net/>
  Rewritten by Duncan McGreggor <http://github.com/oubiwann/>
EOF
	exit 1
}

function version {
	echo "v${VERSION}"
	exit 1
}

function error {
	echo
	echo "ERROR: $1"
	echo
	usage
}

while :; do
  case $1 in
    -h | --help )        usage;;
    -s | --script )      APPSCRIPT="$2"; shift ;;
    -n | --name )        APPNAME="$2"; shift ;;
    -i | --icons )       APPICONS="$2"; shift ;;
    -d | --document )    APPDOC="$2"; shift ;;
    -p | --persistence ) APPPERSIST="$2"; shift ;;
    -v | --version )     version;;
    -- )                 shift; break ;;
    * )                  break ;;
  esac
  shift
done

if [ -z ${APPSCRIPT+nil} ]; then
	error "the script to appify must be provided!"
fi

if [ ! -f "$APPSCRIPT" ]; then
	error "the can't find the script '$APPSCRIPT'"
fi

if [ -a "$APPNAME.app" ]; then
    rm -rf "$APPNAME.app"
fi

APPDIR="$APPNAME.app/Contents"

mkdir -vp "$APPDIR"/{MacOS,Resources}
cp -v "$APPICONS" "$APPDIR/Resources/icon.icns"
cp -v "$APPDOC" "$APPDIR/Resources/default.config"
cp -v "$APPSCRIPT" "$APPDIR/MacOS/$APPNAME"
cp -v "$APPPERSIST" "$APPDIR/PkgInfo"
chmod +x "$APPDIR/MacOS/$APPNAME"

cat <<EOF > "$APPDIR/Info.plist"
<?xml version="1.0" encoding="UTF-8"?>
<!DOCTYPE plist PUBLIC "-//Apple//DTD PLIST 1.0//EN" "http://www.apple.com/DTDs/PropertyList-1.0.dtd">
<plist version="1.0">
  <dict>
    <key>CFBundleExecutable</key>
    <string>$APPNAME</string>
    <key>CFBundleGetInfoString</key>
    <string>$APPNAME</string>
    <key>CFBundleIconFile</key>
    <string>icon</string>
    <key>CFBundleName</key>
    <string>$APPNAME</string>
    <key>CFBundlePackageType</key>
    <string>APPL</string>
    <key>CFBundleSignature</key>
    <string>4242</string>
  </dict>
</plist>
EOF

echo "Application bundle created at '$PWD/$APPNAME.app'"

