package decoder

type Frame struct {
	Payload []byte
}

type Decoder struct{ buffer []byte }

func (d *Decoder) Decode(input []byte) Frame {
	d.buffer = append(d.buffer[:0], input...)
	return Frame{Payload: d.buffer}
}
