
pkgname=penguin
pkgver=0.1.0
pkgrel=1
pkgdesc="Helwan Penguin platform game"
arch=('x86_64')
url="https://github.com/helwan-linux/helwan-penguin"
license=('GPL-3.0-or-later')
depends=('glibc' 'gcc-libs' 'alsa-lib' 'mesa' 'libx11' 'libxcursor' 'libxrandr' 'libxinerama' 'libxi' 'libgl')
makedepends=('go' 'git' 'pkgconf')

source=("git+https://github.com/helwan-linux/helwan-penguin.git")
sha256sums=('SKIP')

pkgver() {
    cd "$srcdir/helwan-penguin"
    printf '0.1.0.r%s.%s' \
        "$(git rev-list --count HEAD)" \
        "$(git rev-parse --short HEAD)"
}

build() {
    cd "$srcdir/helwan-penguin"

    export CGO_ENABLED=1
    export GOFLAGS="-buildmode=pie"

    go build -trimpath \
        -ldflags="-s -w" \
        -o penguin \
        ./cmd/helwan-penguin
}

package() {
    cd "$srcdir/helwan-penguin"

    install -Dm755 penguin \
        "$pkgdir/usr/bin/penguin"

    install -Dm644 cmd/helwan-penguin/icon.png \
        "$pkgdir/usr/share/icons/hicolor/256x256/apps/penguin.png"

    install -Dm644 packaging/helwan-penguin.desktop \
        "$pkgdir/usr/share/applications/penguin.desktop"
}
