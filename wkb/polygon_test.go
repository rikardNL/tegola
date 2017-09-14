package wkb_test

import (
	"encoding/binary"
	"fmt"
	"testing"

	"github.com/rikardNL/tegola/wkb"
)

func cmpPolygons(p1, p2 wkb.Polygon) (bool, string) {
	if len(p1) != len(p2) {

		return false, fmt.Sprintf(
			"Polygon lengths do not match. "+
				"Number of lines in Polygon1: %v\n"+
				"Number of lines in Polygon2: %v",
			len(p1),
			len(p2),
		)
	}
	for j := range p1 {
		if ok, err := cmpLines(p1[j], p2[j]); !ok {
			return false, fmt.Sprintf("Line %v did not match: %v", j, err)
		}
	}
	return true, ""
}

func TestPolygon(t *testing.T) {
	testcases := TestCases{
		{
			bytes: []byte{
				//01    02    03    04    05    06    07    08
				0x01, 0x00, 0x00, 0x00, // Number of Rings 1
				0x05, 0x00, 0x00, 0x00, // Number of Points 5
				0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x3e, 0x40, // X1 30
				0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x24, 0x40, // Y1 10
				0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x44, 0x40, // X2 40
				0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x44, 0x40, // Y2 40
				0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x34, 0x40, // X3 20
				0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x44, 0x40, // Y3 40
				0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x24, 0x40, // X4 10
				0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x34, 0x40, // Y4 20
				0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x3e, 0x40, // X5 30
				0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x24, 0x40, // Y5 10
			},
			bom: binary.LittleEndian,
			expected: &wkb.Polygon{
				wkb.LineString{
					wkb.NewPoint(30, 10),
					wkb.NewPoint(40, 40),
					wkb.NewPoint(20, 40),
					wkb.NewPoint(10, 20),
					wkb.NewPoint(30, 10),
				},
			},
		},
		{
			bytes: []byte{
				//01    02    03    04    05    06    07    08
				0x02, 0x00, 0x00, 0x00, // Number of Lines (2)
				0x05, 0x00, 0x00, 0x00, // Number of Points in Line1 (5)
				0x00, 0x00, 0x00, 0x00, 0x00, 0x80, 0x41, 0x40, // X1 35
				0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x24, 0x40, // Y1 10
				0x00, 0x00, 0x00, 0x00, 0x00, 0x80, 0x46, 0x40, // X2 45
				0x00, 0x00, 0x00, 0x00, 0x00, 0x80, 0x46, 0x40, // Y2 45
				0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x2e, 0x40, // X3 15
				0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x44, 0x40, // Y3 40
				0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x24, 0x40, // X4 10
				0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x34, 0x40, // Y4 20
				0x00, 0x00, 0x00, 0x00, 0x00, 0x80, 0x41, 0x40, // X5 35
				0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x24, 0x40, // Y5 10
				0x04, 0x00, 0x00, 0x00, // Number of Points in Line2 (4)
				0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x34, 0x40, // X1 20
				0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x3e, 0x40, // Y1 30
				0x00, 0x00, 0x00, 0x00, 0x00, 0x80, 0x41, 0x40, // X2 35
				0x00, 0x00, 0x00, 0x00, 0x00, 0x80, 0x41, 0x40, // Y2 35
				0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x3e, 0x40, // X3 30
				0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x34, 0x40, // Y3 20
				0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x34, 0x40, // X4 20
				0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x3e, 0x40, // Y4 30
			},
			bom: binary.LittleEndian,
			expected: &wkb.Polygon{
				wkb.LineString{
					wkb.NewPoint(35, 10),
					wkb.NewPoint(45, 45),
					wkb.NewPoint(15, 40),
					wkb.NewPoint(10, 20),
					wkb.NewPoint(35, 10),
				},
				wkb.LineString{
					wkb.NewPoint(20, 30),
					wkb.NewPoint(35, 35),
					wkb.NewPoint(30, 20),
					wkb.NewPoint(20, 30),
				},
			},
		},
	}
	testcases.RunTests(t, func(num int, tcase *TestCase) {

		var p, expected wkb.Polygon
		if cexp, ok := tcase.expected.(*wkb.Polygon); !ok {
			t.Errorf("Bad test case %v", num)
			return
		} else {
			expected = *cexp
		}
		if err := p.Decode(tcase.bom, tcase.Reader()); err != nil {
			t.Errorf("Got unexpected error %v", err)
			return
		}
		if ok, err := cmpPolygons(expected, p); !ok {
			t.Errorf("Failed Polygon Test %v: %v", num, err)
		}
	})
}

func TestMultiPolygon(t *testing.T) {
	testcases := TestCases{
		{
			bytes: []byte{
				//01    02    03    04    05    06    07    08
				0x02, 0x00, 0x00, 0x00, // Number of Polygons (2)
				0x01,                   // Byte Encoding Little
				0x03, 0x00, 0x00, 0x00, // Type Polygon1 (3)
				0x01, 0x00, 0x00, 0x00, // Number of Lines (1)
				0x04, 0x00, 0x00, 0x00, // Number of Points (4)
				0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x3e, 0x40, // X1 30
				0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x34, 0x40, // Y1 20
				0x00, 0x00, 0x00, 0x00, 0x00, 0x80, 0x46, 0x40, // X2 45
				0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x44, 0x40, // Y2 40
				0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x24, 0x40, // X3 10
				0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x44, 0x40, // Y3 40
				0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x3e, 0x40, // X4 30
				0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x34, 0x40, // X4 20
				0x01,                   // Byte Encoding Little
				0x03, 0x00, 0x00, 0x00, // Type Polygon2 (3)
				0x01, 0x00, 0x00, 0x00, // Number of Lines (1)
				0x05, 0x00, 0x00, 0x00, // Number of Points (5)
				0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x2e, 0x40, // X1 15
				0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x14, 0x40, // Y1  5
				0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x44, 0x40, // X2 40
				0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x24, 0x40, // Y2 10
				0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x24, 0x40, // X3 10
				0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x34, 0x40, // Y3 20
				0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x14, 0x40, // X4  5
				0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x24, 0x40, // Y4 10
				0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x2e, 0x40, // X5 15
				0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x14, 0x40, // Y5  5
			},
			bom: binary.LittleEndian,
			expected: &wkb.MultiPolygon{
				wkb.Polygon{
					wkb.LineString{
						wkb.NewPoint(30, 20),
						wkb.NewPoint(45, 40),
						wkb.NewPoint(10, 40),
						wkb.NewPoint(30, 20),
					},
				},
				wkb.Polygon{
					wkb.LineString{
						wkb.NewPoint(15, 5),
						wkb.NewPoint(40, 10),
						wkb.NewPoint(10, 20),
						wkb.NewPoint(5, 10),
						wkb.NewPoint(15, 5),
					},
				},
			},
		},
		{
			bytes: []byte{
				//01    02    03    04    05    06    07    08
				0x02, 0x00, 0x00, 0x00, // Number of Polygons (2)
				0x01,                   // Byte order marker little
				0x03, 0x00, 0x00, 0x00, // type Polygon (3)
				0x01, 0x00, 0x00, 0x00, // Number of Lines (1)
				0x04, 0x00, 0x00, 0x00, // Number of Points (4)
				0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x44, 0x40, // X1 40
				0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x44, 0x40, // Y1 40
				0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x34, 0x40, // X2 20
				0x00, 0x00, 0x00, 0x00, 0x00, 0x80, 0x46, 0x40, // Y2 45
				0x00, 0x00, 0x00, 0x00, 0x00, 0x80, 0x46, 0x40, // X3 45
				0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x3e, 0x40, // Y3 30
				0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x44, 0x40, // X4 40
				0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x44, 0x40, // Y4 40
				0x01,                   // Byte order marker little
				0x03, 0x00, 0x00, 0x00, // Type Polygon(3)
				0x02, 0x00, 0x00, 0x00, // Number of Lines(2)
				0x06, 0x00, 0x00, 0x00, // Number of Points(6)
				0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x34, 0x40, // X1 20
				0x00, 0x00, 0x00, 0x00, 0x00, 0x80, 0x41, 0x40, // Y1 35
				0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x24, 0x40, // X2 10
				0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x3e, 0x40, // Y2 30
				0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x24, 0x40, // X3 10
				0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x24, 0x40, // Y3 10
				0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x3e, 0x40, // X4 30
				0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x14, 0x40, // Y4 5
				0x00, 0x00, 0x00, 0x00, 0x00, 0x80, 0x46, 0x40, // X5 45
				0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x34, 0x40, // Y5 20
				0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x34, 0x40, // X6 20
				0x00, 0x00, 0x00, 0x00, 0x00, 0x80, 0x41, 0x40, // Y6 35
				0x04, 0x00, 0x00, 0x00, // Number of Points(4)
				0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x3e, 0x40, // X1 30
				0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x34, 0x40, // Y1 20
				0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x34, 0x40, // X2 20
				0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x2e, 0x40, // Y2 15
				0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x34, 0x40, // X3 20
				0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x39, 0x40, // Y3 25
				0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x3e, 0x40, // X4 30
				0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x34, 0x40, // Y4 20

			},
			bom: binary.LittleEndian,
			expected: &wkb.MultiPolygon{
				wkb.Polygon{
					wkb.LineString{
						wkb.NewPoint(40, 40),
						wkb.NewPoint(20, 45),
						wkb.NewPoint(45, 30),
						wkb.NewPoint(40, 40),
					},
				},
				wkb.Polygon{
					wkb.LineString{
						wkb.NewPoint(20, 35),
						wkb.NewPoint(10, 30),
						wkb.NewPoint(10, 10),
						wkb.NewPoint(30, 5),
						wkb.NewPoint(45, 20),
						wkb.NewPoint(20, 35),
					},
					wkb.LineString{
						wkb.NewPoint(30, 20),
						wkb.NewPoint(20, 15),
						wkb.NewPoint(20, 25),
						wkb.NewPoint(30, 20),
					},
				},
			},
		},
	}
	testcases.RunTests(t, func(num int, tcase *TestCase) {
		var p, expected wkb.MultiPolygon
		if cexp, ok := tcase.expected.(*wkb.MultiPolygon); !ok {
			t.Errorf("Bad test case %v", num)
			return
		} else {
			expected = *cexp
		}
		if err := p.Decode(tcase.bom, tcase.Reader()); err != nil {
			t.Errorf("Got unexpected error %v", err)
			return
		}
		if len(expected) != len(p) {
			t.Errorf("Length of Multipolygon do not match for Test %v", num)
		}
		for i := range expected {
			if ok, err := cmpPolygons(expected[i], p[i]); !ok {
				t.Errorf("Failed Multipolygon test %v for polygon %v: %v", num, i, err)
			}
		}
	})

}
