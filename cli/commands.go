package cli

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"

	"github.com/manifoldco/promptui"

	ct "github.com/exampleorg/config-tool/configtool"
)

type ConfigFlags struct {
	ProjectName string
	LocalPath   string
}

type Project struct {
	ProjectName          string `json:"projectname"`
	KinstaUserName       string `json:"kinstausername"`
	KinstaIP             string `json:"kinstaip"`
	KinstaPortStaging    string `json:"kinstaport-staging"`
	KinstaPortLive       string `json:"kinstaport-live"`
	ThemeName            string `json:"themename"`
	ImageSizeFacade      []int  `json:"imagesize-facade"`
	ImageSizeFloorplan   []int  `json:"imagesize-floorplan"`
	ImageSizeLot         []int  `json:"imagesize-lot"`
	Integrator           bool   `json:"integrator"`
	IntegratorUrlStaging string `json:"integratorurl-staging"`
	IntegratorUrlLive    string `json:"integratorurl-live"`
	LocalPath            string
}

// Library function to get a list of known projects

func GetProjects() ([]string, error) {
	path, err := ct.ConfigRootDir()
	if err != nil {
		return nil, err
	}

	configPath := fmt.Sprintf("%s/project-config/config", path)

	entries, err := os.ReadDir(configPath)
	if err != nil {
		return nil, err
	}

	var projects = make([]string, 0)
	for _, e := range entries {
		projects = append(projects, e.Name())
	}

	return projects, nil
}

// Show list of projects

func ListProjects() error {

	projects, err := GetProjects()
	if err != nil {
		return err
	}

	for _, e := range projects {
		fmt.Println(e)
	}

	return nil
}

func PrettyEncode(data interface{}, out io.Writer) error {
	enc := json.NewEncoder(out)
	enc.SetIndent("", "    ")
	if err := enc.Encode(data); err != nil {
		return err
	}
	return nil
}

// Show the config info of a selected project

func ProjectInfo(projectName string) error {

	prjConfig, err := GetProjectConfig(projectName)
	if err != nil {
		return err
	}

	tmpl := `

	Project Name: %s
	Kinsta User Name: %s
	Kinsta IP: %s
	Kinsta Port - Staging: %s
	Kinsta Port - Live: %s
	Theme Name: %s
	Image Size(W x H) - Facade: %d x %d
	Image Size(W x H) - Floorplan: %d x %d
	Image Size(W x H) - Lot: %d x %d
	Integrator: %t
	Intergrator Url - Staging: %s
	Intergrator Url - Live: %s
	Local Codebase: %s

`

	fmt.Printf(tmpl,
		prjConfig.ProjectName,
		prjConfig.KinstaUserName,
		prjConfig.KinstaIP,
		prjConfig.KinstaPortStaging,
		prjConfig.KinstaPortLive,
		prjConfig.ThemeName,
		prjConfig.ImageSizeFacade[0],
		prjConfig.ImageSizeFacade[1],
		prjConfig.ImageSizeFloorplan[0],
		prjConfig.ImageSizeFloorplan[1],
		prjConfig.ImageSizeLot[0],
		prjConfig.ImageSizeLot[1],
		prjConfig.Integrator,
		prjConfig.IntegratorUrlStaging,
		prjConfig.IntegratorUrlLive,
		prjConfig.LocalPath,
	)

	return nil
}

// Library function to get the config info of selected project

func GetProjectConfig(prjName string) (Project, error) {
	projectconfig := Project{}

	path, err := ct.ConfigRootDir()
	if err != nil {
		return projectconfig, err
	}

	// if target project not provided, prompt user
	project := prjName
	if prjName == "" {
		projlist, err := GetProjects()
		if err != nil {
			return projectconfig, err
		}

		t := promptui.Select{
			Label: "Select project",
			Items: projlist,
		}

		_, project, err = t.Run()
		if err != nil {
			return projectconfig, err
		}

		// fmt.Println(project, "selected")
	}

	configPath := fmt.Sprintf("%s/project-config/config", path)
	projectFilePath := fmt.Sprintf("%s/%s/project.json", configPath, project)
	localFilePath := fmt.Sprintf("%s/%s/local.json", configPath, project)

	//
	// read project config file
	//
	projectFile, err := os.Open(projectFilePath)
	if err != nil {
		return projectconfig, err
	}
	// fmt.Println(fmt.Sprintf("Successfully opened project.json for %s", project))
	defer projectFile.Close()

	byteValue, err := io.ReadAll(projectFile)
	if err != nil {
		return projectconfig, err
	}

	json.Unmarshal(byteValue, &projectconfig)

	//
	// read local settings config
	//
	localFile, err := os.Open(localFilePath)
	if err != nil {
		log.Println("warning - cannot find local.json config file for", project)
		log.Println("to create one run the 'hconfig setloc' or 'hconfig setlocwd' command")
		return projectconfig, nil
	}
	// fmt.Println(fmt.Sprintf("Successfully opened local.json for %s", project))
	defer localFile.Close()

	byteValue2, err := io.ReadAll(localFile)
	if err != nil {
		return projectconfig, err
	}

	json.Unmarshal(byteValue2, &projectconfig)

	return projectconfig, nil
}

// Manually set the local path

func SetLocal(flags *ConfigFlags) error {
	path, err := ct.ConfigRootDir()
	if err != nil {
		return err
	}

	configPath := fmt.Sprintf("%s/project-config/config", path)
	localFile := fmt.Sprintf("%s/%s/local.json", configPath, flags.ProjectName)

	lines := [4]string{
		"{",
		" \"projectname\": \"" + flags.ProjectName + "\",",
		" \"localpath\": \"" + flags.LocalPath + "\"",
		"}",
	}

	f, _ := os.Create(localFile)
	w := bufio.NewWriter(f)

	for i, s := range lines {
		fmt.Println(i, s)
		w.WriteString(s + "\n")
	}

	w.Flush()

	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Local config file created\n")

	return nil
}

// Set the local path based on current working directory

func SetLocalCodebase() error {
	path, err := ct.ConfigRootDir()
	if err != nil {
		return err
	}

	configPath := fmt.Sprintf("%s/project-config/config", path)
	localPath, err := os.Getwd()
	if err != nil {
		return err
	}

	projlist, err := GetProjects()
	if err != nil {
		return err
	}

	t := promptui.Select{
		Label: "Select project",
		Items: projlist,
	}

	_, project, err := t.Run()
	fmt.Println(project)

	localFile := fmt.Sprintf("%s/%s/local.json", configPath, project)

	lines := [4]string{
		"{",
		" \"projectname\": \"" + project + "\",",
		" \"localpath\": \"" + localPath + "/\"",
		"}",
	}

	f, _ := os.Create(localFile)
	w := bufio.NewWriter(f)

	for i, s := range lines {
		fmt.Println(i, s)
		w.WriteString(s + "\n")
	}

	w.Flush()

	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Local config file created\n")

	return nil
}
