#!/bin/sh
# Usage: [sudo] [BINDIR=/usr/local/bin] ./install.sh [<BINDIR>]
#
# Example:
#     1. sudo ./install.sh /usr/local/bin
#     2. sudo ./install.sh /usr/bin
#     3. ./install.sh $HOME/usr/bin
#     4. BINDIR=$HOME/usr/bin ./install.sh
#
# Default BINDIR=/usr/bin

set -euf

if [ -n "${DEBUG-}" ]; then
    set -x
fi

: "${BINDIR:=/usr/bin}"

if [ $# -gt 0 ]; then
  BINDIR=$1
fi

_can_install() {
  if [ ! -d "${BINDIR}" ]; then
    mkdir -p "${BINDIR}" 2> /dev/null
  fi
  [ -d "${BINDIR}" ] && [ -w "${BINDIR}" ]
}

if ! _can_install && [ "$(id -u)" != 0 ]; then
  printf "Run script as sudo\n"
  exit 1
fi

if ! _can_install; then
  printf -- "Can't install to %s\n" "${BINDIR}"
  exit 1
fi

machine=$(uname -m)

case $(uname -s) in
    Linux)
        os="Linux"
        ;;
    Darwin)
        os="macOS"
        ;;
    *)
        printf "OS not supported\n"
        exit 1
        ;;
esac

printf "Fetching latest version\n"
latest="$(curl -sL 'https://api.github.com/repos/profclems/glab/releases/latest' | grep 'tag_name' | grep --only 'v[0-9\.]\+' | cut -c 2-)"
tempFolder="/tmp/glab_v${latest}"

printf -- "Found version %s\n" "${latest}"

mkdir -p "${tempFolder}" 2> /dev/null
printf -- "Downloading glab_%s_%s_%s.tar.gz\n" "${latest}" "${os}" "${machine}"
curl -sL "https://github.com/profclems/glab/releases/download/v${latest}/glab_${latest}_${os}_${machine}.tar.gz" | tar -C "${tempFolder}" -xzf -

printf -- "Installing...\n"
install -m755 "${tempFolder}/bin/glab" "${BINDIR}/glab"

printf "Cleaning up temp files\n"
rm -rf "${tempFolder}"

printf -- "Successfully installed glab into %s/\n" "${BINDIR}"
