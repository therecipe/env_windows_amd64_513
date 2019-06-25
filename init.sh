#!/bin/bash

set -ev

QT_ROOT=$HOME/.wine/drive_c
QT_VERSION=5.13.0

rm -rf ./${QT_VERSION}
rm -rf ./Licenses
rm -rf ./Tools

rsync -avz $QT_ROOT/Qt/Qt${QT_VERSION}/${QT_VERSION}/mingw73_64 ./${QT_VERSION}/
rsync -avz $QT_ROOT/Qt/Qt${QT_VERSION}/Licenses .
rsync -avz $QT_ROOT/Qt/Qt${QT_VERSION}/Tools . #using 7.1.0, 7.2.0, 7.3.0 and 8.1.0 will make cgo take ages to generate code
#GCC=x86_64-6.4.0-release-posix-seh-rt_v5-rev0.7z
#curl -sL -O https://sourceforge.net/projects/mingw-w64/files/Toolchains%20targetting%20Win64/Personal%20Builds/mingw-builds/6.4.0/threads-posix/seh/$GCC
#mkdir -p ./Tools && cd ./Tools && p7zip -d ../$GCC && mv mingw64 mingw730_64 && cd .. && rm -f $GCC

rm -rf ./${QT_VERSION}/mingw73_64/{doc,phrasebooks}
rm -rf ./${QT_VERSION}/mingw73_64/lib/{cmake,pkgconfig,libQt5Bootstrap.a}


for v in *.jsc *.log *.pro *.pro.user *.qmake.stash *.qmlc .DS_Store *_debug* *.la *.prl; do
	find . -name ${v} -exec rm -rf {} \;
done

mkdir -p ./${QT_VERSION}/mingw73_64/_bin
for v in qtenv2.bat windeployqt.exe lconvert.exe moc.exe qmake.exe qmlcachegen.exe qmlimportscanner.exe qt.conf rcc.exe uic.exe *.dll; do
	mv ./${QT_VERSION}/mingw73_64/bin/${v} ./${QT_VERSION}/mingw73_64/_bin/
done
rm -rf ./${QT_VERSION}/mingw73_64/bin && mv ./${QT_VERSION}/mingw73_64/_bin ./${QT_VERSION}/mingw73_64/bin

find . -type f -name "*d.dll*" -exec sh -c 'test -e $(echo {} | sed -e "s/d\./\./g") && rm -rf {}' \;
find . -type f -name "*d.a" -exec sh -c 'test -e $(echo {} | sed -e "s/d\./\./g") && rm -rf {}' \;

echo "module github.com/therecipe/env_windows_amd64_513/Tools" > ./Tools/go.mod
echo "package Tools" > ./Tools/mod.go

set +e
mv $QT_ROOT/Qt $QT_ROOT/Qt_orig
set -e

wine go run ./patch.go



du -sh ./5*

#wine $(wine go env GOPATH)/bin/qtsetup
