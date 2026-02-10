# Migration Guide: PostgreSQL to MySQL (or vice versa)

Aplikasi ini sekarang mendukung **PostgreSQL** dan **MySQL**. Berikut panduan untuk migrasi antar database.

## Perbedaan Penting

### 1. UUID Generation

#### PostgreSQL

- Menggunakan `gen_random_uuid()` (built-in)
- Otomatis di-generate saat INSERT

#### MySQL

- MySQL 8.0+: Gunakan `UUID()` function atau generate UUID di aplikasi
- Alternatif: Gunakan GORM hooks untuk auto-generate UUID

### 2. Case Sensitivity

#### PostgreSQL

- `ILIKE` untuk case-insensitive search (sudah di-handle oleh GORM)

#### MySQL

- `LIKE` secara default case-insensitive (tergantung collation)
- GORM otomatis handle perbedaan ini

### 3. Connection String Format

#### PostgreSQL

```
postgres://username:password@localhost:5432/database_name?sslmode=disable
```

#### MySQL

```
username:password@tcp(localhost:3306)/database_name?charset=utf8mb4&parseTime=True&loc=Local
```

## Database Schema

### PostgreSQL DDL

```sql
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE categories (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(255) NOT NULL UNIQUE,
    description TEXT
);

CREATE TABLE products (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(255) NOT NULL,
    description TEXT,
    price INTEGER NOT NULL,
    stock INTEGER NOT NULL,
    category_id UUID NOT NULL REFERENCES categories(id),
    picture_url TEXT
);

CREATE TABLE transactions (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    total_amount INTEGER NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE transaction_details (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    transaction_id UUID NOT NULL REFERENCES transactions(id),
    product_id UUID NOT NULL REFERENCES products(id),
    quantity INTEGER NOT NULL,
    subtotal INTEGER NOT NULL
);
```

### MySQL DDL

```sql
CREATE TABLE categories (
    id CHAR(36) PRIMARY KEY DEFAULT (UUID()),
    name VARCHAR(255) NOT NULL UNIQUE,
    description TEXT
);

CREATE TABLE products (
    id CHAR(36) PRIMARY KEY DEFAULT (UUID()),
    name VARCHAR(255) NOT NULL,
    description TEXT,
    price INT NOT NULL,
    stock INT NOT NULL,
    category_id CHAR(36) NOT NULL,
    picture_url TEXT,
    FOREIGN KEY (category_id) REFERENCES categories(id)
);

CREATE TABLE transactions (
    id CHAR(36) PRIMARY KEY DEFAULT (UUID()),
    total_amount INT NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE transaction_details (
    id CHAR(36) PRIMARY KEY DEFAULT (UUID()),
    transaction_id CHAR(36) NOT NULL,
    product_id CHAR(36) NOT NULL,
    quantity INT NOT NULL,
    subtotal INT NOT NULL,
    FOREIGN KEY (transaction_id) REFERENCES transactions(id),
    FOREIGN KEY (product_id) REFERENCES products(id)
);
```

## GORM Auto Migration

Untuk membuat schema otomatis menggunakan GORM:

```go
import (
    "labkoding.my.id/kasir-api/models"
)

// Di main.go setelah InitDB
err = database.AutoMigrate(db,
    &models.Category{},
    &models.Product{},
    &models.Transaction{},
    &models.TransactionDetail{},
)
if err != nil {
    log.Fatal("Failed to auto-migrate:", err)
}
```

**⚠️ CATATAN PENTING:**

- Untuk MySQL, UUID default value `DEFAULT (UUID())` hanya tersedia di MySQL 8.0+
- Untuk versi MySQL lebih lama, UUID akan di-generate oleh GORM saat INSERT
- Pastikan GORM tags di models sudah sesuai (sudah di-setup dengan benar)

## Migrasi Data

### Export dari PostgreSQL

```bash
pg_dump -h localhost -U username -d dbname -t categories -t products -t transactions -t transaction_details --data-only --column-inserts > data.sql
```

### Import ke MySQL

```bash
# Edit data.sql untuk adjust syntax differences
# Kemudian:
mysql -h localhost -u username -p dbname < data.sql
```

## Testing

Setelah switch database:

1. Test koneksi: Jalankan aplikasi dan cek log
2. Test endpoints: Create, Read, Update, Delete
3. Test transaction: Pastikan ACID properties bekerja
4. Test query performance

## Rollback Strategy

Selalu backup database sebelum migrasi:

### PostgreSQL

```bash
pg_dump -h localhost -U username -d dbname > backup.sql
```

### MySQL

```bash
mysqldump -h localhost -u username -p dbname > backup.sql
```
