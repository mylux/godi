package container

import (
	"fmt"
	"reflect"
)

var registry map[reflect.Type]interface{} = map[reflect.Type]interface{}{}

func Wire(inter interface{}, impl interface{}) error {
	var fType reflect.Type

	if fType = reflect.TypeOf(inter); fType.Kind() == reflect.Ptr {
		fType = reflect.Indirect(reflect.ValueOf(inter)).Type()
	}

	if fType.Kind() == reflect.Interface {
		constructor := reflect.MakeFunc(reflect.FuncOf([]reflect.Type{}, []reflect.Type{fType}, false), func(args []reflect.Value) (results []reflect.Value) {
			return []reflect.Value{reflect.ValueOf(impl)}
		})
		registry[fType] = constructor.Interface()
		return nil
	} else {
		return fmt.Errorf("%v is not an interface, but %v instead", fType, fType.Kind())
	}
}

func WireFactory(fn interface{}) error {
	fType := reflect.TypeOf(fn)

	if fType.Kind() != reflect.Func {
		return fmt.Errorf("argument is not a function")
	}

	for i := 0; i < fType.NumOut(); i++ {
		rType := fType.Out(i)
		if _, exists := registry[rType]; exists {
			return fmt.Errorf("type %s already defined", rType.String())
		}
		registry[rType] = fn
	}

	return nil
}

func AutoWire(obj interface{}) error {
	value := reflect.ValueOf(obj).Elem()
	oType := value.Type()

	if oType.Kind() == reflect.Interface {
		constructor, exists := registry[oType]
		if !exists {
			return fmt.Errorf("type %v not implemented", oType)
		}
		oType = findImplementation(constructor, oType)
	}

	for i := 0; i < oType.NumField(); i++ {
		field := oType.Field(i)
		if _, exists := field.Tag.Lookup("autowired"); exists {
			fieldType := field.Type
			constructedItems, err := construct(fieldType)
			if err != nil {
				return err
			}
			r := getFirstByType(constructedItems, fieldType)
			if r != nil {
				value.Field(i).Set(*r)
			}
			if fieldType.Kind() == reflect.Struct {
				AutoWire(value.Field(i).Addr().Interface())
			}
		}
	}
	return nil
}

func Construct[T any]() *T {
	fType := reflect.TypeFor[T]()
	rs, err := construct(fType)

	if err != nil {
		return nil
	}

	if len(rs) > 0 {
		var found bool
		var r T
		for _, v := range rs {
			r, found = v.Interface().(T)
			if found {
				break
			}
		}
		if found {
			return &r
		} else {
			return nil
		}
	} else {
		return nil
	}
}

func CleanRegistry() {
	for t := range registry {
		delete(registry, t)
	}

}

func findImplementation(constructor interface{}, t reflect.Type) reflect.Type {
	tConstructor := reflect.TypeOf(constructor)
	for i := 0; i < tConstructor.NumOut(); i++ {
		if tConstructor.Out(i).AssignableTo(t) {
			return tConstructor.Out(i)
		}
	}
	return nil
}

func construct(fType reflect.Type) ([]reflect.Value, error) {
	if constructor, exists := registry[fType]; exists {
		constructorType := reflect.TypeOf(constructor)
		argSize := constructorType.NumIn()
		var args []reflect.Value
		if argSize > 0 {
			args = make([]reflect.Value, argSize)
			for i := range argSize {
				a, err := construct(constructorType.In(i))
				if err != nil {
					return nil, err
				}
				args[i] = *getFirstByType(a, constructorType.In(i))
			}
		} else {
			args = nil
		}
		return reflect.ValueOf(constructor).Call(args), nil
	} else if fType.Kind() != reflect.Interface {
		return []reflect.Value{reflect.Indirect(reflect.New(fType))}, nil
	}
	return nil, fmt.Errorf("no constructor found for type %v", fType)
}

func getFirstByType(values []reflect.Value, oType reflect.Type) *reflect.Value {
	for _, v := range values {
		if v.Type() == oType {
			return &v
		}
	}
	return nil
}
