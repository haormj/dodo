package util

import (
	"log"
	"reflect"
	"testing"
)

func TestInitPointer(t *testing.T) {
	log.SetFlags(log.Lshortfile)
	var str *string
	ss := InitPointer(reflect.TypeOf(str))
	log.Println(ss.Interface())
}

func TestSetStructField(t *testing.T) {
	type World struct {
		Name string
	}
	type Hello struct {
		Name  string
		World World
	}
	hello := &Hello{}
	ok := SetStructField(hello, "I'm Hello Name", "Name")
	log.Println(ok, hello)
	ok = SetStructField(hello, "I'm World Name", "World", "Name")
	log.Println(ok, hello)
}

func TestGetStructField(t *testing.T) {
	type World struct {
		Name string
	}
	type Hello struct {
		Name  *string
		World World
	}
	xx := "xxxxxx"
	hello := &Hello{
		Name: &xx,
		World: World{
			Name: "I'm World Name",
		},
	}
	helloName, ok := GetStructField(hello, "Name").String()
	if ok {
		log.Println("Hello.Name: ", helloName)
	}
	worldName, ok := GetStructField(hello, "World", "Name").String()
	if ok {
		log.Println("Hello.World.Name: ", worldName)
	}
}
