#!/usr/bin/env bash

set -eu -o pipefail

if [[ ! -z ${DEBUG-} ]]; then
    set -x
fi

: ${PREFIX:=/usr/local}
BINDIR="$PREFIX/bin"

if [[ $# -gt 0 ]]; then
  BINDIR=$1
fi

_can_install() {
  if [[ ! -d "$BINDIR" ]]; then
    mkdir -p "$BINDIR" 2> /dev/null
  fi
  [[ -d "$BINDIR" && -w "$BINDIR" ]]
}

if ! _can_install && [[ $EUID != 0 ]]; then
  echo "Run script as sudo"
  exit 1
fi

if ! _can_install; then
  echo "Can't install to $BINDIR"
  exit 1
fi

machine=$(uname -m)

case $(uname -s) in
    Linux)
        os="Linux"
        ;;
    Darwin)
        os="Darwin"
        ;;
    *)
        echo "OS not supported"
        exit 1
        ;;
esac

latest="$(curl -sL 'https://api.github.com/repos/profclems/glab/releases/latest' | grep 'tag_name' | grep --only 'v[0-9\.]\+' | cut -c 2-)"
echo $machine
curl -sL "https://github.com/profclems/glab/releases/download/v${latest}/glab_${latest}_${os}_${machine}.tar.gz" | tar -C /tmp/ -xzf -
install -m755 /tmp/glab $BINDIR/glab
echo "Successfully installed glab into $BINDIR/"
