/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package cmd

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"regexp"
	"strings"

	"github.com/culturadevops/jgt/jfile"
	"github.com/culturadevops/jgt/jgit"
	"github.com/spf13/cobra"

	homedir "github.com/mitchellh/go-homedir"
	"github.com/spf13/viper"
)

var cfgFile string

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "ghr",
	Short: "A brief description of your application",
	Long: `A longer description that spans multiple lines and likely contains
examples and usage of using your application. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	//	Run: func(cmd *cobra.Command, args []string) { },
}

func MySearchFiles(rutaFinal string, resource_folder string) {
	//fmt.Println(mf.count)

	if _, err := os.Stat(rutaFinal); !os.IsNotExist(err) {

		files, err := ioutil.ReadDir(rutaFinal)
		if err != nil {
			log.Fatal(err)
		}
		cargarfile(rutaFinal, resource_folder)
		for _, f := range files {
			//mf.count++
			path := rutaFinal + "/" + f.Name()
			fmt.Println(path)
			if f.IsDir() {
				cargarfile(path, resource_folder)
				MySearchFiles(path, resource_folder)
			}
		}
	}
}

func cargarfile(folders_repos_init string, resource_folder string) {
	var J *jfile.Jfile
	J = new(jfile.Jfile)
	J.PrepareInitDefaultLog()
	/*J.Log = &jlog.Jlog{
		IsDebug:       false,
		PrinterLogs:   true,
		PrinterScreen: true,
	}
	J.Log.SetInitProperty()*/
	J.Map = make(map[string]string)
	/*busca el archivo Readme.md
	 no esta: copia recursos/Readme.md
	 si esta:
		buscar No tiene la palabra "# Mis Libros:" regexp.MatchString("# Mis Libros:", word)
			no esta: agrega el archivo al final	*/
	//folder base exist
	if J.FileExist(folders_repos_init) {
		var readmeName string
		readmeName = "/README.md"
		fmt.Println("si existe!, procediendo a buscar palabra...")
		archivodestino := folders_repos_init + readmeName
		fragmentoTexto := J.ReadFile(resource_folder + readmeName)
		// NOT file Readme exist
		if !J.FileExist(archivodestino) {
			fmt.Println("Creando archivo " + readmeName)
			J.CreateFile(archivodestino, fragmentoTexto)
		} else {
			existe := AddIfNotExist(J, archivodestino, "# Mis Libros:", fragmentoTexto)

			if !existe {
				fmt.Println("NO existe la palabra... PROCEDIENDO A AGREGAR")
			}
			if existe {
				fmt.Println("si existe la palabra...")
			}
		}
	}
}

/*todo:jaivic copia este en una clase*/
func AddIfNotExist(J *jfile.Jfile, DestineName string, MatchString string, ContentToAdd string) bool {
	word := J.ReadFile(DestineName)
	existe, _ := regexp.MatchString(MatchString, word)
	if !existe {
		J.AppEndToFile(DestineName, ContentToAdd)
	}
	return existe
}

/*todo:jaivic debes mover esta funcion a una clase especial para leer archivos*/
func GetJsonFileWithStruct(jsonFileName string, WithStruct interface{}) {
	jsonFile, err := os.Open(jsonFileName)
	if err != nil {
		if err.Error() == "open "+jsonFileName+": no such file or directory" {
			fmt.Println("no se encontro el archivo de configuracion")
		}
	}
	defer jsonFile.Close()
	byteValue, _ := ioutil.ReadAll(jsonFile)
	json.Unmarshal([]byte(byteValue), WithStruct)

}
func RemovePathForDir(filename string, path string) string {

	index := strings.Index(filename, path)

	leng := len(filename)

	if index > -1 {

		return filename[index+1 : leng]
	}
	return filename
}
func RemoveExtention(filename string, ext string) string {
	index := strings.Index(filename, ext)
	if index > -1 {
		return filename[0:index]
	}
	return filename
}
func recursiveIndexPath(filename string, path string) string {
	for i := 0; strings.Index(filename, path) > -1; i++ {
		filename = RemovePathForDir(filename, path)
	}
	return filename
}

type FilesStruct struct {
	Name  string
	Git   string
	Path  string
	Ready string
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	jg := &jgit.Jgit{}
	jg.PrepareInit()
	jg.Branch = "master"
	finalfolder := "test/"
	//resource_folder := "recursos"
	var baseRepoJson []FilesStruct
	GetJsonFileWithStruct("reposDescription.json", &baseRepoJson)

	for _, folder := range baseRepoJson {

		fmt.Println(folder.Git)
		JobName := RemoveExtention(recursiveIndexPath(folder.Git, "/"), ".git")
		fmt.Println(JobName)
		fmt.Println(jg)

		jg.CloneB(folder.Git, finalfolder+JobName)
		//MySearchFiles(finalfolder+JobName, resource_folder)
		/*jg.AddAll()
		jg.Commit("prueba de build")
		jg.Push("master")*/
		//jg.FinalPath = ""
	}

}

func init() {
	cobra.OnInitialize(initConfig)

	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.ghr.yaml)")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		home, err := homedir.Dir()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		// Search config in home directory with name ".ghr" (without extension).
		viper.AddConfigPath(home)
		viper.SetConfigName(".ghr")
	}

	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		fmt.Println("Using config file:", viper.ConfigFileUsed())
	}
}
