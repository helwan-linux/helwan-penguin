
pkgname=penguin
pkgver=1.0.0
pkgrel=1
pkgdesc="Helwan Penguin platform game"
arch=('x86_64')
url="https://github.com/helwan-linux/helwan-penguin"
license=('GPL-3.0-or-later')
depends=('glibc' 'gcc-libs')
makedepends=('go' 'git')
source=("$pkgname-$pkgver.tar.gz::https://github.com/helwan-linux/helwan-penguin/archive/refs/tags/v$pkgver.tar.gz")
sha256sums=('SKIP')

build() {
    cd "$srcdir/helwan-penguin-$pkgver"
    export CGO_ENABLED=1
    go build -trimpath -buildmode=pie \
        -ldflags="-s -w -linkmode=external -extldflags=-Wl,-O1,--sort-common,--as-needed,-z,relro,-z,now" \
        -o penguin ./cmd/helwan-penguin
}

package() {
    cd "$srcdir/helwan-penguin-$pkgver"

    install -Dm755 penguin "$pkgdir/usr/bin/penguin"
    install -Dm644 cmd/helwan-penguin/icon.png \
        "$pkgdir/usr/share/icons/hicolor/256x256/apps/penguin.png"
    install -Dm644 packaging/helwan-penguin.desktop \
        "$pkgdir/usr/share/applications/penguin.desktop"
}
