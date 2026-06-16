# User Management API

REST API สำหรับจัดการข้อมูลผู้ใช้งาน สร้างด้วย Go + Gin และเชื่อมต่อ PostgreSQL ผ่าน `pgxpool` โดยมีฟีเจอร์หลักสำหรับสร้าง อ่าน แก้ไข ลบ และ login ด้วย `user_id` / `password`

## Tech Stack

- Go `1.25.0`
- Gin HTTP framework
- PostgreSQL
- pgx / pgxpool
- Viper สำหรับโหลด config และอ่าน environment variable
- Zap logger
- Goose-style SQL migration

## Project Structure

```text
cmd/api/                  entrypoint สำหรับรัน HTTP API
internal/app/             bootstrap server, router, graceful shutdown
internal/config/          embedded application config
internal/db/              PostgreSQL connection pool
internal/userprofile/     handler, service, repository, model ของ user profile
internal/validation/      JSON binding และ custom validation
migrations/               SQL migration สำหรับ schema
```

## Configuration

แอปโหลด config จากไฟล์ที่ embed อยู่ที่ `internal/config/config.yaml`

ค่าหลักที่ใช้งาน:

```yaml
apiServer:
  address: 0.0.0.0:8080

posgresql:
  host: localhost
  port: 5432
  username: sa
  password: psdsystem
  database: appdb
  sslmode: disable
```

สามารถ override ผ่าน environment variable ได้ตาม key ของ config เช่น:

```powershell
$env:APISERVER_ADDRESS="0.0.0.0:8081"
$env:POSGRESQL_HOST="localhost"
$env:POSGRESQL_PORT="5432"
```

> หมายเหตุ: key ใน config ปัจจุบันสะกดเป็น `posgresql` ตามโค้ด จึงควรใช้ชื่อนี้ให้ตรงกัน

## Database

สร้าง database ให้ตรงกับ config ก่อนรันแอป:

```sql
CREATE DATABASE appdb;
```

จากนั้นรัน migration ใน `migrations/001_create_userprofile.sql` เพื่อสร้างตาราง `userprofile`

ถ้าใช้ `psql`:

```powershell
psql "postgres://sa:psdsystem@localhost:5432/appdb?sslmode=disable" -f migrations/001_create_userprofile.sql
```

หรือถ้าใช้ `goose`:

```powershell
goose -dir migrations postgres "postgres://sa:psdsystem@localhost:5432/appdb?sslmode=disable" up
```

## Run

ติดตั้ง dependency:

```powershell
go mod download
```

รัน API:

```powershell
go run .\cmd\api\.
```

หรือใช้ make:

```powershell
make run
```

เมื่อรันสำเร็จ API จะเปิดที่:

```text
http://localhost:8080
```

## Health Check

```http
GET /health
```

ตัวอย่าง:

```powershell
curl http://localhost:8080/health
```

Response:

```json
{
  "status": "UP",
  "timestamp": "2026-06-16T14:00:00Z",
  "uptime": "10s"
}
```

## User Profile API

Base path:

```text
/api/v1/userprofiles
```

| Method | Path | Description |
| --- | --- | --- |
| `POST` | `/api/v1/userprofiles` | สร้าง user profile |
| `GET` | `/api/v1/userprofiles` | ดูรายการ user profile ทั้งหมด |
| `GET` | `/api/v1/userprofiles/{user_id}` | ดู user profile ตาม `user_id` |
| `PUT` | `/api/v1/userprofiles/{user_id}` | แก้ไข user profile |
| `DELETE` | `/api/v1/userprofiles/{user_id}` | ลบ user profile |
| `POST` | `/api/v1/userprofiles/login` | login ด้วย `user_id` และ `password` |

### Create User

```powershell
curl -X POST http://localhost:8080/api/v1/userprofiles `
  -H "Content-Type: application/json" `
  -d '{
    "user_id": "bob01",
    "password": "secret123",
    "first_name": "Bob",
    "last_name": "Five",
    "address": "Bangkok",
    "birthdate": "1995-01-15",
    "email": "bob@example.com"
  }'
```

### List Users

```powershell
curl http://localhost:8080/api/v1/userprofiles
```

### Get User By ID

```powershell
curl http://localhost:8080/api/v1/userprofiles/bob01
```

### Update User

```powershell
curl -X PUT http://localhost:8080/api/v1/userprofiles/bob01 `
  -H "Content-Type: application/json" `
  -d '{
    "password": "new-secret123",
    "first_name": "Bob",
    "last_name": "Five",
    "address": "Chiang Mai",
    "birthdate": "1995-01-15",
    "email": "bob.five@example.com"
  }'
```

### Delete User

```powershell
curl -X DELETE http://localhost:8080/api/v1/userprofiles/bob01
```

### Login

```powershell
curl -X POST http://localhost:8080/api/v1/userprofiles/login `
  -H "Content-Type: application/json" `
  -d '{
    "user_id": "bob01",
    "password": "secret123"
  }'
```

## Request Validation

- `user_id` จำเป็น และยาวไม่เกิน 20 ตัวอักษร
- `password` จำเป็น และยาวไม่เกิน 20 ตัวอักษร
- `first_name` จำเป็น และยาวไม่เกิน 150 ตัวอักษร
- `last_name` จำเป็น และยาวไม่เกิน 150 ตัวอักษร
- `birthdate` ใช้รูปแบบ `YYYY-MM-DD` และต้องมีอายุอย่างน้อย 18 ปี
- `email` ต้องเป็น email format ที่ถูกต้อง และยาวไม่เกิน 255 ตัวอักษร

## Response Format

Success response:

```json
{
  "timestamp": "2026-06-16T21:00:00+07:00",
  "error_code": "200",
  "error_message": "Success",
  "data": {
    "user_id": "bob01",
    "first_name": "Bob",
    "last_name": "Five",
    "address": "Bangkok",
    "birthdate": "1995-01-15",
    "email": "bob@example.com",
    "created_at": "2026-06-16T21:00:00+07:00",
    "updated_at": "2026-06-16T21:00:00+07:00"
  }
}
```

Error response:

```json
{
  "timestamp": "2026-06-16T21:00:00+07:00",
  "error_code": "400",
  "error_message": "Field validation Error",
  "error_detail": {
    "email": "[email] Email invalid format"
  }
}
```

## Test

รัน test ทั้งหมด:

```powershell
go test ./...
```
