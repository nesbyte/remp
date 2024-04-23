#!/usr/bin/env bash


# Convenience function to build various binaries
buildTargetBinary(){
    env GOOS=$1 GOARCH=$2 go build -ldflags="-s -w -X 'main.Version=$GITHUB_REF_NAME'" -o ./builds/remp-$1-$2-$GITHUB_REF_NAME
    tar -zcvf ./releases/remp-$1-$2.tar.gz -C ./builds remp-$1-$2-$GITHUB_REF_NAME
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