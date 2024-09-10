package templates

import (
	"html/template"
	"log"
	"net/http"
)

// Define a global template map
var Templates = map[string]*template.Template{}

var BaseLocation = ""

// Initialize templates
func InitTemplates(location string) {
	log.Print("Initializing")
	BaseLocation = location
	// Load and parse base template
	baseTemplate, err := template.ParseFiles(BaseLocation + "/base.html")
	if err != nil {
		log.Fatalf("Failed to parse base template: %v", err)
	}
	Templates["base"] = baseTemplate

	// Parse page-specific templates
	pageTemplates := []string{
		BaseLocation + "/index.html",
		BaseLocation + "/pricing.html",
		BaseLocation + "/about.html",
		BaseLocation + "/alerts.html",
		BaseLocation + "/profile.html",
		BaseLocation + "/reset-password-sent.html",
		BaseLocation + "/reset-password-success.html",
		BaseLocation + "/subscription-success.html",
		BaseLocation + "/subscription-cancel.html",
		BaseLocation + "/token-expired.html",
		BaseLocation + "/docs.html",
		BaseLocation + "/404.html",
		BaseLocation + "/error.html",
	}

	for _, file := range pageTemplates {
		tmpl, err := template.Must(baseTemplate.Clone()).ParseFiles(file)
		if err != nil {
			log.Fatalf("Failed to parse page template %s: %v", file, err)
		}
		Templates[file] = tmpl
	}
}

func RenderTemplate(w http.ResponseWriter, templateName string, data map[string]interface{}) {
	tmpl, ok := Templates[templateName]
	if !ok {
		http.Error(w, "Template not found", http.StatusNotFound)
		return
	}

	err := tmpl.ExecuteTemplate(w, "base.html", data)
	if err != nil {
		log.Printf("Failed to execute template: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
