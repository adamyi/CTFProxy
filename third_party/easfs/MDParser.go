package main

import (
	"bytes"
	"fmt"
	"net/http"
	"regexp"

	"github.com/flosch/pongo2"
	"github.com/gholt/blackfridaytext"
	"gopkg.in/russross/blackfriday.v2"
)

func ParseMD(w http.ResponseWriter, content []byte, requestPath string) error {
	tmpl, err := pongo2.FromFile("/templates/page-article.html")
	if err != nil {
		fmt.Println(err.Error())
		return err
	}
	context := pongo2.Context{}
	context["bodyClass"] = "devsite-doc-page"
	context["requestPath"] = requestPath
	context["isProd"] = flagProd
	meta, offset := blackfridaytext.MarkdownMetadata(content)
	content = content[offset:]

	for _, val := range meta {
		context[val[0]] = val[1]
	}

	//TODO: include
	content = RenderContent(content)

	// add datetime
	publishDate := regexp.MustCompile(`{# published_on: (\d{4}-\d{2}-\d{2}) #}`).FindSubmatch(content)
	if len(publishDate) > 0 && !bytes.Equal(publishDate[1], []byte("1900-01-01")) {
		context["publishDate"] = string(publishDate[1])
	}
	updateDate := regexp.MustCompile(`{# updated_on: (\d{4}-\d{2}-\d{2}) #}`).FindSubmatch(content)
	if len(updateDate) > 0 && !bytes.Equal(updateDate[1], []byte("1900-01-01")) {
		context["updateDate"] = string(updateDate[1])
	}

	// remove comment
	content = regexp.MustCompile(`{#.+?#}`).ReplaceAll(content, []byte(""))
	content = regexp.MustCompile(`(?ms){% comment %}.*?{% endcomment %}`).ReplaceAll(content, []byte(""))

	// render callouts
	content = regexp.MustCompile(`(?ms)^Note: (.*?)\n^\n`).ReplaceAll(content, []byte("<aside class=\"note\" markdown=\"1\"><strong>Note:</strong> <span>$1</span></aside>\n"))
	content = regexp.MustCompile(`(?ms)^Caution: (.*?)\n^\n`).ReplaceAll(content, []byte("<aside class=\"caution\" markdown=\"1\"><strong>Caution:</strong> <span>$1</span></aside>\n"))
	content = regexp.MustCompile(`(?ms)^Warning: (.*?)\n^\n`).ReplaceAll(content, []byte("<aside class=\"warning\" markdown=\"1\"><strong>Warning:</strong> <span>$1</span></aside>\n"))
	content = regexp.MustCompile(`(?ms)^Key Point: (.*?)\n^\n`).ReplaceAll(content, []byte("<aside class=\"key-point\" markdown=\"1\"><strong>Key Point:</strong> <span>$1</span></aside>\n"))
	content = regexp.MustCompile(`(?ms)^Key Term: (.*?)\n^\n`).ReplaceAll(content, []byte("<aside class=\"key-term\" markdown=\"1\"><strong>Key Term:</strong> <span>$1</span></aside>\n"))
	content = regexp.MustCompile(`(?ms)^Objective: (.*?)\n^\n`).ReplaceAll(content, []byte("<aside class=\"objective\" markdown=\"1\"><strong>Objective:</strong> <span>$1</span></aside>\n"))
	content = regexp.MustCompile(`(?ms)^Success: (.*?)\n^\n`).ReplaceAll(content, []byte("<aside class=\"success\" markdown=\"1\"><strong>Success:</strong> <span>$1</span></aside>\n"))
	content = regexp.MustCompile(`(?ms)^Dogfood: (.*?)\n^\n`).ReplaceAll(content, []byte("<aside class=\"dogfood\" markdown=\"1\"><strong>Dogfood:</strong> <span>$1</span></aside>\n"))
	// fmt.Println(context)
	// fmt.Println(string(content))
	var extensions blackfriday.Extensions
	extensions |= blackfriday.Tables
	extensions |= blackfriday.FencedCode
	extensions |= blackfriday.Strikethrough
	extensions |= blackfriday.LaxHTMLBlocks
	// extensions |= blackfriday.AutoHeadingIDs
	var hp blackfriday.HTMLRendererParameters
	// hp.Flags |= blackfriday.TOC
	renderer := blackfriday.NewHTMLRenderer(hp)
	markdown := blackfriday.New(blackfriday.WithExtensions(extensions))
	ast := markdown.Parse(content)
	// content = blackfriday.Run(content, blackfriday.WithExtensions(extensions), blackfriday.WithRenderer(renderer))
	// context["content"] = string(content)
	var buf bytes.Buffer
	// renderer.RenderHeader(&buf, ast)
	// toc := buf.Bytes()
	// toc = regexp.MustCompile(`(?s)<nav>(.*?<ul>){2}`).ReplaceAll(toc, []byte(""))
	// toc = regexp.MustCompile(`(?s)</ul>\s*</li>\s*</ul>\s*</nav>`).ReplaceAll(toc, []byte(""))
	// toc = regexp.MustCompile(`<ul>`).ReplaceAll(toc, []byte("<ul class=\"devsite-page-nav-list\">"))
	// toc = regexp.MustCompile(`<a href`).ReplaceAll(toc, []byte("<a class=\"devsite-nav-title\" href"))
	// toc = regexp.MustCompile(`<li>`).ReplaceAll(toc, []byte("<li class=\"devsite-page-nav-item\">"))
	// context["renderedTOC"] = string(toc)
	// buf.Reset()
	ast.Walk(func(node *blackfriday.Node, entering bool) blackfriday.WalkStatus {
		return renderer.RenderNode(&buf, node, entering)
	})
	renderer.RenderFooter(&buf, ast)
	content = buf.Bytes()
	context["content"] = string(content)
	// fmt.Println(string(content))
	project, parentProject, err := ParseProject(context["project_path"].(string))
	if err != nil {
		return err
	}
	book, err := ParseBook(context["book_path"].(string))
	if err != nil {
		return err
	}
	context["projectYaml"] = *project
	context["bookYaml"] = book
	context["lowerTabs"] = GetLowerTabs(requestPath, book)
	context["footerBanner"], context["footerPromos"], context["footerLinks"], err = ParseFooter(project.FooterPath)
	if err != nil {
		return err
	}
	context["renderedLeftNav"] = GetLeftNav(requestPath, book)

	context["logoRowIcon"] = project.Icon.Path
	context["logoRowIconType"] = project.Icon.Type
	if parentProject != nil {
		context["logoRowTitle"] = parentProject.Name
	} else {
		context["logoRowTitle"] = project.Name
	}
	context["headerTitle"] = project.Name
	// context["headerDescription"] = project.Description
	pageTitle := project.Name
	if parentProject != nil {
		pageTitle += " | " + parentProject.Name
	}
	// titleRO := regexp.MustCompile(`<h1 class="page-title".*?>(.*?)<\/h1>`).FindSubmatch(content)
	titleRO := regexp.MustCompile(`<h1>(.*?)<\/h1>`).FindSubmatch(content)
	if len(titleRO) > 0 {
		pageTitle = string(titleRO[1]) + " | " + pageTitle
	}
	context["pageTitle"] = pageTitle
	return tmpl.ExecuteWriter(context, w)
}
