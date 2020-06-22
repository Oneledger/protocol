package main

import (
	"database/sql"
	"errors"
	"fmt"
	"reflect"

	"github.com/Oneledger/protocol/chains/bitcoin"
	"github.com/Oneledger/protocol/chains/ethereum"
	"github.com/Oneledger/protocol/consensus"
	"github.com/Oneledger/protocol/data/fees"
	"github.com/Oneledger/protocol/data/governance"
	"github.com/Oneledger/protocol/data/ons"
)

func main() {
	updateOptions()
}

func explore(s interface{}) {
	reflectType := reflect.TypeOf(s).Elem()
	reflectValue := reflect.ValueOf(s).Elem()
	//fmt.Println("Exploring : ", reflectType.Name())

	for i := 0; i < reflectType.NumField(); i++ {
		typeName := reflectType.Field(i).Name
		valueType := reflectValue.Field(i).Type()
		valueKind := reflectValue.Field(i).Kind()
		v := reflectValue.Field(i).Addr()
		if typeName == "ERCContractAddress" {
			fmt.Printf("%s | %s | %s \n", typeName, valueType, valueKind)
			if reflectValue.Field(i).CanSet() {
				switch valueKind {
				case reflect.Int:
					fmt.Println("Int : ", v.Elem())
					v.SetInt(1)
				case reflect.Array:
					fmt.Println("Array : ", v.Elem())
					newaddr := []byte{0x0}
					v.SetBytes(newaddr)
				}
			}
		}
		if reflectValue.Field(i).Kind() == reflect.Struct {
			explore(v.Interface())
		}
	}
}

func deepFields(reflectType reflect.Type) []reflect.StructField {
	var fields []reflect.StructField
	if reflectType = indirectType(reflectType); reflectType.Kind() == reflect.Struct {
		for i := 0; i < reflectType.NumField(); i++ {
			v := reflectType.Field(i)
			if v.Anonymous {
				fields = append(fields, deepFields(v.Type)...)
			} else {
				fields = append(fields, v)
			}
		}
	}

	return fields
}
func updateOptions() {

	gov := consensus.GovernanceState{
		FeeOption:   fees.FeeOption{},
		ETHCDOption: ethereum.ChainDriverOption{},
		BTCCDOption: bitcoin.ChainDriverOption{},
		ONSOptions:  ons.Options{},
		PropOptions: governance.ProposalOptionSet{},
	}

	//updates := map[string]string{}
	//updates["Currency"] = "TestNew"
	//for k, v := range updates {
	//	fmt.Println(k, v)
	//}

	//s := reflect.ValueOf(&gov).Elem()
	fmt.Println(gov.FeeOption.FeeCurrency.Chain)
	explore(&gov)
	fmt.Println(gov.FeeOption.FeeCurrency.Chain)
	//typeof := s.Type()

	//o := reflect.ValueOf(&options).Elem()
	//reflect.ValueOf(&options).Elem().FieldByName(update.Key).SetString(update.Value)
	//fmt.Printf("%+v\n", options)

	//ps := reflect.ValueOf(&gov)
	//// struct
	//s := ps.Elem()
	//if s.Kind() == reflect.Struct {
	//	// exported field
	//	f := s.FieldByName(gov.Key)
	//	if f.IsValid() {
	//		if f.CanSet() {
	//			// change value of N
	//			if f.Kind() == reflect.String {
	//				x := gov.Value
	//				f.SetString(x)
	//			}
	//		}
	//	}
	//}

	//updateOption := fmt.Sprintf("{%s : %s}", update.Key, update.Value)
	//fmt.Println(updateOption)
	//updateData := []byte(`
	//{
	//	"currency": "Test"
	//}
	//`)
	//
	//var updatedFields ons.Options
	//err := json.Unmarshal(updateData, &updatedFields)
	//if err != nil {
	//	fmt.Println("An error occured: %v", err)
	//}
	//fmt.Printf("%+v\n", updatedFields)
}

func Copy(toValue interface{}, fromValue interface{}) (err error) {
	var (
		isSlice bool
		amount  = 1
		from    = indirect(reflect.ValueOf(fromValue))
		to      = indirect(reflect.ValueOf(toValue))
	)

	if !to.CanAddr() {
		return errors.New("copy to value is unaddressable")
	}

	// Return is from value is invalid
	if !from.IsValid() {
		return
	}

	fromType := indirectType(from.Type())
	toType := indirectType(to.Type())

	// Just set it if possible to assign
	// And need to do copy anyway if the type is struct
	if fromType.Kind() != reflect.Struct && from.Type().AssignableTo(to.Type()) {
		to.Set(from)
		return
	}

	if fromType.Kind() != reflect.Struct || toType.Kind() != reflect.Struct {
		return
	}

	if to.Kind() == reflect.Slice {
		isSlice = true
		if from.Kind() == reflect.Slice {
			amount = from.Len()
		}
	}

	for i := 0; i < amount; i++ {
		var dest, source reflect.Value

		if isSlice {
			// source
			if from.Kind() == reflect.Slice {
				source = indirect(from.Index(i))
			} else {
				source = indirect(from)
			}
			// dest
			dest = indirect(reflect.New(toType).Elem())
		} else {
			source = indirect(from)
			dest = indirect(to)
		}

		// check source
		if source.IsValid() {
			fromTypeFields := deepFields(fromType)
			fmt.Printf("%#v", fromTypeFields)
			// Copy from field to field or method
			for _, field := range fromTypeFields {
				name := field.Name

				if fromField := source.FieldByName(name); fromField.IsValid() {
					// has field
					if toField := dest.FieldByName(name); toField.IsValid() {
						if toField.CanSet() {
							if !fromField.IsZero() {
								if !set(toField, fromField) {
									if err := Copy(toField.Addr().Interface(), fromField.Interface()); err != nil {
										return err
									}
								}
							}
						}
					} else {
						// try to set to method
						var toMethod reflect.Value
						if dest.CanAddr() {
							toMethod = dest.Addr().MethodByName(name)
						} else {
							toMethod = dest.MethodByName(name)
						}

						if toMethod.IsValid() && toMethod.Type().NumIn() == 1 && fromField.Type().AssignableTo(toMethod.Type().In(0)) {
							toMethod.Call([]reflect.Value{fromField})
						}
					}
				}
			}

			// Copy from method to field
			for _, field := range deepFields(toType) {
				name := field.Name

				var fromMethod reflect.Value
				if source.CanAddr() {
					fromMethod = source.Addr().MethodByName(name)
				} else {
					fromMethod = source.MethodByName(name)
				}

				if fromMethod.IsValid() && fromMethod.Type().NumIn() == 0 && fromMethod.Type().NumOut() == 1 {
					if toField := dest.FieldByName(name); toField.IsValid() && toField.CanSet() {
						values := fromMethod.Call([]reflect.Value{})
						if len(values) >= 1 {
							set(toField, values[0])
						}
					}
				}
			}
		}
		if isSlice {
			if dest.Addr().Type().AssignableTo(to.Type().Elem()) {
				to.Set(reflect.Append(to, dest.Addr()))
			} else if dest.Type().AssignableTo(to.Type().Elem()) {
				to.Set(reflect.Append(to, dest))
			}
		}
	}
	return
}

func indirect(reflectValue reflect.Value) reflect.Value {
	for reflectValue.Kind() == reflect.Ptr {
		reflectValue = reflectValue.Elem()
	}
	return reflectValue
}

func indirectType(reflectType reflect.Type) reflect.Type {
	for reflectType.Kind() == reflect.Ptr || reflectType.Kind() == reflect.Slice {
		reflectType = reflectType.Elem()
	}
	return reflectType
}

func set(to, from reflect.Value) bool {
	if from.IsValid() {
		if to.Kind() == reflect.Ptr {
			//set `to` to nil if from is nil
			if from.Kind() == reflect.Ptr && from.IsNil() {
				to.Set(reflect.Zero(to.Type()))
				return true
			} else if to.IsNil() {
				to.Set(reflect.New(to.Type().Elem()))
			}
			to = to.Elem()
		}

		if from.Type().ConvertibleTo(to.Type()) {
			to.Set(from.Convert(to.Type()))
		} else if scanner, ok := to.Addr().Interface().(sql.Scanner); ok {
			err := scanner.Scan(from.Interface())
			if err != nil {
				return false
			}
		} else if from.Kind() == reflect.Ptr {
			return set(to, from.Elem())
		} else {
			return false
		}
	}
	return true
}
