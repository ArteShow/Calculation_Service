package calculate

import (
	"errors"
	"os"
	"testing"
)

func test_createLogger(t *testing.T){
	_, err := createLogger("./log", "TestLog")
	if err != nil{
		errors.New("Soimething went wrong")
	}
	_, err2 := os.Stat("./log")
	if err2 != nil{
		t.Fatalf("No logger in log")
	}
}

func test_Postpone(t* testing.T){
	example := make([]float64, 5)
	for i := 1; i <= 5; i++{
		example[i-1] = float64(i)
	}
	expected := []float64{1, 2, 4, 5}
	got := Postpone(example, 2)
	for i := 0; i <= 2; i++{
		if expected[0] != got[0]{
			t.Fatalf("Fail in Postpone")
		}
	}
}

func test_PostponeStringSlice(t* testing.T){
	example := []string{"He", "llo", "Worl", "d!"}
	expected := []string{"He", "llo", "d!"}
	got := PostponeStringSlice(example, 2)
	for i := 0; i <= 2; i++{
		if expected[0] != got[0]{
			t.Fatalf("Fail in Postpone")
		}
	}
}

func test_ClacBasic(t*testing.T){
	example := []string{"1+2", "3-2", "5*3", "15/3"}
	expected := []float64{3, 1, 15, 5 }
	for i := 0; i <= 3; i++{
		got, err := CalcBasic(example[i])
		if err != nil{
			panic(err)
		}
		if got != expected[i]{
			t.Fatalf("Fail by testing Calc Basic")
		}
	}
}

func test_Clac(t*testing.T){
	example := []string{"(1+2)*2", "(4-3)*3", "(5*3)-4", "15/(3-2)"}
	expected := []float64{6, 3, 11, 15}
	for i := 0; i <= 3; i++{
		got, err := Calc(example[i])
		if err != nil{
			panic(err)
		}
		if got != expected[i]{
			t.Fatalf("Fail by testing Calc Basic")
		}
	}
}