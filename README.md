# kasir-api

API sederhana untuk mengelola kategori dan produk, dibuat dengan Go + chi dan PostgreSQL.

## Ringkasan

- Bahasa: Go
- Router: chi
- Database: PostgreSQL (implementasi saat ini menggunakan fitur `RETURNING` dan `json_agg`)
- Environment variables:
  - `PORT` — port server (contoh: `3000`)
  - `DB_CONN` — connection string ke database Postgres (contoh: `postgres://user:pass@host:port/dbname?sslmode=disable`)

Menjalankan server (contoh):

```bash
export PORT=3000
export DB_CONN="postgres://user:pass@host:5432/dbname?sslmode=disable"
go run main.go
```

Semua endpoint yang menerima body JSON harus mengirim header:

- `Content-Type: application/json`

---

## Struktur Endpoints (base: http://localhost:{PORT})

1. Health

- GET `/`
  - Deskripsi: health check / welcome
  - Response: plain text `welcome`
  - Contoh:
    ```bash
    curl http://localhost:3000/
    # welcome
    ```

---

2. Categories

a) GET `/categories`

- Deskripsi: Ambil semua kategori. Setiap kategori menyertakan `products` (array) dan `product_count`.
- Response contoh:

```json
[
  {
    "id": "60a974b9-ee9e-4fe7-80cc-4331d41ad275",
    "name": "Minuman",
    "description": "Kategori minuman",
    "product_count": 2,
    "products": [
      {
        "id": "11111111-2222-3333-4444-555555555555",
        "name": "Teh Botol",
        "description": "Teh manis",
        "price": 5000,
        "stock": 10,
        "category_id": "60a974b9-ee9e-4fe7-80cc-4331d41ad275",
        "category_name": "Minuman"
      }
    ]
  }
]
```

b) POST `/categories`

- Deskripsi: Buat kategori baru. Mengembalikan `id` kategori yang baru dibuat.
- Request body:

```json
{
  "name": "Makanan Ringan",
  "description": "Snack dan cemilan"
}
```

- Response contoh:

```json
{ "id": "e6f0a2d1-xxxx-xxxx-xxxx-xxxxxxxxxxxx" }
```

- Catatan: Implementasi repo menggunakan `RETURNING id` (Postgres) untuk mendapatkan id.

c) GET `/categories/{id}`

- Deskripsi: Ambil detail kategori berdasarkan `id`.
- Response contoh:

```json
{
  "id": "60a974b9-ee9e-4fe7-80cc-4331d41ad275",
  "name": "Minuman",
  "description": "Kategori minuman"
}
```

d) PUT `/categories/update/{id}`

- Deskripsi: Update kategori berdasarkan `id`.
- Request body:

```json
{
  "name": "Minuman & Jus",
  "description": "Minuman dingin dan jus"
}
```

- Response: Objek kategori yang diupdate (handler saat ini meng-encode objek hasil update).

e) DELETE `/categories/delete/{id}`

- Deskripsi: Hapus kategori berdasarkan `id`.
- Response contoh:

```json
{ "message": "Category deleted successfully" }
```

---

3. Products

a) GET `/products`

- Deskripsi: Ambil semua produk.
- Response contoh:

```json
[
  {
    "id": "11111111-2222-3333-4444-555555555555",
    "name": "Teh Botol",
    "description": "Teh manis",
    "price": 5000,
    "stock": 10,
    "category_id": "60a974b9-ee9e-4fe7-80cc-4331d41ad275",
    "category_name": "Minuman"
  }
]
```

b) POST `/products`

- Deskripsi: Buat produk baru.
- Request body contoh:

```json
{
  "name": "Teh Botol",
  "description": "Teh manis",
  "price": 5000,
  "stock": 10,
  "category_id": "60a974b9-ee9e-4fe7-80cc-4331d41ad275"
}
```

- Response: Handler saat ini meng-encode object produk yang diterima. Jika ingin ID dikembalikan, perlu menyesuaikan repo/service untuk menggunakan `RETURNING id`.

c) GET `/products/{id}`

- Deskripsi: Ambil produk berdasarkan `id`.
- Response contoh:

```json
{
  "id": "11111111-2222-3333-4444-555555555555",
  "name": "Teh Botol",
  "description": "Teh manis",
  "price": 5000,
  "stock": 10,
  "category_id": "60a974b9-ee9e-4fe7-80cc-4331d41ad275",
  "category_name": "Minuman"
}
```

d) PUT `/products/update/{id}`

- Deskripsi: Update produk berdasarkan `id` (ID diambil dari path param).
- Request body contoh:

```json
{
  "name": "Teh Botol (Update)",
  "description": "Teh manis dingin",
  "price": 5500,
  "stock": 15,
  "category_id": "60a974b9-ee9e-4fe7-80cc-4331d41ad275"
}
```

- Contoh curl:

```bash
curl -X PUT http://localhost:3000/products/update/60a974b9-ee9e-4fe7-80cc-4331d41ad275 \
  -H "Content-Type: application/json" \
  -d '{"name":"Teh Botol Baru","description":"Deskripsi","price":6000,"stock":20,"category_id":"60a974b9-ee9e-4fe7-80cc-4331d41ad275"}'
```

e) DELETE `/products/delete/{id}`

- Deskripsi: Hapus produk berdasarkan `id`.
- Response contoh:

```json
{ "message": "produk berhasil dihapus" }
```

---

4. Transactions

a) POST `/transactions/checkout`

- Deskripsi: Buat transaksi baru.
- Request body contoh:

```json
{
  "items": [
    {
      "product_id": "60a974b9-ee9e-4fe7-80cc-4331d41ad275",
      "quantity": 1
    }
  ]
}
```

---

4. Reports

a) GET `/reports/today`

- Deskripsi: Ambil ringkasan laporan untuk hari ini.
- Response contoh:

```json
{
  "total_revenue": 150000,
  "total_transactions": 12,
  "best_selling_products": {
    "name": "Teh Botol",
    "qty_sold": 20
  }
}
```

b) GET `/reports/range?start=YYYY-MM-DD&end=YYYY-MM-DD`

- Deskripsi: Ambil ringkasan laporan untuk rentang tanggal (inklusif). Parameter `start` dan `end` harus dalam format `YYYY-MM-DD`.
- Contoh:

```bash
curl "http://localhost:3000/reports/range?start=2026-01-01&end=2026-01-31"
```

- Response contoh:

```json
{
  "total_revenue": 4500000,
  "total_transactions": 120,
  "best_selling_products": {
    "name": "Nasi Goreng",
    "qty_sold": 150
  }
}
```
