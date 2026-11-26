package amqp

import (
	"reflect"
	"testing"
)

func TestAMQPRequest_ToMap(t *testing.T) {
	type fields struct {
		EventType string
		Enabled   string
	}
	tests := []struct {
		name   string
		fields fields
		want   map[string]interface{}
	}{
		{"empty map", fields{EventType: "", Enabled: ""}, map[string]interface{}{}},
		{"map with eventType", fields{EventType: "test", Enabled: ""}, map[string]interface{}{"event_type": "test"}},
		{"map with enabled flag", fields{EventType: "", Enabled: "true"}, map[string]interface{}{"enabled": "true"}},
		{"map with enabled flag and eventType", fields{EventType: "test", Enabled: "true"}, map[string]interface{}{"event_type": "test", "enabled": "true"}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			request := AMQPRequest{
				EventType: tt.fields.EventType,
				Enabled:   tt.fields.Enabled,
			}
			if got := request.ToMap(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ToMap() = %v, want %v", got, tt.want)
			}
		})
	}
}
