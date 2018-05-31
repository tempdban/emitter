/**********************************************************************************
* Copyright (c) 2009-2017 Misakai Ltd.
* This program is free software: you can redistribute it and/or modify it under the
* terms of the GNU Affero General Public License as published by the  Free Software
* Foundation, either version 3 of the License, or(at your option) any later version.
*
* This program is distributed  in the hope that it  will be useful, but WITHOUT ANY
* WARRANTY;  without even  the implied warranty of MERCHANTABILITY or FITNESS FOR A
* PARTICULAR PURPOSE.  See the GNU Affero General Public License  for  more details.
*
* You should have  received a copy  of the  GNU Affero General Public License along
* with this program. If not, see<http://www.gnu.org/licenses/>.
************************************************************************************/

package message

import (
	"sort"

	"github.com/golang/snappy"
	"github.com/kelindar/binary"
)

// Frame represents a message frame which is sent through the wire to the
// remote server and contains a set of messages.
type Frame []Message

// Sort sorts the frame
func (f Frame) Sort() {
	sort.Slice(f, func(i, j int) bool { return f[i].Time < f[j].Time })
}

// Limit limits the frame to a specific number of elements
func (f *Frame) Limit(n int) {
	if len(*f) > n {
		*f = (*f)[:n]
	}
}

// Message represents a message which has to be forwarded or stored.
type Message struct {
	Time    int64  `json:"ts,omitempty"`   // The timestamp of the message
	Ssid    Ssid   `json:"ssid,omitempty"` // The Ssid of the message
	Channel []byte `json:"chan,omitempty"` // The channel of the message
	Payload []byte `json:"data,omitempty"` // The payload of the message
	TTL     uint32 `json:"ttl,omitempty"`  // The time-to-live of the message
}

// NewFrame creates a new frame with the specified capacity
func NewFrame(capacity int) Frame {
	return make(Frame, 0, capacity)
}

// Size returns the byte size of the message.
func (m *Message) Size() int64 {
	return int64(len(m.Payload))
}

// Encode encodes the message frame
func (f *Frame) Encode() (out []byte) {
	var enc []byte
	enc, err := binary.Marshal(f)
	if err != nil {
		panic(err) // This should never happen unless there's some terrible bug in the encoder
	}

	return snappy.Encode(out, enc)
}

// Append appends the message to a frame.
func (f *Frame) Append(time int64, ssid Ssid, channel, payload []byte) {
	*f = append(*f, Message{Time: time, Ssid: ssid, Channel: channel, Payload: payload})
}

// DecodeFrame decodes the message frame from the decoder.
func DecodeFrame(buf []byte) (out Frame, err error) {
	// TODO: optimize
	var buffer []byte
	if buf, err = snappy.Decode(buffer, buf); err == nil {
		out = make(Frame, 0, 64)
		err = binary.Unmarshal(buf, &out)
	}
	return
}