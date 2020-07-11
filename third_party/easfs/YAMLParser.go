package main

import (
	"fmt"
	"net/http"

	"github.com/flosch/pongo2"
	"gopkg.in/yaml.v2"
)

type YAMLPage struct {
	ProjectPath string `yaml:"project_path"`
	BookPath    string `yaml:"book_path"`
	Title       string `yaml:"title"`
	LandingPage struct {
		CustomCSSPath string `yaml:"custom_css_path,omitempty"`
		CustomJSPath  string `yaml:"custom_js_path,omitempty"`
		MetaTags      []struct {
			Name    string `yaml:"name"`
			Content string `yaml:"content"`
		} `yaml:"meta_tags"`
		Header struct {
			Name        string `yaml:"name,omitempty"`
			Description string `yaml:"description,omitempty"`
			CustomHTML  string `yaml:"custom_html,omitempty"`
		} `yaml:"header"`
		Rows []struct {
			ClassName string `yaml:"classname,omitempty"`
			Items     []struct {
				ClassName   string `yaml:"classname,omitempty"`
				Heading     string `yaml:"heading,omitempty"`
				Description string `yaml:"description,omitempty"`
				ImagePath   string `yaml:"image_path,omitempty"`
				Path        string `yaml:"path,omitempty"`
				CustomHTML  string `yaml:"custom_html,omitempty"`
				Buttons     []struct {
					Label     string `yaml:"label"`
					Target    string `yaml:"target"`
					Path      string `yaml:"path"`
					ClassName string `yaml:"classname"`
				} `yaml:"buttons,omitempty"`
			} `yaml:"items"`
			Heading     string `yaml:"heading,omitempty"`
			Background  string `yaml:"background,omitempty"`
			Description string `yaml:"description,omitempty"`
			CustomHTML  string `yaml:"custom_html,omitempty"`
			ItemCount   int    `yaml:"item_count,omitempty"`
		} `yaml:"rows"`
	} `yaml:"landing_page"`
}

func ParseYAML(w http.ResponseWriter, content []byte, requestPath string) error {
	tmpl, err := pongo2.FromFile("/templates/page-landing.html")
	if err != nil {
		fmt.Println(err.Error())
		return err
	}
	context := pongo2.Context{}
	context["bodyClass"] = "devsite-landing-page"
	context["requestPath"] = requestPath
	context["isProd"] = flagProd
	parsedYAML := YAMLPage{}
	err = yaml.Unmarshal(content, &parsedYAML)
	if err != nil {
		return err
	}
	for idx := range parsedYAML.LandingPage.Rows {
		if parsedYAML.LandingPage.Rows[idx].CustomHTML != "" {
			parsedYAML.LandingPage.Rows[idx].CustomHTML = string(RenderContent([]byte(parsedYAML.LandingPage.Rows[idx].CustomHTML)))
		}
		parsedYAML.LandingPage.Rows[idx].ItemCount = len(parsedYAML.LandingPage.Rows[idx].Items)
		for idx2 := range parsedYAML.LandingPage.Rows[idx].Items {
			if parsedYAML.LandingPage.Rows[idx].Items[idx2].CustomHTML != "" {
				parsedYAML.LandingPage.Rows[idx].Items[idx2].CustomHTML = string(RenderContent([]byte(parsedYAML.LandingPage.Rows[idx].Items[idx2].CustomHTML)))
			}
		}
	}
	// fmt.Println(parsedYAML)
	context["rows"] = parsedYAML.LandingPage.Rows
	context["customJSPath"] = parsedYAML.LandingPage.CustomJSPath
	context["customCSSPath"] = parsedYAML.LandingPage.CustomCSSPath
	context["metaTags"] = parsedYAML.LandingPage.MetaTags
	project, parentProject, err := ParseProject(parsedYAML.ProjectPath)
	if err != nil {
		return err
	}
	book, err := ParseBook(parsedYAML.BookPath)
	if err != nil {
		return err
	}
	context["projectYaml"] = *project
	context["logoRowIcon"] = project.Icon.Path
	context["logoRowIconType"] = project.Icon.Type
	if parsedYAML.LandingPage.Header.Name != "" {
		context["logoRowTitle"] = parsedYAML.LandingPage.Header.Name
	} else if parentProject != nil {
		context["logoRowTitle"] = parentProject.Name
	} else {
		context["logoRowTitle"] = project.Name
	}
	context["customHeader"] = string(RenderContent([]byte(parsedYAML.LandingPage.Header.CustomHTML)))
	if parsedYAML.Title != "" {
		context["headerTitle"] = parsedYAML.Title
	} else if project.Name != "" {
		context["headerTitle"] = project.Name
	} else if parentProject != nil {
		context["headerTitle"] = parentProject.Name
	}
	if parsedYAML.LandingPage.Header.Description != "" {
		context["headerDescription"] = parsedYAML.LandingPage.Header.Description
	} else if project.Description != "" {
		context["headerDescription"] = project.Description
	} else if parentProject != nil {
		context["headerDescription"] = parentProject.Description
	}
	// TODO: header buttons
	if parsedYAML.Title != "" {
		context["pageTitle"] = parsedYAML.Title + " | " + project.Name
	} else {
		context["pageTitle"] = project.Name
	}
	context["bookYaml"] = book
	context["lowerTabs"] = GetLowerTabs(requestPath, book)
	context["footerBanner"], context["footerPromos"], context["footerLinks"], err = ParseFooter(project.FooterPath)
	if err != nil {
		return err
	}

	// fmt.Println(context)

	return tmpl.ExecuteWriter(context, w)
}
