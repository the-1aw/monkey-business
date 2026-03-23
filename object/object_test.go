package object

import "testing"

func TestStringHashKey(t *testing.T) {
	const testValue1 = "Hello World"
	const testValue2 = "My name is hohnny"
	hello1 := &String{Value: testValue1}
	hello2 := &String{Value: testValue1}
	diff1 := &String{Value: testValue2}
	diff2 := &String{Value: testValue2}

	if hello1.HashKey() != hello2.HashKey() {
		t.Errorf("strings with same content have different hash keys")
	}

	if diff1.HashKey() != diff2.HashKey() {
		t.Errorf("strings with same content have different hash keys")
	}

	if hello1.HashKey() == diff1.HashKey() {
		t.Errorf("strings with different content have same hash keys")
	}
}

func TestBooleanHashKey(t *testing.T) {
	true1 := &Boolean{Value: true}
	true2 := &Boolean{Value: true}
	diff1 := &Boolean{Value: false}
	diff2 := &Boolean{Value: false}

	if true1.HashKey() != true2.HashKey() {
		t.Errorf("bool with same content have different hash keys")
	}

	if diff1.HashKey() != diff2.HashKey() {
		t.Errorf("bool with same content have different hash keys")
	}

	if true1.HashKey() == diff1.HashKey() {
		t.Errorf("bool with different content have same hash keys")
	}
}
func TestIntegerHashKey(t *testing.T) {
	one1 := &Integer{Value: 1}
	one2 := &Integer{Value: 1}
	two1 := &Integer{Value: 2}
	two2 := &Integer{Value: 2}

	if one1.HashKey() != one2.HashKey() {
		t.Errorf("bool with same content have different hash keys")
	}

	if two1.HashKey() != two2.HashKey() {
		t.Errorf("bool with same content have different hash keys")
	}

	if one1.HashKey() == two1.HashKey() {
		t.Errorf("bool with different content have same hash keys")
	}
}
