#!/usr/bin/env bash


# Convenience function to build various binaries
buildTargetBinary(){
    env GOOS=$1 GOARCH=$2 go build -ldflags="-s -w -X 'main.Version=$GITHUB_REF'" -o ./builds/remp-$1-$2-$GITHUB_REF

    if [[ $1 == "windows" ]]; then
        zip -r ./releases/remp-$1-$2-$GITHUB_REF.zip ./builds/remp-$1-$2-$GITHUB_REF
    else 
        tar -zcvf ./releases/remp-$1-$2-$GITHUB_REF.tar.gz ./builds/remp-$1-$2-$GITHUB_REF
    fi
} 

rm -rf ./builds
rm -rf ./releases
mkdir ./releases

buildTargetBinary linux amd64
buildTargetBinary linux 386
buildTargetBinary linux arm64
buildTargetBinary linux arm

buildTargetBinary darwin arm64
buildTargetBinary darwin amd64

buildTargetBinary windows amd64
buildTargetBinary windows 386