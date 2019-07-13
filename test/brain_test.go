package test

import (
	"github.com/stretchr/testify/suite"
	"testing"
)

func Test(t *testing.T) {
	suite.Run(t, New())
}

/*
//	a := "2733d1c1b714e141"
//	latBytes := a[:8]
//	lonBytes := a[8:]

	src := []byte("2733d1c1")

	dst := make([]byte, hex.DecodedLen(len(src)))
	n, err := hex.Decode(dst, src)
	if err != nil {
		return
	}

	fmt.Printf("%s\n", dst[:n])

	 bits := binary.LittleEndian.Uint32(dst[:n])
	 final := math.Float32frombits(bits)

	fmt.Println(final)

*/
