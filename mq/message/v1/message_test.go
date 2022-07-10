package v1

import (
	"testing"

	"google.golang.org/protobuf/runtime/protoimpl"
)

func TestMessage_generateValidSum(t *testing.T) {
	type fields struct {
		state         protoimpl.MessageState
		sizeCache     protoimpl.SizeCache
		unknownFields protoimpl.UnknownFields
		Data          []byte
		Md            map[string]string
		Error         string
		Key           string
		ValidSum      uint32
	}
	tests := []struct {
		name   string
		fields fields
		want   uint32
	}{
		{
			name: "1",
			fields: fields{
				Data:     []byte("hello world"),
				Md:       nil,
				Error:    "",
				Key:      "123",
				ValidSum: 123,
			},
			want: 1831907993,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			x := &Message{
				state:         tt.fields.state,
				sizeCache:     tt.fields.sizeCache,
				unknownFields: tt.fields.unknownFields,
				Data:          tt.fields.Data,
				Md:            tt.fields.Md,
				Error:         tt.fields.Error,
				Key:           tt.fields.Key,
				ValidSum:      tt.fields.ValidSum,
			}
			if got := x.generateValidSum(); got != tt.want {
				t.Errorf("generateValidSum() = %v, want %v", got, tt.want)
			}
			t.Log(tt.fields.ValidSum)
			if x.ValidSum != tt.fields.ValidSum {
				t.Errorf("valid sum should not be changed")
			}
		})
	}
}
