package http

import (
	"path/filepath"

	"github.com/gin-contrib/multitemplate"
)

// loadTemplates, Gin için özel bir HTML renderer oluşturur.
// Bu yapı, sayfaların birbirinin "content" bloklarını ezmesini engeller.
// Her sayfa, kendi "base" şablonuyla birleştirilerek izole edilir.
func NewRenderer() multitemplate.Renderer {
	r := multitemplate.NewRenderer()

	// Şablon Dizini
	templatesDir := "web/templates"
	layoutDir := filepath.Join(templatesDir, "layout")
	
	// 1. Layout Dosyaları
	baseLayout := filepath.Join(layoutDir, "base.html")

	// 2. Sayfa Tanımları (Manual Mapping - Mükemmel Kontrol İçin)
	
	// Dashboard
	r.AddFromFiles("dashboard/dashboard.html", baseLayout, filepath.Join(templatesDir, "dashboard", "dashboard.html"))

	// Auth (Standalone - Base KULLANMAZ)
	r.AddFromFiles("auth/login.html", filepath.Join(templatesDir, "auth", "login.html"))

	// Clients
	r.AddFromFiles("clients/list.html", baseLayout, filepath.Join(templatesDir, "clients", "list.html"))
	r.AddFromFiles("clients/create.html", baseLayout, filepath.Join(templatesDir, "clients", "create.html"))
	r.AddFromFiles("clients/edit.html", baseLayout, filepath.Join(templatesDir, "clients", "edit.html"))

	// Dava Dosyaları (Cases)
	r.AddFromFiles("cases/list.html", baseLayout, filepath.Join(templatesDir, "cases", "list.html"))
	r.AddFromFiles("cases/create.html", baseLayout, filepath.Join(templatesDir, "cases", "create.html"))
	r.AddFromFiles("cases/edit.html", baseLayout, filepath.Join(templatesDir, "cases", "edit.html"))

	// Duruşmalar (Hearings)
	r.AddFromFiles("hearings/list.html", baseLayout, filepath.Join(templatesDir, "hearings", "list.html"))
	r.AddFromFiles("hearings/create.html", baseLayout, filepath.Join(templatesDir, "hearings", "create.html"))
	r.AddFromFiles("hearings/edit.html", baseLayout, filepath.Join(templatesDir, "hearings", "edit.html"))

	// Yeni sayfalar buraya eklenecek...

	return r
}
