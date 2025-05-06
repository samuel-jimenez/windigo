/*
 * Copyright (C) 2019 The windigo Authors. All Rights Reserved.
 * Copyright (C) 2010-2013 Allen Dang. All Rights Reserved.
 */

package windigo

type Event struct {
	Sender Controller
	Data   interface{}
}

func NewEvent(sender Controller, data interface{}) *Event {
	return &Event{Sender: sender, Data: data}
}
