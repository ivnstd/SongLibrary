#!/bin/bash
cd ./song_library || exit
go run ./cmd/main.go &
cd ../music_info_mock || exit
go run ./main.go