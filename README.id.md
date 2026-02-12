# HeadControl

Dashboard admin Headscale yang minimalis, dibangun dengan Go dan HTMX.

[Read in English](README.md)

Project ini dibuat untuk pemakaian pribadi. Saya memilih Go dan HTMX sebagai
fondasi karena saya butuh sesuatu yang benar-benar ringan dan cepat, tanpa
framework JavaScript berat, tanpa build pipeline yang rumit, cukup satu
binary yang melayani semuanya. Kalau kamu merasa ini berguna atau ingin
memodifikasinya untuk kebutuhanmu sendiri, silakan.

*Disclaimer: saya tidak begitu jago Go, jadi kalau baca kodenya dan bertanya
"kenapa begini?", jawabannya kemungkinan besar "karena yang ini yang berhasil jalan."*

---

## Stack

- **Go** — HTTP server, API client, template rendering, penyimpanan SQLite
- **HTMX** — partial page update tanpa menulis JavaScript
- **Lucide Icons** — icon set yang dimuat via CDN
- **SQLite** — penyimpanan konfigurasi lokal

## Fitur

- Koneksi ke Headscale instance manapun via API key
- Dashboard dengan statistik node dan user
- Manajemen user (buat, rename, hapus)
- Manajemen node (rename, expire, hapus, tags, routes)
- Detail view per node
- 16 tema warna bawaan
- Layout responsif (desktop, tablet, mobile)

---

## Persyaratan

- Go 1.21 atau lebih baru
- GCC (dibutuhkan oleh go-sqlite3, lihat catatan di bawah)
- Server Headscale yang sudah berjalan dengan API key

### GCC di Windows

Driver `go-sqlite3` membutuhkan CGO. Di Windows, install
[MSYS2](https://www.msys2.org/) atau [TDM-GCC](https://jmeubank.github.io/tdm-gcc/),
lalu pastikan `gcc` tersedia di PATH.

---

## Instalasi

Clone repository:

```
git clone https://github.com/ahmadzip/headcontrol.git
cd headcontrol
```

Install dependency:

```
go mod tidy
```

Build:

```
go build .
```

Perintah ini menghasilkan `headcontrol.exe` di Windows atau `headcontrol` di Linux/macOS.

Jalankan:

```
./headcontrol
```

Server berjalan di `http://localhost:8080` secara default.

### Flag command-line

| Flag | Default | Keterangan |
|------|---------|------------|
| `-port` | `8080` | Port server |
| `-db` | `headcontrol.db` | Path database SQLite |

Contoh:

```
./headcontrol -port 3000 -db /data/headcontrol.db
```

---

## Pertama 

1. Buka `http://localhost:8080` di browser.
2. Kamu akan diarahkan ke halaman setup.
3. Masukkan URL server Headscale (contoh: `https://headscale.example.com`).
4. Masukkan API key (dibuat dengan `headscale apikeys create`).
5. Klik "Test Connection" untuk verifikasi.
6. Klik "Save" untuk masuk ke dashboard.

---

## Development

Untuk hot-reload saat development, install [Air](https://github.com/air-verse/air):

```
go install github.com/air-verse/air@latest
```

Lalu jalankan:

```
air
```

Air memantau perubahan file dan rebuild secara otomatis.
Konfigurasi ada di `.air.toml`.

---

## Struktur Project

```
headcontrol/
  main.go                          entrypoint, registrasi route
  internal/
    handler/
      handler.go                   struct inti, template engine, middleware
      helpers.go                   render helper, format waktu
      setup.go                     handler halaman setup
      dashboard.go                 handler halaman dashboard
      users.go                     handler manajemen user
      nodes.go                     handler manajemen node
      settings.go                  handler halaman settings
    headscale/
      client.go                    API client headscale
    model/
      models.go                    struktur data
    store/
      store.go                     layer penyimpanan SQLite
  templates/
    layout/layout.html             layout dasar dengan sidebar
    pages/                         template halaman penuh
    partials/                      template partial HTMX
  static/
    css/app.css                    stylesheet utama
    css/theme/                     file tema warna
    js/app.js                      logika client-side
```

---

## Tema

HeadControl menyediakan 16 tema bawaan. Ganti tema lewat selector di
navigasi atas. Pilihan kamu disimpan di localStorage.

---

## Saran dan Kritik

Dibuat di waktu istirahat. Saran, kritik, atau pull request diterima.

---

## Lisensi

Gunakan sesuka hati.
