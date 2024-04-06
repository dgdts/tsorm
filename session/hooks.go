package session

import (
	"reflect"
	"tsorm/log"
)

// Constants representing different lifecycle events.
const (
	BeforeQuery  = "BeforeQuery"
	AfterQuery   = "AfterQuery"
	BeforeUpdate = "BeforeUpdate"
	AfterUpdate  = "AfterUpdate"
	BeforeDelete = "BeforeDelete"
	AfterDelete  = "AfterDelete"
	BeforeInsert = "BeforeInsert"
	AfterInsert  = "AfterInsert"
)

// CallMethod calls the specified method on the value using reflection.
// It accepts the method name and the value on which the method should be called.
func (s *Session) CallMethod(method string, value interface{}) {
	// Get the method of the model associated with the session.
	fm := reflect.ValueOf(s.RefTable().Model).MethodByName(method)
	// If a specific value is provided, get the method from that value instead.
	if value != nil {
		fm = reflect.ValueOf(value).MethodByName(method)
	}

	// Prepare the parameters for the method call.
	param := []reflect.Value{reflect.ValueOf(s)}

	// Check if the method is valid.
	if fm.IsValid() {
		// Call the method with the session as the parameter.
		if v := fm.Call(param); len(v) > 0 {
			// If the method returns an error, log it.
			if err, ok := v[0].Interface().(error); ok {
				log.Error(err)
			}
		}
	}
}
