package templates

import (
	"html/template"
	"log"
	"net/http"
)

// Define a global template map
var templates = map[string]*template.Template{}

// Initialize templates
func InitTemplates() {
	log.Printf("Initializing")

	// Load and parse base template
	baseTemplate, err := template.ParseFiles("./templates/base.html")
	if err != nil {
		log.Fatalf("Failed to parse base template: %v", err)
	}
	templates["base"] = baseTemplate

	// Parse page-specific templates
	pageTemplates := []string{
		"./templates/index.html",
		"./templates/pricing.html",
		"./templates/about.html",
		"./templates/404.html",
	}

	for _, file := range pageTemplates {
		tmpl, err := template.Must(baseTemplate.Clone()).ParseFiles(file)
		if err != nil {
			log.Fatalf("Failed to parse page template %s: %v", file, err)
		}
		templates[file] = tmpl
	}
	log.Printf("templates:", templates)
}

// Render the specified page template within the base layout
func RenderTemplate(w http.ResponseWriter, templateName string, title string) {
	tmpl, ok := templates[templateName]
	if !ok {
		http.Error(w, "Template not found", http.StatusNotFound)
		return
	}

	err := tmpl.ExecuteTemplate(w, "base.html", map[string]interface{}{
		"Title":   title,
		"Content": templateName,
	})
	if err != nil {
		log.Printf("Failed to execute template: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
