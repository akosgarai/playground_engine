package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
)

const (
	downLoadBaseUrl       = "https://raw.githubusercontent.com/akosgarai/playground_engine/master/"
	shaderBaseDirectory   = "pkg/shader/apps/"
	modelBaseDirectory    = "pkg/model/assets/"
	modelTargetDirectory  = "assets"
	shaderTargetDirectory = "shaders"
)

var (
	commands = map[string][]string{
		"help":    []string{},
		"install": []string{"shaders", "models"},
	}

	modelImages = []string{
		"concrete-wall.jpg",
		"crystal-ball.png",
		"door.jpg",
		"metal.jpg",
	}
	shaderApps = []string{
		"material.frag",
		"material.vert",
		"point.frag",
		"point.vert",
		"texturecolor.frag",
		"texturecolor.vert",
		"texture.frag",
		"texturemat.frag",
		"texturemat.vert",
		"texture.vert",
	}
)

func createTarget(targetPath string) {
	if _, err := os.Stat(targetPath); os.IsNotExist(err) {
		fmt.Printf("Creating '%s' directory.\n", targetPath)
		os.MkdirAll(targetPath, os.ModePerm)
	}
}
func fileExists(filepath string) bool {
	if _, err := os.Stat(filepath); os.IsNotExist(err) {
		return false
	}
	return true
}
func downloadShaders() {
	createTarget(shaderTargetDirectory)
	fmt.Println("Downloading shaders ...")
	for _, img := range shaderApps {
		target := shaderTargetDirectory + "/" + img
		if fileExists(target) {
			continue
		}
		url := fmt.Sprintf("%s%s%s", downLoadBaseUrl, shaderBaseDirectory, img)
		fmt.Printf("Downloading '%s'\n", url)
		err := downloadFile(url, target)
		if err != nil {
			fmt.Printf("Something happened during download: '%s'.\n", err.Error())
			fmt.Println("Please try to install again.")
		}
	}
}
func downloadFile(from, target string) error {
	// create the file
	to, err := os.Create(target)
	if err != nil {
		return err
	}
	defer to.Close()
	// get the data from the given url
	resp, err := http.Get(from)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// write  the downloaded data to the file
	_, err = io.Copy(to, resp.Body)
	if err != nil {
		return err
	}
	return nil
}
func downloadModels() {
	createTarget(modelTargetDirectory)
	fmt.Println("Downloading models ...")
	for _, img := range modelImages {
		target := modelTargetDirectory + "/" + img
		if fileExists(target) {
			continue
		}
		url := fmt.Sprintf("%s%s%s", downLoadBaseUrl, modelBaseDirectory, img)
		fmt.Printf("Downloading '%s'\n", url)
		err := downloadFile(url, target)
		if err != nil {
			fmt.Printf("Something happened during download: '%s'.\n", err.Error())
			fmt.Println("Please try to install again.")
		}
	}
}

func printCommandNames() {
	for cmd, _ := range commands {
		fmt.Printf("\t%s\n", cmd)

	}
}
func helpMenu() {
	fmt.Println("Welcome in the Playground Engine.\n")
	fmt.Println("This CLI tool helps You to get the necessary assets.")
	fmt.Println("How to use it?\n")
	fmt.Println("\tcmd [command] [params]\n")
	fmt.Println("Commands:\n")
	printCommandNames()
	fmt.Println("")
}
func helpMenuInstall() {
	fmt.Println("With the install tool, You can download the necessary stuff for your application.\n")
	for _, param := range commands["install"] {
		fmt.Printf("If You want to download the %s, use the following command:\n\n", param)
		fmt.Printf("\tplayground_engine install %s\n\n", param)
	}
}
func installCommandHandler(args []string) {
	if len(args) == 0 {
		helpMenuInstall()
		return
	}
	param := args[0]
	switch param {
	case "help":
		break
	case "shaders":
		downloadShaders()
		break
	case "models":
		downloadModels()
		break
	default:
	}
}

func main() {
	// Get the arguments without the program name
	Args := os.Args[1:]
	if len(Args) == 0 {
		helpMenu()
		return
	}
	Cmd := Args[0]
	switch Cmd {
	case "help":
		helpMenu()
		break
	case "install":
		installCommandHandler(Args[1:])
		break
	default:
		fmt.Println("Unsupported command. Please use one of the followings:\n")
		printCommandNames()
		break
	}

}