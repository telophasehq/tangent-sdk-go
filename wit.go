package tangent_sdk

import (
	"fmt"
	"reflect"
	"strings"
	"time"
)

func WITOf[T any](typName string) (string, error) {
	var zero T
	return witOfType(reflect.TypeOf(zero), typName)
}

func witOfType(rt reflect.Type, name string) (string, error) {
	def, err := witTypedef(rt)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("/// auto-generated\nrecord %s {\n%s}\n", name, def), nil
}

func witTypedef(rt reflect.Type) (string, error) {
	switch rt.Kind() {
	case reflect.Bool:
		return "  // bool\n", nil
	case reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
		reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return "  // integer\n", nil
	case reflect.Float32, reflect.Float64:
		return "  // float\n", nil
	case reflect.String:
		return "  // string\n", nil
	case reflect.Slice:
		// []byte -> list<u8>, else list<...>
		if rt.Elem().Kind() == reflect.Uint8 {
			return "  // list<u8>\n", nil
		}
		inner, err := witFieldType(rt.Elem())
		if err != nil {
			return "", err
		}
		return fmt.Sprintf("  // list<%s>\n", inner), nil
	case reflect.Map:
		// only map[string]T supported → list<tuple<string,T>>
		if rt.Key().Kind() != reflect.String {
			return "", fmt.Errorf("only map[string] supported")
		}
		inner, err := witFieldType(rt.Elem())
		if err != nil {
			return "", err
		}
		return fmt.Sprintf("  // list<tuple<string,%s>>\n", inner), nil
	case reflect.Pointer:
		// *T -> option<T>
		inner, err := witFieldType(rt.Elem())
		if err != nil {
			return "", err
		}
		return fmt.Sprintf("  // option<%s>\n", inner), nil
	case reflect.Struct:
		// special-case time.Time → i64 millis
		if rt.PkgPath() == "time" && rt.Name() == "Time" {
			return "  // i64  // epoch_ms\n", nil
		}
		var b strings.Builder
		for i := 0; i < rt.NumField(); i++ {
			f := rt.Field(i)
			if f.PkgPath != "" {
				continue
			} // unexported
			n := f.Tag.Get("json")
			if n == "" {
				n = f.Name
			}
			tstr, err := witFieldType(f.Type)
			if err != nil {
				return "", err
			}
			fmt.Fprintf(&b, "  %s: %s,\n", jsonName(n), tstr)
		}
		return b.String(), nil
	default:
		return "", fmt.Errorf("unsupported kind: %s", rt.Kind())
	}
}

func witFieldType(rt reflect.Type) (string, error) {
	switch rt.Kind() {
	case reflect.Bool:
		return "bool", nil
	case reflect.Int8:
		return "s8", nil
	case reflect.Int16:
		return "s16", nil
	case reflect.Int32:
		return "s32", nil
	case reflect.Int64:
		return "s64", nil
	case reflect.Uint8:
		return "u8", nil
	case reflect.Uint16:
		return "u16", nil
	case reflect.Uint32:
		return "u32", nil
	case reflect.Uint64:
		return "u64", nil
	case reflect.Float32:
		return "float32", nil
	case reflect.Float64:
		return "float64", nil
	case reflect.String:
		return "string", nil
	case reflect.Slice:
		if rt.Elem().Kind() == reflect.Uint8 {
			return "list<u8>", nil
		}
		inner, _ := witFieldType(rt.Elem())
		return "list<" + inner + ">", nil
	case reflect.Map:
		if rt.Key().Kind() != reflect.String {
			return "", fmt.Errorf("map key must be string")
		}
		inner, _ := witFieldType(rt.Elem())
		return "list<tuple<string," + inner + ">>", nil
	case reflect.Pointer:
		inner, _ := witFieldType(rt.Elem())
		return "option<" + inner + ">", nil
	case reflect.Struct:
		if rt == reflect.TypeOf(time.Time{}) {
			return "s64", nil
		} // epoch_ms
		// nested anonymous struct → inline record
		def, _ := witTypedef(rt)
		return "record {\n" + def + "}", nil
	default:
		return "", fmt.Errorf("unsupported: %s", rt.Kind())
	}
}

func jsonName(n string) string {
	if idx := strings.Index(n, ","); idx >= 0 {
		n = n[:idx]
	}
	if n == "-" {
		n = ""
	}
	return n
}
