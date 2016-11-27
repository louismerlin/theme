package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"os/user"
	"strings"
)

func extractData(file os.FileInfo, flavor map[string]string) error {
	usr, err := user.Current()
	if err != nil {
		return (err)
	}

	home := usr.HomeDir

	dir := home + "/.theme/templates/" + file.Name()
	input, err := ioutil.ReadFile(dir)
	if err != nil {
		return (err)
	}

	var dat map[string]interface{}
	if err := json.Unmarshal(input, &dat); err != nil {
		return err
	}

	fmt.Print(dat["name"], ": ")
	err = readAndReplace(dat, flavor)
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println("ok")
	}

	return nil
}

func readAndReplace(file map[string]interface{}, flavor map[string]string) error {
	loc := file["location"].(string)
	input, err := ioutil.ReadFile(loc)
	if err != nil {
		return err
	}

	lines := file["lines"].([]interface{})

	for _, line := range lines {
		line := line.(map[string]interface{})
		prev := line["prev"].(string)
		li := strings.Split(string(input), "\n")
		for i, l := range li {
			if strings.Contains(l, prev) {
				next := line["next"].(string)
				color := flavor[line["color"].(string)]
				li[i] = fmt.Sprintf(next, color)
			}
		}
		output := strings.Join(li, "\n")
		err = ioutil.WriteFile(loc, []byte(output), 0644)
		if err != nil {
			return err
		}
	}

	return nil
}

func flavorNotFound(flavor string, flavors []os.FileInfo) {
	if flavor != "" {
		fmt.Println("Flavor not found : ", flavor)
	}

	for _, fl := range flavors {
		fmt.Println(fl.Name()[:len(fl.Name())-5])
	}
}

func main() {
	args := os.Args[1:]

	usr, err := user.Current()
	if err != nil {
		fmt.Println(err)
		return
	}

	home := usr.HomeDir

	dir := home + "/.theme/flavors/"
	flavors, err := ioutil.ReadDir(dir)
	if err != nil {
		fmt.Println(err)
		return
	}

	if len(args) == 0 {
		flavorNotFound("", flavors)
		return
	}

	var flavor map[string]string
	for _, fl := range flavors {
		if fl.Name()[:len(fl.Name())-5] == args[0] {
			dir = home + "/.theme/flavors/" + fl.Name()
			input, err := ioutil.ReadFile(dir)
			if err != nil {
				fmt.Println(err)
				return
			}
			if err := json.Unmarshal(input, &flavor); err != nil {
				fmt.Println(err)
				return
			}
		}
	}

	if flavor == nil {
		flavorNotFound(args[0], flavors)
		return
	}

	dir = home + "/.theme/templates/"
	templates, err := ioutil.ReadDir(dir)
	if err != nil {
		fmt.Println(err)
		return
	}

	for _, template := range templates {
		err := extractData(template, flavor)
		if err != nil {
			fmt.Println(err)
			return
		}
	}

}
