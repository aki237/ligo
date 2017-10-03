#!/bin/bash

OUT="$HOME/ligo"
Black="\e[1;30m"
Blue="\e[1;34m"
Green="\e[1;32m"
Cyan="\e[1;36m"
Red="\e[1;31m"
Purple="\e[1;35m"
Brown="\e[1;33m"
Sky="\e[1;34m"
Green="\e[1;32m"
Cyan="\e[1;36m"
Red="\e[1;31m"
Purple="\e[1;35m"
Brown="\e[1;33m"
Reset="\e[m"


function buildPackage () {
    go install github.com/aki237/ligo/pkg/ligo/ && echo "Built the base ligo library" "Built the ligo package"
}

function buildPlugins () {
    for i in $(ls -1 --color=no);do
        if [ -d "$i" ]; then
            PKG=$(basename $i)
            [[ -d "$OUT/lib/$PKG" ]] || mkdir -p "$OUT/lib/$PKG"
            echo -en Building $Red$PKG$Reset...
            cd $i && go build -buildmode=plugin -o $OUT/lib/$PKG/$PKG.plg && cd ..
            cp $i/*.lg $OUT/lib/$PKG/
            echo " Done"
        fi
    done
}


function build () {
    [[ -d "$OUT/lib/" ]] || mkdir -p $OUT

    buildPlugins
    exit;
}


build
