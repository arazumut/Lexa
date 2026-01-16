# ⚔️ LEXA: PROJECT STATUS REPORT
**Tarih:** 16.01.2026 (Refactored)
**Durum:** Kritk Güvenlik Açığı Giderildi, Document Storage (Evrak) Modülü Eklendi.

---

## 🏗️ 1. MİMARİ VE TEKNOLOJİ YIĞINI (TECH STACK)
Proje, **Clean Architecture (Temiz Mimari)** prensiplerine sadık kalınarak geliştirilmiştir.

### 🔧 Backend - Güncellemeler
*   **Security:** Hardcoded JWT secret kaldırıldı. `.env` üzerinden `JWT_SECRET` okunuyor.
*   **File Storage:** `google/uuid` tabanlı dosya isimlendirme ve `web/static/uploads` yerel depolama sistemi kuruldu.
*   **Modüller:**
    *   Auth (Tamam)
    *   Client (Tamam)
    *   Case (Tamam)
    *   Hearing (Tamam) - *Dashboard'a entegre.*
    *   Accounting/Transaction (Tamam) - *Dashboard'da grafikler aktif.*
    *   Document (YENİ) - *Evrak yükleme ve listeleme altyapısı hazır.*

### 📂 Klasör Yapısı
```text
LEXA/
├── cmd/app/main.go            # Dependency Injection ve Config burada yönetiliyor.
├── config/                    # Env ve Config yönetimi.
├── internal/
│   ├── domain/                # Saf Go structları (User, Client, Case, Document...).
│   ├── repository/            # GORM implementasyonları.
│   ├── service/               # İş mantığı (Upload, Calc Balance vb.).
│   └── transport/http/        # Gin Handler'lar.
├── web/
│   └── static/uploads/        # Yüklenen evraklar burada tutulur.
└── .env                       # Hassas bilgiler (Git-ignored).
```

---

## ✅ 2. TAMAMLANAN KRİTİK GELİŞTİRMELER (DONE)

### 🔴 ACİL GÜVENLİK DÜZELTMESİ
*   [x] `main.go` içindeki hardcoded anahtar temizlendi.
*   [x] `Config` paketi `.env` desteği ile güncellendi.
*   [x] 256-bit secure hex key oluşturulup `.env` dosyasına yazıldı.

### 📄 FAZ 4: Evrak Yönetimi (Document Management)
*   [x] **Domain:** `Document` entity oluşturuldu (Dosya Adı, Tipi, Yolu, Yükleyen).
*   [x] **Repository:** Dosyaları davaya göre (`FindByCaseID`) getiren repo yazıldı.
*   [x] **Service:**
    *   `multipart/form-data` işleme mantığı.
    *   UUID ile benzersiz dosya adı oluşturma (`uuid.v4`).
    *   Fiziksel diskten ve DB'den silme (`os.Remove`).
*   [x] **API:** `/api/documents/upload` ve `/api/cases/:id/documents` uçları hazır.

---

## 🚀 3. SIRADAKİ ADIMLAR (TODO)
Kod şu an backend tarafında **%95 tamamlandı**. Sadece UI eksikleri kaldı.

1.  **UI Entegrasyonu (Document):**
    *   Dava detay sayfasına (`cases/detail.html` - *henüz yok*) veya edit sayfasına "Dosyalar" sekmesi eklenecek.
    *   AJAX ile dosya yükleme scripti yazılacak.
2.  **Test Yazımı:**
    *   Hiç test yok. Kritik servisler için unit test yazılmalı.
3.  **Deploy Hazırlığı:**
    *   Dockerfile `uploads` klasörü permission ayarları kontrol edilecek (Render'da volume gerekebilir).
