package any

import (
	"bytes"
	"io"
	"reflect"
	"testing"

	proto "github.com/golang/protobuf/proto"
	"github.com/jimlambrt/any/any_test"
	"github.com/matryer/is"
)

func Test_QueueBuffer(t *testing.T) {
	t.Parallel()
	is := is.New(t)
	b := QueueBuffer{}

	_, err := b.Write([]byte("bye"))
	is.NoErr(err)
	is.True(b.Len() == len([]byte("bye")))

	_, err = b.Write([]byte("hello"))
	is.NoErr(err)
	is.True(b.Len() == len([]byte("bye"))+len([]byte("hello")))

	bye := make([]byte, 3)
	_, err = b.Read(bye)
	is.NoErr(err)
	is.True(bytes.Equal(bye, []byte("bye")))

	hello := b.Next(len([]byte("hello")))
	is.NoErr(err)
	is.True(bytes.Equal(hello, []byte("hello")))

	t.Log(string(hello))
	t.Log(string(bye))
}

func Test_Queue(t *testing.T) {
	t.Parallel()
	is := is.New(t)
	types, err := NewTypeCatalog(
		new(any_test.TestUser),
		new(any_test.TestCar),
		new(any_test.TestRental),
		new(TestShirt),
		new(TestPants),
	)
	is.NoErr(err)
	queue := Queue{Catalog: types}

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

	err = queue.Add(&user)
	is.NoErr(err)
	err = queue.Add(&car)
	is.NoErr(err)
	err = queue.Add(&rental)
	is.NoErr(err)

	err = queue.Add(&shirt)
	is.NoErr(err)
	err = queue.Add(&pants)
	is.NoErr(err)

	queuedUser, err := queue.Remove()
	is.NoErr(err)
	origUser, err := comparable(&user)
	is.NoErr(err)
	is.True(reflect.DeepEqual(*origUser.(*any_test.TestUser), *queuedUser.(*any_test.TestUser)))

	queuedCar, err := queue.Remove()
	is.NoErr(err)
	origCar, err := comparable(&car)
	is.NoErr(err)
	is.True(reflect.DeepEqual(*origCar.(*any_test.TestCar), *queuedCar.(*any_test.TestCar)))

	queuedRental, err := queue.Remove()
	is.NoErr(err)
	origRental, err := comparable(&rental)
	is.NoErr(err)
	is.True(reflect.DeepEqual(*origRental.(*any_test.TestRental), *queuedRental.(*any_test.TestRental)))

	queuedShirt, err := queue.Remove()
	is.NoErr(err)
	is.True(reflect.DeepEqual(&shirt, queuedShirt.(*TestShirt)))

	queuedPants, err := queue.Remove()
	is.NoErr(err)
	is.True(reflect.DeepEqual(&pants, queuedPants.(*TestPants)))

	_, err = queue.Remove()
	is.True(err == io.EOF)
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
