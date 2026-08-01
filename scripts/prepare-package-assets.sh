#!/usr/bin/env bash

set -euo pipefail

mkdir -p build/completions build/scripts

sed \
    -e 's/__start_gc/__start_gitcode/g' \
    -e 's/__gc_/__gitcode_/g' \
    -e 's/__start_gitcode gc/__start_gitcode gitcode/g' \
    -e 's/ gc$/ gitcode/g' \
    -e 's/ for gc / for gitcode /g' \
    completions/gc.bash > build/completions/gitcode.bash

sed \
    -e 's/#compdef gc/#compdef gitcode/' \
    -e 's/compdef _gc gc/compdef _gitcode gitcode/' \
    -e 's/_gc/_gitcode/g' \
    -e 's/__gc_/__gitcode_/g' \
    -e 's/ for gc / for gitcode /g' \
    completions/gc.zsh > build/completions/gitcode.zsh

sed \
    -e 's/__gc_/__gitcode_/g' \
    -e 's/ for gc / for gitcode /g' \
    -e 's/-c gc/-c gitcode/g' \
    completions/gc.fish > build/completions/gitcode.fish

sed 's/\r$//' scripts/postinstall.sh > build/scripts/postinstall.sh
chmod 755 build/scripts/postinstall.sh
