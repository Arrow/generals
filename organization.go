package generals

import (
	"encoding/xml"
	"fmt"
	"log"
	"os"
)

var (
	filename = "org.xml"
	Groups   = make([]Group, 0)
)

type Group struct {
	Name     string
	Subgroup []Subgroup
	//Subgroup string
	//DefaultSub int
}

type Subgroup struct {
	SubName    string
	DefaultSub int
}

func init() {
	file, err := os.Open(filename)
	if err != nil {
		log.Fatal(err)
	}
	//groups := make([]Group,0)
	d := xml.NewDecoder(file)
	for {
		err = d.Decode(&Groups)
		if err != nil {
			return
		}
		fmt.Println(Groups)
	}
	fmt.Println(Groups)
}
