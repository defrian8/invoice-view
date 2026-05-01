package service

import (
	"bytes"
	"context"
	"fmt"
	"invoice-view/model"
	"invoice-view/util"
	"log"
	"path/filepath"
	"strings"
	"text/template"
)

type TemplateService interface {
	Render(ctx context.Context, invoice *model.Invoice) (bytes.Buffer, error)
}

type templateService struct {
	funcMap template.FuncMap
	tmplMap map[string]*template.Template
}

func NewTemplateService() TemplateService {
	funcMap := template.FuncMap{
		"formatCurrency": util.FormatCurrency,
		"dateFormat":     util.DateFormat,
		"statusClass":    util.StatusClassCSS,
	}

	t := &templateService{funcMap: funcMap, tmplMap: make(map[string]*template.Template)}
	t.loadTemplates("template/*.html")

	return t
}

func (s *templateService) loadTemplates(pattern string) {
	matches, err := filepath.Glob(pattern)
	if err != nil {
		panic(fmt.Sprintf("failed to glob templates: %v", err))
	}

	if len(matches) == 0 {
		panic("no template files found in " + pattern)
	}

	for _, file := range matches {
		base := filepath.Base(file)
		name := strings.TrimSuffix(base, ".html")

		tmpl, err := template.New(base).Funcs(s.funcMap).ParseFiles(file)
		if err != nil {
			panic(fmt.Sprintf("failed to parse template %s: %v", file, err))
		}

		s.tmplMap[name] = tmpl
		log.Printf("[router] loaded template: %s -> %s", name, file)
	}
}

func (s *templateService) chooseTemplate(config model.Config) string {
	defaultTemplate := "invoice"
	name := strings.TrimSpace(config.Template)
	if name == "" {
		return defaultTemplate
	}

	if _, ok := s.tmplMap[name]; !ok {
		return defaultTemplate
	}

	return name
}

func (s *templateService) Render(ctx context.Context, invoice *model.Invoice) (bytes.Buffer, error) {
	templateName := s.chooseTemplate(invoice.Config)
	tmpl := s.tmplMap[templateName]
	defineName := templateName + ".html"

	var buf bytes.Buffer
	if err := tmpl.ExecuteTemplate(&buf, defineName, invoice); err != nil {
		log.Printf("[router] failed to render %s: %v", defineName, err)
		return buf, err
	}

	return buf, nil
}
