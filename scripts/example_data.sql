-- =====================================================
-- Search Engine Service - Örnek Veri Dosyası
-- =====================================================
-- Bu dosya veritabanı şemasını, örnek verileri ve 
-- yararlı sorguları içerir.
-- =====================================================

-- Veritabanını oluştur
CREATE DATABASE IF NOT EXISTS search_engine_service 
CHARACTER SET utf8mb4 
COLLATE utf8mb4_unicode_ci;

USE search_engine_service;

-- =====================================================
-- TABLO ŞEMALARI
-- =====================================================

-- İçerik tablosu
CREATE TABLE IF NOT EXISTS contents (
    id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    provider VARCHAR(100) NOT NULL COMMENT 'Veri sağlayıcısı (json_provider, xml_provider)',
    provider_id VARCHAR(255) NOT NULL COMMENT 'Sağlayıcıdaki benzersiz ID',
    title VARCHAR(500) NOT NULL COMMENT 'İçerik başlığı',
    description TEXT COMMENT 'İçerik açıklaması',
    url VARCHAR(1000) NOT NULL COMMENT 'İçerik URL\'i',
    type ENUM('video', 'text') NOT NULL COMMENT 'İçerik türü',
    views INT UNSIGNED DEFAULT 0 COMMENT 'Görüntülenme sayısı',
    likes INT UNSIGNED DEFAULT 0 COMMENT 'Beğeni sayısı',
    duration INT UNSIGNED DEFAULT 0 COMMENT 'Süre (saniye)',
    tags VARCHAR(1000) COMMENT 'Etiketler (virgülle ayrılmış)',
    language VARCHAR(10) DEFAULT 'tr' COMMENT 'Dil kodu',
    published_at TIMESTAMP NULL COMMENT 'Yayınlanma tarihi',
    final_score DECIMAL(10,4) DEFAULT 0.0000 COMMENT 'Hesaplanan final skor',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP COMMENT 'Oluşturulma tarihi',
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT 'Güncellenme tarihi',
    
    -- İndeksler
    INDEX idx_provider_provider_id (provider, provider_id),
    INDEX idx_type (type),
    INDEX idx_final_score (final_score DESC),
    INDEX idx_published_at (published_at DESC),
    INDEX idx_views (views DESC),
    INDEX idx_likes (likes DESC),
    FULLTEXT idx_search (title, description, tags),
    
    -- Benzersizlik kısıtlaması
    UNIQUE KEY uk_provider_provider_id (provider, provider_id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='İçerik tablosu';

-- =====================================================
-- ÖRNEK VERİLER
-- =====================================================

-- Video içerikleri (JSON Provider'dan)
INSERT INTO contents (provider, provider_id, title, description, url, type, views, likes, duration, tags, language, published_at, final_score) VALUES
('json_provider', 'video_001', 'Go Programlama Dili Temelleri', 'Bu videoda Go programlama dilinin temel kavramlarını öğreneceksiniz. Değişkenler, fonksiyonlar, struct\'lar ve daha fazlası...', 'https://example.com/videos/go-basics', 'video', 15420, 892, 1800, 'go,golang,programlama,backend', 'tr', '2024-01-15 10:30:00', 85.2340),
('json_provider', 'video_002', 'Docker ile Containerization', 'Docker kullanarak uygulamalarınızı nasıl containerize edeceğinizi öğrenin. Dockerfile, docker-compose ve best practices...', 'https://example.com/videos/docker-guide', 'video', 8920, 456, 2400, 'docker,container,devops,deployment', 'tr', '2024-01-20 14:15:00', 72.1560),
('json_provider', 'video_003', 'MySQL Veritabanı Optimizasyonu', 'MySQL veritabanınızın performansını nasıl artıracağınızı öğrenin. İndeksler, sorgu optimizasyonu ve monitoring...', 'https://example.com/videos/mysql-optimization', 'video', 12340, 678, 2100, 'mysql,database,optimization,performance', 'tr', '2024-01-25 09:45:00', 78.8900),
('json_provider', 'video_004', 'RESTful API Tasarım Prensipleri', 'Modern ve ölçeklenebilir REST API\'ler nasıl tasarlanır? HTTP metodları, status kodları ve best practices...', 'https://example.com/videos/rest-api-design', 'video', 18760, 1023, 2700, 'api,rest,web,backend,design', 'tr', '2024-02-01 16:20:00', 92.4560),
('json_provider', 'video_005', 'Microservices Mimarisi', 'Microservices mimarisinin avantajları ve dezavantajları. Service discovery, load balancing ve monitoring...', 'https://example.com/videos/microservices', 'video', 9870, 534, 3300, 'microservices,architecture,distributed,scalability', 'tr', '2024-02-05 11:10:00', 76.2340);

-- Metin içerikleri (XML Provider'dan)
INSERT INTO contents (provider, provider_id, title, description, url, type, views, likes, duration, tags, language, published_at, final_score) VALUES
('xml_provider', 'article_001', 'Yapay Zeka ve Gelecek', 'Yapay zeka teknolojilerinin gelecekte hayatımızı nasıl değiştireceğini inceleyen kapsamlı bir analiz. Machine learning, deep learning ve etik konular...', 'https://example.com/articles/ai-future', 'text', 23450, 1234, 0, 'ai,artificial-intelligence,machine-learning,future', 'tr', '2024-01-10 08:00:00', 89.1230),
('xml_provider', 'article_002', 'Blockchain Teknolojisi Rehberi', 'Blockchain teknolojisinin temelleri, kripto para birimleri ve merkezi olmayan uygulamalar (DApps) hakkında detaylı bir rehber...', 'https://example.com/articles/blockchain-guide', 'text', 15670, 789, 0, 'blockchain,cryptocurrency,dapps,decentralized', 'tr', '2024-01-18 12:30:00', 81.5670),
('xml_provider', 'article_003', 'Cloud Computing Avantajları', 'Bulut bilişimin işletmelere sağladığı avantajlar. AWS, Azure, Google Cloud karşılaştırması ve migration stratejileri...', 'https://example.com/articles/cloud-computing', 'text', 19890, 987, 0, 'cloud,aws,azure,google-cloud,migration', 'tr', '2024-01-28 15:45:00', 85.7890),
('xml_provider', 'article_004', 'DevOps Kültürü ve Araçları', 'DevOps kültürünün benimsenmesi ve kullanılan araçlar. CI/CD pipeline\'ları, monitoring ve automation...', 'https://example.com/articles/devops-culture', 'text', 14560, 654, 0, 'devops,ci-cd,automation,monitoring,culture', 'tr', '2024-02-03 10:20:00', 77.8900),
('xml_provider', 'article_005', 'Güvenlik Açıkları ve Korunma Yöntemleri', 'Web uygulamalarında karşılaşılan güvenlik açıkları ve bunlardan korunma yöntemleri. OWASP Top 10 ve best practices...', 'https://example.com/articles/security-vulnerabilities', 'text', 22340, 1123, 0, 'security,owasp,vulnerabilities,web-security,protection', 'tr', '2024-02-08 13:15:00', 91.2340);

-- =====================================================
-- YARARLI SORGULAR
-- =====================================================

-- 1. En popüler içerikler (final_score'a göre)
SELECT 
    id,
    title,
    type,
    views,
    likes,
    final_score,
    published_at
FROM contents 
ORDER BY final_score DESC 
LIMIT 10;

-- 2. İçerik türüne göre istatistikler
SELECT 
    type,
    COUNT(*) as total_content,
    AVG(views) as avg_views,
    AVG(likes) as avg_likes,
    AVG(final_score) as avg_score,
    MAX(final_score) as max_score
FROM contents 
GROUP BY type;

-- 3. Sağlayıcıya göre içerik dağılımı
SELECT 
    provider,
    COUNT(*) as content_count,
    AVG(final_score) as avg_score
FROM contents 
GROUP BY provider;

-- 4. Son 30 günde yayınlanan içerikler
SELECT 
    title,
    type,
    views,
    likes,
    final_score,
    published_at
FROM contents 
WHERE published_at >= DATE_SUB(NOW(), INTERVAL 30 DAY)
ORDER BY published_at DESC;

-- 5. En çok görüntülenen içerikler
SELECT 
    title,
    type,
    views,
    likes,
    final_score
FROM contents 
ORDER BY views DESC 
LIMIT 10;

-- 6. Etiketlere göre içerik analizi
SELECT 
    SUBSTRING_INDEX(SUBSTRING_INDEX(tags, ',', numbers.n), ',', -1) as tag,
    COUNT(*) as content_count,
    AVG(final_score) as avg_score
FROM contents 
CROSS JOIN (
    SELECT 1 n UNION ALL SELECT 2 UNION ALL SELECT 3 UNION ALL SELECT 4 UNION ALL SELECT 5
) numbers
WHERE CHAR_LENGTH(tags) - CHAR_LENGTH(REPLACE(tags, ',', '')) >= numbers.n - 1
GROUP BY tag
ORDER BY content_count DESC;

-- 7. Aylık içerik yayınlama istatistikleri
SELECT 
    DATE_FORMAT(published_at, '%Y-%m') as month,
    COUNT(*) as content_count,
    AVG(final_score) as avg_score
FROM contents 
WHERE published_at IS NOT NULL
GROUP BY month
ORDER BY month DESC;

-- 8. Skor aralıklarına göre içerik dağılımı
SELECT 
    CASE 
        WHEN final_score >= 90 THEN '90-100 (Mükemmel)'
        WHEN final_score >= 80 THEN '80-89 (Çok İyi)'
        WHEN final_score >= 70 THEN '70-79 (İyi)'
        WHEN final_score >= 60 THEN '60-69 (Orta)'
        ELSE '0-59 (Düşük)'
    END as score_range,
    COUNT(*) as content_count
FROM contents 
GROUP BY score_range
ORDER BY MIN(final_score) DESC;

-- 9. Dil bazında içerik analizi
SELECT 
    language,
    COUNT(*) as content_count,
    AVG(final_score) as avg_score,
    AVG(views) as avg_views
FROM contents 
GROUP BY language;

-- 10. Sağlayıcı ve tür kombinasyonuna göre analiz
SELECT 
    provider,
    type,
    COUNT(*) as content_count,
    AVG(final_score) as avg_score,
    SUM(views) as total_views,
    SUM(likes) as total_likes
FROM contents 
GROUP BY provider, type
ORDER BY provider, type;

-- =====================================================
-- PERFORMANS İYİLEŞTİRME SORGULARI
-- =====================================================

-- İndeks kullanımını kontrol et
EXPLAIN SELECT * FROM contents WHERE type = 'video' AND final_score > 80;

-- Yavaş sorguları bul
SELECT 
    table_schema,
    table_name,
    index_name,
    cardinality
FROM information_schema.statistics 
WHERE table_schema = 'search_engine_service' 
AND table_name = 'contents';

-- =====================================================
-- BAKIM SORGULARI
-- =====================================================

-- Eski içerikleri temizle (6 aydan eski)
-- DELETE FROM contents WHERE published_at < DATE_SUB(NOW(), INTERVAL 6 MONTH);

-- Duplicate içerikleri bul
SELECT 
    provider, 
    provider_id, 
    COUNT(*) as duplicate_count
FROM contents 
GROUP BY provider, provider_id 
HAVING COUNT(*) > 1;

-- Boş veya eksik verileri bul
SELECT 
    id,
    title,
    provider,
    provider_id
FROM contents 
WHERE title IS NULL 
   OR title = '' 
   OR url IS NULL 
   OR url = '';

-- =====================================================
-- TEST VERİLERİ
-- =====================================================

-- Test için ek içerikler
INSERT INTO contents (provider, provider_id, title, description, url, type, views, likes, duration, tags, language, published_at, final_score) VALUES
('test_provider', 'test_001', 'Test Video İçeriği', 'Bu bir test içeriğidir.', 'https://test.com/video1', 'video', 100, 10, 120, 'test,video', 'tr', NOW(), 50.0000),
('test_provider', 'test_002', 'Test Metin İçeriği', 'Bu bir test metin içeriğidir.', 'https://test.com/article1', 'text', 50, 5, 0, 'test,text', 'tr', NOW(), 45.0000);

-- Test verilerini temizle
-- DELETE FROM contents WHERE provider = 'test_provider';

-- =====================================================
-- VERİTABANI İSTATİSTİKLERİ
-- =====================================================

-- Tablo boyutu
SELECT 
    table_name,
    ROUND(((data_length + index_length) / 1024 / 1024), 2) AS 'Size (MB)',
    table_rows
FROM information_schema.tables 
WHERE table_schema = 'search_engine_service' 
AND table_name = 'contents';

-- Toplam kayıt sayısı
SELECT COUNT(*) as total_contents FROM contents;

-- Son güncelleme zamanı
SELECT MAX(updated_at) as last_update FROM contents; 