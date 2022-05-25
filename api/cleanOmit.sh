#!/bin/sh
# remove all "omitempty" from the generated json types 
# such that false is not always removed

cd ./ent
if [[ $OSTYPE == 'darwin'* ]]; 
then
    echo 'executing on macOS'
    sed -i '' 's/,omitempty//g' *.go
else
    echo 'executing on linux'    
    sed -i 's/,omitempty//g' *.go
fi

