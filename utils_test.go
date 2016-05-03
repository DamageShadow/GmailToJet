package main

import "testing"

func TestDifference(t *testing.T) {

	t.Log("Testing differentiator")

	slice1 := []string{"test", "test2", "test3", "test5"}
	slice2 := []string{"test", "test2", "test4"}

	dif1, dif2 := difference(slice1, slice2)

	properDif1 := []string{"test3", "test5"}
	properDif2 := []string{"test4"}

	t.Log(dif1, " and ", dif2)

	if (dif1[0] != properDif1[0] && dif1[1] != properDif1[1]) {
		t.Error("Dif1 is wrong")
	}

	if (dif2[0] != properDif2[0]) {
		t.Error("Dif2 is wrong")
	}

}
