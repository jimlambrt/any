package main

import (
	"encoding/json"
	"fmt"
	"io"
	"reflect"

	"github.com/gogo/protobuf/proto"
	"github.com/jimlambrt/any"
)

func main() {
	types, _ := any.NewTypeCatalog(
		new(any.TestUser),
		new(any.TestCar),
		new(any.TestRental),
		new(TestShirt),
		new(TestPants),
	)
	queue := any.Queue{Catalog: types}

	user := any.TestUser{
		Id:          1,
		Name:        "Alice",
		PhoneNumber: "867-5309",
		Email:       "alice@bob.com",
	}
	car := any.TestCar{
		Id:    2,
		Model: "Jeep",
		Mpg:   25,
	}
	rental := any.TestRental{
		UserId: 1,
		CarId:  2,
	}
	shirt := TestShirt{
		Name:         "plain-white-simple",
		NeckSize:     16.25,
		Material:     "cotton",
		SleeveLength: 34,
		Pocket:       false,
	}
	pants := TestPants{
		Name:     "jeans",
		Waist:    32,
		Inseam:   36,
		Material: "denim",
	}

	queue.Add(&user, &car, &rental, &shirt, &pants)

	queuedUser, _ := queue.Remove()
	origUser, _ := comparable(&user)
	printJSON(queuedUser, origUser)
	if !reflect.DeepEqual(*origUser.(*any.TestUser), *queuedUser.(*any.TestUser)) {
		panic("users should be equal")
	}

	queuedCar, _ := queue.Remove()
	origCar, err := comparable(&car)
	printJSON(queuedCar, origCar)
	if !reflect.DeepEqual(*origCar.(*any.TestCar), *queuedCar.(*any.TestCar)) {
		panic("cars should be equal")
	}

	queuedRental, _ := queue.Remove()
	origRental, _ := comparable(&rental)
	printJSON(queuedRental, origRental)
	if !reflect.DeepEqual(*origRental.(*any.TestRental), *queuedRental.(*any.TestRental)) {
		panic("rentals should be equal")
	}

	queuedShirt, _ := queue.Remove()
	printJSON(queuedShirt, shirt)
	if !reflect.DeepEqual(&shirt, queuedShirt.(*TestShirt)) {
		panic("shirts should be equal")
	}

	queuedPants, _ := queue.Remove()
	printJSON(queuedPants, pants)
	if !reflect.DeepEqual(&pants, queuedPants.(*TestPants)) {
		panic("pants should be equal")
	}

	_, err = queue.Remove()
	if err != io.EOF {
		panic("should get EOF")
	}

}

type TestShirt struct {
	Name         string
	Material     string
	Pocket       bool
	SleeveLength int
	NeckSize     float32
}

type TestPants struct {
	Name     string
	Waist    int
	Inseam   int
	Material string
}

func printJSON(things ...interface{}) {
	for _, t := range things {
		j, _ := json.MarshalIndent(t, "", "\t")
		fmt.Println(string(j))
	}
}
func comparable(m proto.Message) (proto.Message, error) {
	data, err := proto.Marshal(m)
	if err != nil {
		return nil, err
	}

	retMsg := reflect.New(reflect.TypeOf(m).Elem()).Interface().(proto.Message)
	if err := proto.Unmarshal(data, retMsg); err != nil {
		return nil, err
	}
	return retMsg, nil
}
