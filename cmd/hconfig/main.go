package main

import (
	_ "embed"
	"log"

	cl "github.com/exampleorg/config-tool/cli"
	ct "github.com/exampleorg/config-tool/configtool"

	"github.com/leaanthony/clir"
)

//go:embed version.txt
var version string

func main() {

	var projectName, localPath string

	// ivytp CLI created
	cli := clir.NewCli("hconfig", "Project Config Tool", version)

	//
	// capture flags
	//
	cli.StringFlag("projectname", "project config to update", &projectName).
		StringFlag("localpath", "local environment codebase folder", &localPath)

	// init command
	initCmd := cli.NewSubCommandInheritFlags("init", "Manage project configuration")
	initCmd.Action(func() error {
		return ct.PullProjectFiles()
	})

	// update command
	updateCmd := cli.NewSubCommandInheritFlags("update", "Get latest project updates")
	updateCmd.Action(func() error {
		return ct.PullProjectFiles()
	})

	// list projects command
	lsCmd := cli.NewSubCommandInheritFlags("ls", "Show list of known projects")
	lsCmd.Action(func() error {
		return cl.ListProjects()
	})

	// project info command
	infoCmd := cli.NewSubCommandInheritFlags("info", "Select project and show config details")
	infoCmd.Action(func() error {
		return cl.ProjectInfo(projectName)
	})

	// set local command
	setlocalCmd := cli.NewSubCommandInheritFlags("setlocal", "Set local codebase path for project")
	setlocalCmd.Action(func() error {
		// capture the flags provided
		cf := &cl.ConfigFlags{
			ProjectName: projectName,
			LocalPath:   localPath,
		}
		return cl.SetLocal(cf)
	})

	// set local path using current folder command
	setlocalbycwdCmd := cli.NewSubCommandInheritFlags("setlocalwd", "Set project local codebase path using current folder path")
	setlocalbycwdCmd.Action(func() error {
		return cl.SetLocalCodebase()
	})

	err := cli.Run()
	if err != nil {
		log.Fatal(err)
	}
}
