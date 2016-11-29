package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"os/user"
	"strings"
)

var flavor map[string]string

func extractData(file os.FileInfo) error {
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
	err = readAndReplace(dat)
	if err != nil {
		fmt.Println("\x1b[31;1m", err, "\x1b[0m")
	} else {
		fmt.Println("\x1b[32;1mok\x1b[0m")
	}

	execute(dat["exec"].([]interface{}))

	return nil
}

func execute(commands []interface{}) {
	usr, err := user.Current()
	if err != nil {
		fmt.Println(err)
	}

	home := usr.HomeDir

	for _, command := range commands {
		fmt.Print("  ", command, ": ")
		com := strings.Split(command.(string), " ")
		var c *exec.Cmd
		if len(com) == 1 {
			c = exec.Command(com[0])
		} else {
			c = exec.Command(com[0], com[1:]...)
		}
		c.Dir = home
		err := c.Run()
		if err != nil {
			fmt.Println("\x1b[31;1m", err, "\x1b[0m")
		} else {
			fmt.Println("\x1b[32;1mok\x1b[0m")
		}
	}
}

func readAndReplace(file map[string]interface{}) error {
	usr, err := user.Current()
	if err != nil {
		return (err)
	}

	home := usr.HomeDir

	loc := home + "/" + file["location"].(string)
	input, err := ioutil.ReadFile(loc)
	if err != nil {
		return err
	}

	lines := file["lines"].([]interface{})
	li := strings.Split(string(input), "\n")

	for _, line := range lines {
		line := line.(map[string]interface{})
		prev := line["prev"].(string)
		for i, l := range li {
			if strings.Contains(l, prev) {
				next := line["next"].(string)
				colors := line["colors"].([]interface{})
				for j, color := range colors {
					colors[j] = flavor[color.(string)]
				}
				li[i] = fmt.Sprintf(next, colors...)
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
			flavor["name"] = fl.Name()[:len(fl.Name())-5]
		}
	}

	if flavor == nil {
		flavorNotFound(args[0], flavors)
		return
	}

	if len(args) > 1 && args[1] == "light" {
		flavor["darkness"] = "light"
	} else {
		flavor["darkness"] = "dark"
	}

	dir = home + "/.theme/templates/"
	templates, err := ioutil.ReadDir(dir)
	if err != nil {
		fmt.Println(err)
		return
	}

	for _, template := range templates {
		err := extractData(template)
		if err != nil {
			fmt.Println(err)
			return
		}
	}

}
