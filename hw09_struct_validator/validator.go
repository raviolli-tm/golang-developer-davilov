package hw09structvalidator

import (
	"fmt"
	"reflect"
)

type ValidationError struct {
	Field string
	Err   error
}

type Validators map[string]func(string) (bool, error)

var validators = Validators{
	"len": func(s string) (bool, error) {
		return true, nil
	},
	"min": func(s string) (bool, error) {
		return true, nil
	},
	"max": func(s string) (bool, error) {
		return true, nil
	},
	"in": func(s string) (bool, error) {
		return true, nil
	},
}

type ValidationErrors []ValidationError

func (v ValidationErrors) Error() string {
	panic("implement me")
}

func Validate(v interface{}) error {
	typeV := reflect.Indirect(reflect.ValueOf(v)).Type()
	err := reflectTypes(typeV)
	if err != nil {
		return err
	}

	fmt.Println()
	return nil
}

func reflectTypes(reflectType reflect.Type) error {
	fmt.Println(" ======= start reflectTypes")
	for i := 0; i < reflectType.NumField(); i++ {
		field := reflectType.Field(i)
		fmt.Print(field.Name, "  ")
		//if field.Tag.Get("validate") == "" {
		//	continue
		//}

		switch field.Type.Kind() {
		case reflect.String:
			fmt.Println("string")
		case reflect.Int:
			fmt.Println("int")
		case reflect.Array:
		case reflect.Slice:
			fmt.Print("slice | array")
			elemType := field.Type.Elem()
			fmt.Printf(", Element type: %s \n", elemType)
		case reflect.Struct:
			fmt.Println("struct")
			reflectTypes(field.Type)
		default:
			fmt.Println("Not supported type")

		}
	}
	fmt.Println(" ======= end reflectTypes")
	return nil
}

func validateField(reflectValue reflect.Value, tag string) error {
	return nil
}
