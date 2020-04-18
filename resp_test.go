package redisprotocol

import (
	"fmt"
	"testing"
)

func Test(t *testing.T) {
	rbs1 := "$-1\r\n"
	rbs2 := "$5\r\nHello\r\n"

	rbs1Obj, _, _ := DecodeBulkStr(rbs1)
	s1, _ := Encode(rbs1Obj)
	if s1 != rbs1 {
		t.Errorf("Error in BulkStr")
	}

	rbs2Obj, _, _ := DecodeBulkStr(rbs2)
	s2, _ := Encode(rbs2Obj)
	if s2 != rbs2 {
		t.Errorf("Error in BulkStr")
	}
	ris := ":5\r\n"

	riObj, _, _ := DecodeInt(ris)
	s3, _ := Encode(riObj)
	if s3 != ris {
		t.Errorf("Error in Decode Int")
	}

	ra1 := "*-1\r\n"
	ra2 := "*0\r\n"
	ra3 := "*2\r\n*3\r\n:1\r\n:2\r\n:3\r\n\r\n*5\r\n:5\r\n:7\r\n+Hello Word\r\n-Err\r\n$6\r\nfoobar\r\n\r\n\r\n"

	ra1Obj, err, _ := DecodeArray(ra1)
	if err != nil {
		t.Errorf("Failed to decode array1")
	}
	ra1S, err := Encode(ra1Obj)
	if err != nil {
		t.Errorf("Failed to encode array1")
	}
	if ra1 != ra1S {
		t.Errorf("Error in Array 1 Cycle")
	}

	ra2Obj, err, _ := DecodeArray(ra2)
	if err != nil {
		t.Errorf("Failed to decode array2")
	}
	ra2S, err := Encode(ra2Obj)
	if err != nil {
		t.Errorf("Failed to encode array2")
	}
	if ra2 != ra2S {
		t.Errorf("Error in Array 2 Cycle")
	}

	ra3Obj, err, _ := DecodeArray(ra3)
	if err != nil {
		t.Errorf("Failed to decode array3 %v", err)
	}
	ra3S, err := Encode(ra3Obj)
	if err != nil {
		t.Errorf("Failed to encode array3 %v", err)
	}
	if ra3 != ra3S {
		fmt.Printf("%v", ra3Obj)
		fmt.Println("=================================================================")
		fmt.Println(ra3)
		fmt.Println("=================================================================")
		fmt.Println("=================================================================")
		fmt.Println(ra3S)
		fmt.Println("=================================================================")
		t.Errorf("Error in Array 3 Cycle")
	}
}
