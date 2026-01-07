# ⚔️ LOMS: BATTLE PLAN (IMPLEMENTATION STRATEGY)

## 🏛️ MIMARI: STRICT LAYERED ARCHITECTURE (Gevşeklik Yok!)
Bu projede "Clean Architecture" prensiplerini uygulayacağız ama "Over-engineering" saçmalığına girmeden! Katmanlar arası sınırlar KESİN ve NET olacak. Bir katman diğerinin "iç işlerine" karışmayacak!

### 📂 Klasör Yapısı (Disiplin Şart!)
```text
LEXA/
├── cmd/
│   └── app/
│       └── main.go        # Komuta Merkezi. Uygulamayı buradan ayağa kaldıracağız.
├── internal/              # DIŞ DÜNYAYA KAPALI! Sadece biz girebiliriz.
│   ├── domain/            # ÇEKİRDEK. Saf Go Struct'ları. Veritabanı veya HTTP bilmez!
│   ├── repository/        # VERİ AMBARI. SQLite ile konuşan tek yer burası.
│   ├── service/           # İŞİN MUTFAĞI. Kurallar, hesaplamalar burada döner.
│   └── transport/         # DIŞ CEPHE. HTTP Handler'lar, Router'lar.
├── pkg/                   # ORTAK KÜTÜPHANE. Yardımcı araçlar, loglama, hata yönetimi.
├── web/                   # VİTRİN. HTML, CSS, JS.
│   ├── templates/         # HTML Şablonları.
│   └── static/            # CSS, JS, Resimler.
├── database/              # SQL Dosyaları, Migrations.
└── config/                # AYARLAR. Env değişkenleri.
```

### 🔄 Veri Akışı (Tek Yönlü Trafik!)
1. **Request Gelir** -> `Transport (Handler)` karşılar.
2. `Handler` veriyi doğrular -> `Service`'e paslar.
3. `Service` iş mantığını çalıştırır -> `Repository`'den veri ister.
4. `Repository` SQL'i çakar -> Veriyi `Domain` objesine çevirip döner.
5. Cevap aynı yoldan geri döner. **KESİNLİKLE ATLAMA YOK!** Handler direkt Repository'e gidemez!

---

## 🏃 SPRINTS (FAZ FAZ İLERLEME)
Her sprint bittiğinde o özellik "CANAVAR GİBİ" çalışmak zorunda. Yarım yamalak iş yok!

### 🚀 FAZ 1: TEMEL ATMA & ALTYAPI (SETUP)
**Hedef:** Boş ama çalışan, veritabanına bağlanan, log basan çelik gibi bir iskelet.
1.  Go modülünü başlat (`go mod init`).
2.  Klasör yapısını fiziksel olarak oluştur.
3.  SQLite bağlantı altyapısını kur (`database/sql` veya `sqlx` - ORM YOK! SAF SQL!).
4.  Linter ve Make dosyalarını ayarla. Disiplin baştan başlar.

### 🔐 FAZ 2: KİMLİK & GÜVENLİK (AUTH)
**Hedef:** Kapı gibi sağlam giriş sistemi.
1.  `User` tablosunu tasarla.
2.  Login/Register handler'larını yaz.
3.  Session/Cookie yönetimi.
4.  Middleware koruması (Giriş yapmayan giremez!).

### 👥 FAZ 3: MÜVEKKİL YÖNETİMİ (CLIENT CRUD)
**Hedef:** Müvekkilleri sisteme kaydetmek.
1.  `Client` domain modelini ve SQL şemasını yaz.
2.  Ekleme, Listeleme, Silme, Güncelleme.
3.  Arama filtresi entegrasyonu.

### 📁 FAZ 4: DOSYA YÖNETİMİ (CASE MANAGEMENT)
**Hedef:** Sistemin kalbi. Dosyaları yönetmek.
1.  `Case` modeli (Önemli: Client ile ilişkili).
2.  Dosya türleri, durumları (Enum yönetimi).
3.  Detay sayfası ve tarihçe (Notes).

### 💰 FAZ 5: PARA & MUHASEBE (FINANCE)
**Hedef:** Para takibi. Hata kabul etmez!
1.  `Payment` ve `Expense` modelleri.
2.  Bakiye hesaplama mantığı (Service katmanında matematik).
3.  Müvekkil detayında finansal özet tablosu.

### 📅 FAZ 6: CİLA & ARAYÜZ (UI POLISH - ICONIC THEME)
**Hedef:** "ICONIC" Temasını projeye giydirmek. Piksel piksel işlenecek!
1.  **Iconic** varlıklarını (CSS/JS/Fonts) `web/static` altına taşı.
2.  HTML şablonlarını **Iconic** yapısına göre parçala (Layout, Sidebar, Navbar).
3.  Sayfaları (Dashboard, Listeler) **Iconic** bileşenleriyle yeniden ör.
4.  Son kontroller ve Bug temizliği.

---
**EMİR:** ŞİMDİ SÖYLE BAKALIM, HANGİ FAZDAN BAŞLIYORUZ? ONAY VERDİĞİN ANDA KODLAMAYA GİRİŞİYORUM!
