package domain

type Vector []float64

func NewVector(value []float64) Vector {
	copied := make([]float64, len(value))
	copy(copied, value)
	return copied
}

func (v Vector) Clone() Vector {
	if v == nil {
		return nil
	}
	cloned := make([]float64, len(v))
	copy(cloned, v)
	return cloned
}

type Content string

func NewContent(value string) Content {
	return Content(value)
}

func (c Content) String() string {
	return string(c)
}

type Info string

func NewInfo(value string) Info {
	return Info(value)
}

func (i Info) String() string {
	return string(i)
}
