package configmanager

import (
	"bufio"
	"log"
	"os"
	"reflect"
	"strconv"
	"strings"
)

// DefaultConfigFilePath - The default configuration file path
var DefaultConfigFilePath = "default.conf"
// DefaultSliceParametersSeparator - The default Slice Parameters Separator 
var DefaultSliceParametersSeparator = "|"

// AssignConfiguration - pass the reference of struct (&s) as an argument, the struct will be modified
func AssignConfiguration(configStruct interface{}) {

	var err error
	ConfigFilePath := ""
	configs := []string{}
	args := os.Args
	if len(args) > 2 {
		if args[1] == "--config-file-path" || args[1] == "-f" {
			log.Println("Reading config file :: ", args[2])
			ConfigFilePath = args[2]
		}
	}

	if ConfigFilePath != "" {
		configs, err = readConfigurationFile(ConfigFilePath)
		if err != nil {
			log.Fatalf("Error while reading specified file < %v > :: %v", ConfigFilePath, err)
		}

	} else {
		// if config file path is not specified read defaultconfigfile
		log.Println("Reading default config file :: ", DefaultConfigFilePath)
		_, err = os.Stat(DefaultConfigFilePath)
		if err != nil {

			log.Printf("Error: Could not find < %v > file :: %v", DefaultConfigFilePath, err)
			log.Fatalf(`ABORTING PROCESS - NEITHER ANY CONFIG FILE IS SPECIFIED NOR THE DEFAULT CONFIG FILE < %v > IS FOUND.
To specify a config file, run the program with arguments : -f <path of config file *.conf>
To use a default config file, create a  < %v > file.			
			`, DefaultConfigFilePath, DefaultConfigFilePath)
		} else {
			ConfigFilePath = DefaultConfigFilePath
			configs, err = readConfigurationFile(ConfigFilePath)
			if err != nil {
				log.Fatalf("Error while Reading Config File < %v > :: %v", ConfigFilePath, err)
			}
		}

	}

	// fmt.Println(configs)

	//assign default
	assignConfiguration(configs, configStruct)

}

//read config file
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
		switch f.Type().String() {

		case "string":
			for ix, vl := range configs {

				_temp := strings.Split(vl, "=")
				key, val := _temp[0], _temp[1]

				if key == typeOfS.Field(i).Name {
					val = strings.Trim(val, " ") // trim white spaces
					v.FieldByName(key).SetString(val)
					break
				}

				if ix+1 == len(configs) {
					log.Fatalf("ERROR: Parameter < %v > is missing in configs", typeOfS.Field(i).Name)
				}

			}

		case "int64":
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

		case "float64":
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

		case "bool":
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

		case "[]string":

			for ix, vl := range configs {

				_temp := strings.Split(vl, "=")
				key, val := _temp[0], _temp[1]

				if key == typeOfS.Field(i).Name {
					valSlice := []string{}
					
					for _,prm := range strings.Split(val, DefaultSliceParametersSeparator){
						prm = strings.Trim(prm," ")
						//filter empty parameters
						if prm != ""{
							valSlice = append(valSlice,prm)
						}
					}

					slice := reflect.MakeSlice(reflect.TypeOf([]string{}), len(valSlice), len(valSlice))

					for ix, vl_ := range valSlice {
						slice.Index(ix).SetString(string(vl_))
					}
					v.FieldByName(key).Set(slice)
					break
				}

				if ix+1 == len(configs) {
					log.Fatalf("ERROR: Parameter < %v > is missing in configs", typeOfS.Field(i).Name)
				}

			}

		case "[]float64":

			for ix, vl := range configs {

				_temp := strings.Split(vl, "=")
				key, val := _temp[0], _temp[1]

				if key == typeOfS.Field(i).Name {
					valSlice := []string{}
					
					for _,prm := range strings.Split(val, DefaultSliceParametersSeparator){
						prm = strings.Trim(prm," ")
						//filter empty parameters
						if prm != ""{
							valSlice = append(valSlice,prm)
						}
					}

					slice := reflect.MakeSlice(reflect.TypeOf([]float64{}), len(valSlice), len(valSlice))

					for ix, vl_ := range valSlice {
						fval, err := strconv.ParseFloat(vl_, 64)
						if err != nil {
							log.Fatalf("ERROR: Invalid parameters < %v > in configs :: %v", key, err)
						}
						slice.Index(ix).SetFloat(fval)
					}
					v.FieldByName(key).Set(slice)
					break
				}

				if ix+1 == len(configs) {
					log.Fatalf("ERROR: Parameter < %v > is missing in configs", typeOfS.Field(i).Name)
				}

			}

		case "[]int64":

			for ix, vl := range configs {

				_temp := strings.Split(vl, "=")
				key, val := _temp[0], _temp[1]

				if key == typeOfS.Field(i).Name {
					valSlice := []string{}
					
					for _,prm := range strings.Split(val, DefaultSliceParametersSeparator){
						prm = strings.Trim(prm," ")
						//filter empty parameters
						if prm != ""{
							valSlice = append(valSlice,prm)
						}
					}

					slice := reflect.MakeSlice(reflect.TypeOf([]int64{}), len(valSlice), len(valSlice))

					for ix, vl_ := range valSlice {
						ival, err := strconv.Atoi(vl_)
						if err != nil {
							log.Fatalf("ERROR: Invalid parameters < %v > in Configs :: %v", key, err)
						}
						slice.Index(ix).SetInt(int64(ival))
					}
					v.FieldByName(key).Set(slice)
					break
				}

				if ix+1 == len(configs) {
					log.Fatalf("ERROR: Parameter < %v > is missing in configs", typeOfS.Field(i).Name)
				}

			}

		default:
			log.Fatalf("ERROR: Type of < %v > is < %v >, which is not supported", typeOfS.Field(i).Name, f.Type().String())

		}

	}
	// fmt.Println(s)

}

// reverse string slice function
func reverse(slice []string) []string {

	newSlice := []string{}
	for ix := range slice {
		newSlice = append(newSlice, slice[len(slice)-1-ix])
	}

	return newSlice
}
