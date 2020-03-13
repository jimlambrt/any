package main

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"reflect"

	proto "google.golang.org/protobuf/proto"

	"github.com/jimlambrt/any"
	"github.com/jimlambrt/any/any_test"
)

func main() {
	types, _ := any.NewTypeCatalog(
		new(any_test.TestUser),
		new(any_test.TestCar),
		new(any_test.TestRental),
		new(TestShirt),
		new(TestPants),
	)
	queue := any.Queue{Catalog: types}

	user := any_test.TestUser{
		Id:          1,
		Name:        "Alice",
		PhoneNumber: "867-5309",
		Email:       "alice@bob.com",
	}
	car := any_test.TestCar{
		Id:    2,
		Model: "Jeep",
		Mpg:   25,
	}
	rental := any_test.TestRental{
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

	_ = queue.Add(&user, &car, &rental, &shirt, &pants)

	queuedUser, _ := queue.Remove()
	origUser, _ := comparable(&user)
	printJSON(queuedUser, origUser)
	if !reflect.DeepEqual(*origUser.(*any_test.TestUser), *queuedUser.(*any_test.TestUser)) {
		panic("users should be equal")
	}

	queuedCar, _ := queue.Remove()
	origCar, _ := comparable(&car)
	printJSON(queuedCar, origCar)
	if !reflect.DeepEqual(*origCar.(*any_test.TestCar), *queuedCar.(*any_test.TestCar)) {
		panic("cars should be equal")
	}

	queuedRental, _ := queue.Remove()
	origRental, _ := comparable(&rental)
	printJSON(queuedRental, origRental)
	if !reflect.DeepEqual(*origRental.(*any_test.TestRental), *queuedRental.(*any_test.TestRental)) {
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

	_, err := queue.Remove()
	if err != io.EOF {
		panic("should get EOF")
	}

	wfileQueue := any.Queue{Catalog: types}
	if err := wfileQueue.Add(&user, &car); err != nil {
		panic(err)
	}

	tmpfile, err := ioutil.TempFile("", "example")
	if err != nil {
		log.Fatal(err)
	}

	defer os.Remove(tmpfile.Name()) // clean up

	if _, err := tmpfile.Write(wfileQueue.QueueBuffer); err != nil {
		log.Fatalln("Failed to write address book:", err)
	}
	_ = tmpfile.Close()

	rfileQueue := any.Queue{Catalog: types}
	in, err := ioutil.ReadFile(tmpfile.Name())
	if err != nil {
		log.Fatalln("Error reading file:", err)
	}
	rfileQueue.QueueBuffer = in
	fileUser, _ := rfileQueue.Remove()
	fmt.Println("from file.")
	printJSON(fileUser, origUser)
	if !reflect.DeepEqual(*origUser.(*any_test.TestUser), *fileUser.(*any_test.TestUser)) {
		panic("users should be equal")
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
