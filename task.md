# 📋 PROJECT CANVAS: LEGAL OFFICE MANAGEMENT SYSTEM (LOMS)
**Type:** Pure CRUD / CRM / ERP
**No AI - No UYAP - No Bullshit**

---

## 🏗️ 1. USER & AUTHENTICATION (GİRİŞ & YETKİ)
*Sisteme giriş kapısı. Avukat ve gerekirse sekreteri.*

- [ ] **Login / Register**
    - Email & Şifre ile giriş.
    - Basit "Beni Hatırla" yapısı.
- [ ] **Profile Management**
    - Avukat bilgileri (Ad, Soyad, Baro Sicil, İletişim).
    - Büro Logosu yükleme (Faturalar/Raporlar için).
- [ ] **Role Management (Opsiyonel)**
    - `Admin` (Avukat): Her şeyi görür.
    - `Staff` (Sekreter/Stajyer): Sadece dosya ekler, muhasebeyi görmez.

---

## 👥 2. CLIENT MANAGEMENT (MÜVEKKİL CRUD)
*İşin kökü müvekkil. Kimin işini yapıyoruz?*

- [ ] **Create Client (Müvekkil Ekle)**
    - Tip Seçimi: Gerçek Kişi (TCKN) / Tüzel Kişi (Vergi No).
    - İletişim: Telefon, Email, Adres.
    - Notlar: Müvekkil hakkında özel not alanı.
- [ ] **Client List & Search**
    - İsimden, TCKN'den anlık filtreleme.
    - Bakiye Gösterimi (Bu adamın bize borcu var mı?).
- [ ] **Client Detail View**
    - Tek ekranda müvekkilin **tüm dosyaları**, **tüm ödemeleri**, **tüm evrakları**.

---

## 📁 3. CASE MANAGEMENT (DOSYA YÖNETİMİ)
*Müvekkile bağlı dava/icra dosyaları.*

- [ ] **Create Case (Dosya Aç)**
    - Müvekkil Seçimi (Dropdown).
    - Dosya Türü: Dava, İcra, Danışmanlık.
    - Mahkeme / İcra Dairesi Bilgisi.
    - Esas No / Dosya No (Manuel giriş).
    - Karşı Taraf Bilgileri (Davalı/Borçlu kim?).
- [ ] **Case Stages (Aşamalar)**
    - Durum Güncelleme: Dava Açıldı -> Ön İnceleme -> Bilirkişi -> Karar -> İstinaf.
- [ ] **Case Notes (Tarihçe)**
    - Dosyaya tarihli not düşme (Örn: "Bugün kalemle görüşüldü, müzekkere yazılmış").

---

## 💰 4. FINANCE & ACCOUNTING (MUHASEBE)
*Para takibi. Avukatın en hassas olduğu yer.*

- [ ] **Add Payment (Tahsilat Gir)**
    - Hangi Müvekkil? Hangi Dosya?
    - Tutar, Tarih, Açıklama.
    - Tür: Nakit, Havale, Kredi Kartı.
- [ ] **Expense Tracking (Masraf Gir)**
    - Dosya için yapılan masraflar (Harç, Yol, Posta).
    - Masrafı kim ödedi? (Bürodan mı çıktı, müvekkil mi verdi?).
- [ ] **Balance Calculation (Hesap Özeti)**
    - Anlaşılan Ücret - (Ödenenler) = **Kalan Bakiye**.
    - Dosya bazlı kâr/zarar durumu.

---

## 📅 5. AGENDA & TASKS (AJANDA)
*Duruşma ve iş takibi.*

- [ ] **Event Creation**
    - Duruşma Tarihi Ekle.
    - Süreli İş Ekle (Örn: "Cevap dilekçesi son gün").
- [ ] **Calendar View**
    - Aylık/Haftalık görünüm.
    - Yaklaşan işler listesi.

---

## 📂 6. DOCUMENT STORAGE (EVRAK)
*Dosyaları klasörlemek için. AI yok, sadece depolama.*

- [ ] **File Upload**
    - Dosyanın içine PDF/Resim yükleme.
    - "Dava Dilekçesi", "Bilirkişi Raporu" diye etiketleme.
- [ ] **Download/View**
    - Yüklenen evrakı indirme.

---

## ⚙️ 7. SYSTEM SPECS (TEKNİK)
- **Backend:** Go (CRUD işlemleri için en hızlısı).
- **Database:** SQLite (Kurulum gerektirmez, tek dosya, hızlı ve taşınabilir).
- **Storage:** Local Disk veya MinIO (Dosyalar için).
- **Frontend:** Server Side Rendering (HTML) - **Tema: ICONIC** (Referans tasarım birebir uygulanacak).