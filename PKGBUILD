pkgname=penguin
pkgver=1.0
pkgrel=2
pkgdesc="Helwan Penguin platform game"
arch=('x86_64')
url="https://github.com/helwan-linux/helwan-penguin"
license=('MIT')

depends=(
    'glibc'
    'gcc-libs'
    'alsa-lib'
    'mesa'
    'libx11'
    'libxcursor'
    'libxrandr'
    'libxinerama'
    'libxi'
    'libgl'
)

makedepends=(
    'go'
    'git'
    'pkgconf'
)

source=("git+https://github.com/helwan-linux/helwan-penguin.git#branch=main")
sha256sums=('SKIP')

build() {
    cd "$srcdir/helwan-penguin"

    export CGO_ENABLED=1

    go build \
        -trimpath \
        -buildmode=pie \
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

    sed -i \
        -e 's/^Exec=.*/Exec=penguin/' \
        -e 's/^Icon=.*/Icon=penguin/' \
        "$pkgdir/usr/share/applications/penguin.desktop"

    grep -q '^Exec=' "$pkgdir/usr/share/applications/penguin.desktop" || \
        echo 'Exec=penguin' >> "$pkgdir/usr/share/applications/penguin.desktop"

    grep -q '^Icon=' "$pkgdir/usr/share/applications/penguin.desktop" || \
        echo 'Icon=penguin' >> "$pkgdir/usr/share/applications/penguin.desktop"
}

