package any

import (
	"encoding/binary"
	"encoding/json"
	"fmt"
	"reflect"
	"sync"

	"github.com/golang/protobuf/proto"
	"github.com/golang/protobuf/ptypes/any"
)

// Queue provides a FIFO queue which is concurrent safe
type Queue struct {
	// QueueBuffer provides a buffer for the queue
	QueueBuffer
	// Catalog provides a catalog of types for the queue
	Catalog *TypeCatalog
	l       sync.Mutex
}

// Add something to the queue
func (r *Queue) Add(things ...interface{}) error {
	for _, m := range things {
		_, isProto := m.(proto.Message)
		type_usl := reflect.TypeOf(m).String()

		var value []byte
		var err error
		if isProto {
			value, err = proto.Marshal(m.(proto.Message))
			if err != nil {
				return err
			}
		} else {
			value, err = json.Marshal(m.(interface{}))
			if err != nil {
				return err
			}
		}
		msg := &Any{
			IsPb: isProto,
			Anything: &any.Any{
				TypeUrl: type_usl,
				Value:   value,
			},
		}

		data, err := proto.Marshal(msg)
		if err != nil {
			return err
		}
		r.l.Lock()
		defer r.l.Unlock()
		err = binary.Write(r, binary.LittleEndian, int32(len(data)))
		if err != nil {
			return err
		}
		n, err := r.Write(data)
		if err != nil {
			return err
		}
		if n != len(data) {
			return fmt.Errorf("failed to write the all the data to the buffer: %d of %d", n, len(data))
		}
	}
	return nil
}

// Remove (FIFO) something from the queue and EOF if empty
func (r *Queue) Remove() (interface{}, error) {
	r.l.Lock()
	defer r.l.Unlock()
	var n int32
	err := binary.Read(r, binary.LittleEndian, &n)
	if err != nil {
		return nil, err
	}
	data := r.Next(int(n))
	msg := new(Any)
	err = proto.Unmarshal(data, msg)
	if err != nil {
		return nil, err
	}
	if msg.Anything.Value == nil {
		return nil, nil
	}
	any, err := r.Catalog.Get(msg.Anything.TypeUrl)
	if err != nil {
		return nil, err
	}
	if msg.IsPb {
		pm := any.(proto.Message)
		err = proto.Unmarshal(msg.Anything.Value, pm)
		return pm, err
	}
	pm := any.(interface{})
	err = json.Unmarshal(msg.Anything.Value, pm)
	return pm, err
}
