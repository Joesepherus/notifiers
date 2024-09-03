package templates

import (
	"html/template"
	"log"
	"net/http"
)

// Define a global template map
var Templates = map[string]*template.Template{}

// Initialize templates
func InitTemplates() {
	log.Printf("Initializing")

	// Load and parse base template
	baseTemplate, err := template.ParseFiles("./templates/base.html")
	if err != nil {
		log.Fatalf("Failed to parse base template: %v", err)
	}
	Templates["base"] = baseTemplate

	// Parse page-specific templates
	pageTemplates := []string{
		"./templates/index.html",
		"./templates/pricing.html",
		"./templates/about.html",
		"./templates/alerts.html",
		"./templates/profile.html",
		"./templates/reset-password-sent.html",
		"./templates/reset-password-success.html",
		"./templates/subscription-success.html",
		"./templates/subscription-cancel.html",
		"./templates/docs.html",
		"./templates/404.html",
	}

	for _, file := range pageTemplates {
		tmpl, err := template.Must(baseTemplate.Clone()).ParseFiles(file)
		if err != nil {
			log.Fatalf("Failed to parse page template %s: %v", file, err)
		}
		Templates[file] = tmpl
	}
	log.Printf("templates:", Templates)
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
