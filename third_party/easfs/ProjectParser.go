package main

import (
	"io/ioutil"

	"gopkg.in/yaml.v2"
)

type Project struct {
	ParentProjectMetadataPath string   `yaml:"parent_project_metadata_path"`
	Name                      string   `yaml:"name"`
	Description               string   `yaml:"description"`
	HomeURL                   string   `yaml:"home_url"`
	Color                     string   `yaml:"color"`
	ContentLicense            string   `yaml:"content_license"`
	FooterPath                string   `yaml:"footer_path"`
	GoogleAnalyticsIds        []string `yaml:"google_analytics_ids"`
	Icon                      struct {
		Path string `yaml:"path"`
		Type string `yaml:"type"`
	} `yaml:"icon"`
	SocialMedia struct {
		Image struct {
			Path   string `yaml:"path"`
			Width  int    `yaml:"width"`
			Height int    `yaml:"height"`
		} `yaml:"image"`
	} `yaml:"social_media"`
}

func ParseProject(filepath string) (*Project, *Project, error) {
	projectContent, err := ioutil.ReadFile(flagSitePath + filepath)
	if err != nil {
		return nil, nil, err
	}
	project := Project{}
	err = yaml.Unmarshal(projectContent, &project)
	if err != nil {
		return nil, nil, err
	}
	if project.ParentProjectMetadataPath == "" {
		return &project, nil, nil
	}
	parentProjectContent, err := ioutil.ReadFile(flagSitePath + project.ParentProjectMetadataPath)
	if err != nil {
		return nil, nil, err
	}
	parentProject := Project{}
	err = yaml.Unmarshal(parentProjectContent, &parentProject)
	if err != nil {
		return nil, nil, err
	}
	return &project, &parentProject, nil
}
