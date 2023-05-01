package cmd

import (
	"bytes"
	"fmt"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v2"
	"html/template"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
)

type WebPage struct {
	filename string

	PageName               string `yaml:"page_name"`
	UseLocalPath           bool   `yaml:"use_local_path"`
	LocalDirectHitLink     string `yaml:"local_direct_hit_link"`
	LocalMultipleHitsLink  string `yaml:"local_multiple_hits_link"`
	LocalNoHitsLink        string `yaml:"local_no_hits_link"`
	PublicDirectHitLink    string `yaml:"public_direct_hit_link"`
	PublicMultipleHitsLink string `yaml:"public_multiple_hits_link"`
	PublicNoHitsLink       string `yaml:"public_no_hits_link"`
}

func generateRobotScript(webpage *WebPage) error {
	// Read the template file as text
	templateFile, err := ioutil.ReadFile("internal/template.robot")
	if err != nil {
		return fmt.Errorf("failed to read template file: %v", err)
	}
	// Parse the template
	tmpl, err := template.New("template").Parse(string(templateFile))
	if err != nil {
		return fmt.Errorf("failed to parse template: %v", err)
	}
	// Render the template with the variables in the WebPage struct
	var renderedScript bytes.Buffer
	if err := tmpl.Execute(&renderedScript, webpage); err != nil {
		return fmt.Errorf("failed to render template: %v", err)
	}
	// Save the rendered script to a file
	if err := ioutil.WriteFile(fmt.Sprintf("./scripts/%s.robot", webpage.PageName), renderedScript.Bytes(), os.ModePerm); err != nil {
		return fmt.Errorf("failed to save rendered script to file: %v", err)
	}
	return nil
}

func scanDirectory(dir *string) []WebPage {

	// Find all YAML files in the directory
	yamlFiles, err := filepath.Glob(filepath.Join(*dir, "*.yml"))
	if err != nil {
		log.Fatal(err)
	}

	// Parse each YAML file and create a WebPage struct for each
	var webPages []WebPage
	for _, yamlFile := range yamlFiles {
		yamlData, err := ioutil.ReadFile(yamlFile)
		if err != nil {
			log.Printf("Error reading file %s: %s\n", yamlFile, err)
			continue
		}

		var webPage WebPage
		err = yaml.Unmarshal(yamlData, &webPage)
		if err != nil {
			log.Printf("Error parsing YAML in file %s: %s\n", yamlFile, err)
			continue
		}
		webPage.filename = yamlFile
		webPages = append(webPages, webPage)
	}

	// Print the WebPage structs to the console
	for _, webPage := range webPages {
		fmt.Printf("%+v\n", webPage)
	}

	return webPages
}

var parseCmd = &cobra.Command{
	Use:   "parse",
	Short: "Parse page config files as robot scripts",
	Run: func(cmd *cobra.Command, args []string) {

		var dir = ".pages"
		if len(args) > 0 {
			dir = args[0]
		}

		// Define command-line flags
		webPages := scanDirectory(&dir)

		for _, page := range webPages {
			err := generateRobotScript(&page)
			if err != nil {
				log.Fatal(err)
			}
		}
	},
}

func init() {
	rootCmd.AddCommand(parseCmd)
}
