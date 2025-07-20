# Search Engine Service v2.0.0

🚀 **Enterprise-Grade Search Engine Service**

---

## ⚡ Hızlı Başlangıç (Her Adımı Takip Edin, Hatasız Çalışır!)

### 1. Gereksinimler
- Go 1.21+ (https://go.dev/dl/)
- MySQL 8.0+ (https://dev.mysql.com/downloads/installer/)
- Git
- (Opsiyonel) Docker & Docker Compose

### 2. Go Kurulumu ve PATH Ayarı (Windows için)
- Go'yu yükledikten sonra, terminale şunu yazın:
  ```powershell
  $env:PATH = "C:\Program Files\Go\bin;" + $env:PATH
  go version
  ```
- `go version` çıktısı görmelisiniz.

### 3. Projeyi Klonlayın
```bash
git clone <repository-url>
cd search-engine-service
```

### 4. MySQL Kullanıcı ve Veritabanı Oluşturun (Hem development hem test için)
```sql
-- MySQL'e root ile bağlanın:
mysql -u root -p

-- Ana veritabanı ve kullanıcı:
CREATE DATABASE search_engine CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;
CREATE USER 'search_user'@'localhost' IDENTIFIED BY 'search_password';
GRANT ALL PRIVILEGES ON search_engine.* TO 'search_user'@'localhost';

-- Test için (integration testleri için):
CREATE DATABASE test_db CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;
CREATE USER 'test_user'@'localhost' IDENTIFIED BY 'test_password';
GRANT ALL PRIVILEGES ON test_db.* TO 'test_user'@'localhost';
FLUSH PRIVILEGES;
```

### 5. .env ve DB Bilgileri Nerelerde Güncellenmeli?
- **Anahtar dosya:** `.env` (proje kök dizininde)
- **Docker ile çalışıyorsanız:** `docker-compose.yml` içindeki `environment` kısmı
- **Testler için:** `.env` dosyasında test DB bilgileri de olmalı (örn. TEST_DB_USER, TEST_DB_PASSWORD, TEST_DB_NAME veya testlerinizde kullanılan DB_USER/DB_NAME)

#### Örnek .env
```env
DB_HOST=localhost
DB_PORT=3306
DB_USER=search_user
DB_PASSWORD=search_password
DB_NAME=search_engine

# Testler için (gerekirse)
TEST_DB_USER=test_user
TEST_DB_PASSWORD=test_password
TEST_DB_NAME=test_db

PROVIDER_JSON_URL=http://localhost:3001/api/videos
PROVIDER_XML_URL=http://localhost:3001/api/articles
```

> **Not:** DB bilgilerini değiştirdiğinizde hem ana uygulama hem testler hem de Docker ortamı için aynı bilgileri kullandığınızdan emin olun.

### 6. Bağımlılıkları Yükleyin
```bash
go mod download
go mod tidy
```

### 7. Veritabanı Şemasını ve Örnek Verileri Yükleyin
```bash
mysql -u search_user -p search_engine < scripts/init.sql
mysql -u search_user -p search_engine < scripts/example_data.sql
```

### 8. Uygulamayı Başlatın
```bash
# Mock server'ı başlatın (ayrı terminalde):
go run cmd/mock-server/main.go

# Ana uygulamayı başlatın:
go run cmd/server/main.go
```

### 9. Testleri Çalıştırın
```bash
# Tüm testler (unit + integration):
go test ./... -v
# Sadece unit testler:
go test ./internal/services/... -v
# Sadece integration testler:
go test ./tests/integration/... -v
```

### 10. Sık Karşılaşılan Hatalar ve Çözümleri
- **go : The term 'go' is not recognized...**
  - Çözüm: Go'yu yükleyin ve yukarıdaki PATH komutunu uygulayın.
- **Access denied for user 'test_user'@'localhost'**
  - Çözüm: MySQL'de test_user ve test_db oluşturun, şifreyi .env ile eşleştirin.
- **Port already in use**
  - Çözüm: 8080 veya 3001 portunda başka uygulama çalışıyorsa kapatın.
- **Cannot connect to MySQL**
  - Çözüm: MySQL servisinin çalıştığından emin olun.

---

## 🐳 Docker ile Tek Komutla Çalıştırmak İçin
```bash
docker-compose up -d
# http://localhost:8080 (API), http://localhost:3001 (Mock), http://localhost:8080/dashboard
```

---

## 📚 Diğer Bilgiler
- API dokümantasyonu, örnek istekler ve gelişmiş kullanım için aşağıya bakın.
- Tüm adımları eksiksiz uygularsanız sistem **hatasız** çalışır.

---

## 🔗 Test için Kullanabileceğiniz Temel URL’ler

| Amaç                | URL                                               |
|---------------------|--------------------------------------------------|
| API ana endpoint    | http://localhost:8080                            |
| Sağlık kontrolü     | http://localhost:8080/health                     |
| Arama               | http://localhost:8080/api/search?q=golang&type=video |
| Popüler içerik      | http://localhost:8080/api/content/popular        |
| Mock server (JSON)  | http://localhost:3001/api/videos                 |
| Mock server (XML)   | http://localhost:3001/api/articles               |
| Dashboard           | http://localhost:8080/dashboard                  |

---

---

## 🏗️ Proje Mimarisi

### Clean Architecture
```
├── cmd/                    # Uygulama giriş noktaları
│   ├── server/            # Ana API sunucusu
│   └── mock-server/       # Test için mock provider
├── internal/              # İç paketler
│   ├── api/              # HTTP handlers ve middleware
│   ├── database/         # Veritabanı modelleri ve repository
│   ├── di/               # Dependency injection container
│   ├── providers/        # Veri sağlayıcıları (JSON/XML)
│   ├── services/         # İş mantığı katmanı
│   └── utils/            # Yardımcı fonksiyonlar
├── scripts/              # Veritabanı scriptleri
├── tests/                # Test dosyaları
└── docs/                 # Dokümantasyon
```

### Teknoloji Stack'i
- **Backend**: Go 1.21+ with Gin framework
- **Database**: MySQL 8.0+ with GORM ORM
- **Architecture**: Clean Architecture with Dependency Injection
- **Testing**: Go testing + Testify + Integration tests
- **Security**: Rate limiting, CORS, input sanitization
- **Logging**: Structured logging with Zap
- **Containerization**: Docker & Docker Compose

---

## 🔍 API Dokümantasyonu

### Temel Endpoints

#### 1. Health Check
```http
GET /health
```
**Response:**
```json
{
  "status": "healthy",
  "timestamp": "2024-01-15T10:30:00Z",
  "version": "2.0.0",
  "database": "connected",
  "providers": ["video_provider", "article_provider"]
}
```

#### 2. Arama API
```http
GET /api/search?q={query}&type={content_type}&page={page}&limit={limit}
```

**Parametreler:**
- `q` (required): Arama terimi
- `type` (optional): `video` veya `text`
- `page` (optional): Sayfa numarası (default: 1)
- `limit` (optional): Sayfa başına sonuç (default: 10)

**Örnek İstek:**
```http
GET /api/search?q=golang&type=video&page=1&limit=5
```

**Response:**
```json
{
  "query": "golang",
  "total_results": 15,
  "page": 1,
  "limit": 5,
  "total_pages": 3,
  "results": [
    {
      "id": 1,
      "title": "Go Programming Tutorial",
      "description": "Learn Go programming from scratch",
      "type": "video",
      "url": "https://example.com/video1",
      "published_at": "2024-01-10T15:30:00Z",
      "views": 15000,
      "likes": 450,
      "reading_time": 0,
      "reactions": 0,
      "scores": {
        "base_score": 19.5,
        "type_multiplier": 1.5,
        "freshness_score": 3.0,
        "engagement_score": 0.3,
        "final_score": 32.55
      }
    }
  ]
}
```

#### 3. Popüler İçerik
```http
GET /api/content/popular?limit={limit}
```

**Response:**
```json
{
  "popular_content": [
    {
      "id": 1,
      "title": "Most Popular Video",
      "type": "video",
      "final_score": 95.5,
      "views": 50000,
      "likes": 1200
    }
  ]
}
```

#### 4. İçerik Detayı
```http
GET /api/content/{id}
```

#### 5. Provider Bilgileri
```http
GET /api/providers
```

---

## 🎯 Puan Hesaplama Algoritması

### Video İçerik Puanı
```
Base Score = (views / 1000) + (likes / 100)
Type Multiplier = 1.5x
Freshness Score = Yayın tarihine göre (1 hafta=5, 1 ay=3, 3 ay=1, eski=0)
Engagement Score = (likes / views) × 10
Final Score = (Base Score × Type Multiplier) + Freshness Score + Engagement Score
```

### Metin İçerik Puanı
```
Base Score = reading_time + (reactions / 50)
Type Multiplier = 1.0x
Freshness Score = Video ile aynı
Engagement Score = (reactions / reading_time) × 5
Final Score = (Base Score × Type Multiplier) + Freshness Score + Engagement Score
```

### Örnek Hesaplama
**Video:**
- Views: 15,000, Likes: 450, Published: 5 gün önce
- Base Score: (15000/1000) + (450/100) = 15 + 4.5 = 19.5
- Type Multiplier: 1.5
- Freshness Score: 5.0 (1 hafta içinde)
- Engagement Score: (450/15000) × 10 = 0.3
- Final Score: (19.5 × 1.5) + 5.0 + 0.3 = 32.55

---

## 🔧 Gelişmiş Konfigürasyon

### Environment Variables
```env
# Database
DB_HOST=localhost
DB_PORT=3306
DB_USER=search_user
DB_PASSWORD=search_password
DB_NAME=search_engine

# Server
SERVER_PORT=8080
SERVER_HOST=localhost
ENVIRONMENT=development

# Logging
LOG_LEVEL=debug

# Security
JWT_SECRET=your-secret-key-here
RATE_LIMIT_REQUESTS=100
RATE_LIMIT_WINDOW=1m

# Providers
VIDEO_PROVIDER_URL=http://localhost:3001/api/videos
ARTICLE_PROVIDER_URL=http://localhost:3002/api/articles
```

### Docker Compose Konfigürasyonu
```yaml
version: '3.8'
services:
  mysql:
    image: mysql:8.0
    environment:
      MYSQL_ROOT_PASSWORD: root_password
      MYSQL_DATABASE: search_engine
      MYSQL_USER: search_user
      MYSQL_PASSWORD: search_password
    ports:
      - "3306:3306"
    volumes:
      - mysql_data:/var/lib/mysql

  app:
    build: .
    ports:
      - "8080:8080"
    environment:
      - DB_HOST=mysql
      - DB_USER=search_user
      - DB_PASSWORD=search_password
      - DB_NAME=search_engine
    depends_on:
      - mysql

volumes:
  mysql_data:
```

---

## 🧪 Test Stratejisi

### Test Kategorileri
1. **Unit Tests**: Her servis ve fonksiyon için
2. **Integration Tests**: API endpoint'leri ve veritabanı
3. **Benchmark Tests**: Performans testleri

### Test Komutları
```bash
# Tüm testler
go test ./... -v

# Sadece unit testler
go test ./internal/services/... -v

# Sadece integration testler
go test ./tests/integration/... -v

# Benchmark testler
go test -bench=. ./internal/services/

# Test coverage
go test -cover ./...

# Race condition testleri
go test -race ./...
```

### Test Veritabanı
- Ayrı test veritabanı kullanılır
- Her test öncesi temizlenir
- Mock provider ile test edilir

---

## 🔒 Güvenlik Özellikleri

### Middleware Katmanı
1. **Rate Limiting**: IP başına istek sınırı
2. **CORS**: Cross-origin resource sharing
3. **Security Headers**: XSS, CSRF koruması
4. **Input Sanitization**: Girdi temizleme
5. **Request ID Tracking**: İstek takibi
6. **Error Recovery**: Hata yakalama

### Güvenlik Headers
```
X-Content-Type-Options: nosniff
X-Frame-Options: DENY
X-XSS-Protection: 1; mode=block
Strict-Transport-Security: max-age=31536000; includeSubDomains
Content-Security-Policy: default-src 'self'
```

---

## 📊 Monitoring ve Logging

### Structured Logging
```go
logger.Info("Search request processed",
    zap.String("query", query),
    zap.String("type", contentType),
    zap.Int("results", len(results)),
    zap.Duration("duration", duration),
)
```

### Health Check Endpoints
- `/health`: Genel sağlık durumu
- `/ready`: Uygulama hazır mı?
- `/live`: Uygulama çalışıyor mu?

### Metrics (Gelecek Özellik)
- Request sayısı
- Response süreleri
- Error oranları
- Database bağlantı durumu

---

## 🚀 Deployment

### Production Checklist
- [ ] Environment variables ayarlanmış
- [ ] Database migration'ları çalıştırılmış
- [ ] SSL sertifikası yüklenmiş
- [ ] Monitoring aktif
- [ ] Backup stratejisi hazır
- [ ] Load balancer konfigürasyonu
- [ ] Security audit tamamlanmış

### Docker Deployment
```bash
# Production build
docker build -t search-engine:latest .

# Run with environment
docker run -d \
  -p 8080:8080 \
  -e DB_HOST=your-db-host \
  -e DB_USER=your-db-user \
  -e DB_PASSWORD=your-db-password \
  search-engine:latest
```

---

## 🤝 Katkıda Bulunma

### Geliştirme Ortamı Kurulumu
1. Fork yapın
2. Feature branch oluşturun
3. Değişikliklerinizi commit edin
4. Testleri çalıştırın
5. Pull request gönderin

### Kod Standartları
- Go fmt kullanın
- Lint kurallarına uyun
- Test coverage %80+ olmalı
- Documentation güncelleyin

---

## 📝 Changelog

### v2.0.0 (2024-01-15)
- ✨ Clean Architecture implementasyonu
- 🔒 Güvenlik middleware'leri eklendi
- 🧪 Kapsamlı test suite
- 📊 Puan hesaplama algoritması
- 🐳 Docker desteği
- 📚 Detaylı dokümantasyon

### v1.0.0 (2024-01-10)
- 🎉 İlk sürüm
- 🔍 Temel arama fonksiyonalitesi
- 📦 Provider entegrasyonu

---

## 📄 Lisans

Bu proje MIT lisansı altında lisanslanmıştır. Detaylar için [LICENSE](LICENSE) dosyasına bakın.

---

## 👥 Ekip

- **Geliştirici**: [Adınız]
- **Email**: [email@example.com]
- **GitHub**: [github.com/kullaniciadi]

---

## 🙏 Teşekkürler

- [Gin Framework](https://github.com/gin-gonic/gin)
- [GORM](https://gorm.io/)
- [Uber Zap](https://github.com/uber-go/zap)
- [Testify](https://github.com/stretchr/testify)

---

**⭐ Bu projeyi beğendiyseniz yıldız vermeyi unutmayın!** 