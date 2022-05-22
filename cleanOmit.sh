#!/bin/bash

cd ./api/ent
if [[ $OSTYPE == 'darwin'* ]]; 
then
    echo 'executing on macOS'
    sed -i '' 's/,omitempty//g' *.go
else
    echo 'executing on linux'    
    sed -i 's/,omitempty//g' *.go
fi

