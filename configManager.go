package configmanager

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"reflect"
	"strconv"
	"strings"
)

const defaultConfigFilePath = "default.conf"

// AssignConfiguration - pass the reference of struct (&s) as an argument, the struct will be modified
func AssignConfiguration(s interface{}) {

	var err error
	ConfigFilePath := ""
	configs := []string{}
	args := os.Args
	if len(args) > 2 {
		if args[1] == "--config-file-path" || args[1] == "-f" {
			fmt.Println(args[2])
			ConfigFilePath = args[2]
		}
	}

	if ConfigFilePath != "" {
		configs, err = readConfigurationFile(ConfigFilePath)
		if err != nil {
			log.Fatalf("Error while reading specified file < %v > :: %v", ConfigFilePath, err)
		}

	} else {
		// if config file path is not specified search in current directory
		_, err = os.Stat(defaultConfigFilePath)
		if err != nil {
			log.Printf(" < default.conf > File not found in current directory :: %v", err)
			log.Fatalf(`ABORTING PROCESS - NEITHER ANY CONFIG FILE IS SPECIFIED NOR A DEFAULT CONFIG FILE < default.conf > IS FOUND IN CURRENT DIRECTORY
To specify a config file, run the program with argument -f <path of config file *.conf>
To use a default config file, create a 'default.conf' file in the working directory			
			`)
		} else {
			ConfigFilePath = defaultConfigFilePath
			configs, err = readConfigurationFile(ConfigFilePath)
			if err != nil {
				log.Fatalf("Error while Reading Config File < %v > :: %v", ConfigFilePath, err)
			}
		}

	}

	// fmt.Println(configs)

	//assign default
	assignConfiguration(configs, s)

}

//raed a file
func readConfigurationFile(filePath string) ([]string, error) {

	configs := []string{}

	fl, err := os.Open(filePath)
	if err != nil {
		return configs, err
	}
	defer fl.Close()

	scanner := bufio.NewScanner(fl)
	lineNo := 0
	for scanner.Scan() {
		lineNo++
		newText := scanner.Text()
		newText = strings.Trim(newText, " ")
		newText = strings.ReplaceAll(newText, " ", "")
		newText = strings.Trim(newText, "\n")

		if strings.Contains(newText, "#") {
			newText = strings.Split(newText, "#")[0]
		}

		if newText == "" {
			continue
		}

		if strings.Count(newText, "=") != 1 {
			log.Fatalf("Invalid configuration at line No: %v < %v >", lineNo, newText)
		}

		configs = append(configs, newText)

	}

	if err = scanner.Err(); err != nil {
		log.Printf("%v", err)
	}

	return configs, nil

}

//config assignment
func assignConfiguration(configs []string, s interface{}) {

	configs = reverse(configs)

	v := reflect.ValueOf(s).Elem()
	typeOfS := v.Type()

	for i := 0; i < v.NumField(); i++ {
		f := v.Field(i)
		// fmt.Println(f.Kind(), typeOfS.Field(i).Name)
		switch type_ := f.Kind(); type_ {

		case reflect.String:
			for ix, vl := range configs {

				_temp := strings.Split(vl, "=")
				key, val := _temp[0], _temp[1]

				if key == typeOfS.Field(i).Name {
					v.FieldByName(key).SetString(val)
					break
				}

				if ix+1 == len(configs) {
					log.Fatalf("ERROR: Parameter < %v > is missing in configs", typeOfS.Field(i).Name)
				}

			}

		case reflect.Int64:
			for ix, vl := range configs {

				_temp := strings.Split(vl, "=")
				key, val := _temp[0], _temp[1]

				if key == typeOfS.Field(i).Name {
					nval, err := strconv.Atoi(val)
					if err != nil {
						log.Fatalf("ERROR: Invalid parameters < %v > in Configs :: %v", key, err)
					}
					v.FieldByName(key).SetInt(int64(nval))
					break
				}

				if ix+1 == len(configs) {
					log.Fatalf("ERROR: Parameter < %v > is missing in configs", typeOfS.Field(i).Name)
				}

			}

		case reflect.Float64:
			for ix, vl := range configs {

				_temp := strings.Split(vl, "=")
				key, val := _temp[0], _temp[1]

				if key == typeOfS.Field(i).Name {
					nval, err := strconv.ParseFloat(val, 64)
					if err != nil {
						log.Fatalf("ERROR: Invalid parameters < %v > in configs :: %v", key, err)
					}
					v.FieldByName(key).SetFloat(nval)
					break
				}

				if ix+1 == len(configs) {
					log.Fatalf("ERROR: Parameter < %v > is missing in configs", typeOfS.Field(i).Name)
				}

			}

		case reflect.Bool:
			for ix, vl := range configs {

				_temp := strings.Split(vl, "=")
				key, val := _temp[0], _temp[1]

				if key == typeOfS.Field(i).Name {
					nval, err := strconv.ParseBool(val)
					if err != nil {
						log.Fatalf("ERROR: Invalid Parameters @ < %v > in Configs :: %v", key, err)
					}
					v.FieldByName(key).SetBool(nval)
					break
				}

				if ix+1 == len(configs) {
					log.Fatalf("ERROR: Parameter < %v > is missing in configs", typeOfS.Field(i).Name)
				}

			}

		}

	}
	// fmt.Println(s)

}

// reverse string slice function
func reverse(slice []string) []string {

	newSlice := []string{}
	for ix, _ := range slice {
		newSlice = append(newSlice, slice[len(slice)-1-ix])
	}

	return newSlice
}
