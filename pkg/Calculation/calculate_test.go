package calculate

import (
	"errors"
	"os"
	"testing"
)

func TestCreateLogger(t *testing.T){
	_, err, _ := createLogger("./log", "TestLog")
	if err != nil{
		errors.New("Something went wrong")
	}
	_, err2 := os.Stat("./log")
	if err2 != nil{
		t.Fatalf("No logger in log")
	}
}

func TestPostpone(t* testing.T){
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

func TestPostponeStringSlice(t* testing.T){
	example := []string{"He", "llo", "Worl", "d!"}
	expected := []string{"He", "llo", "d!"}
	got := PostponeStringSlice(example, 2)
	for i := 0; i <= 2; i++{
		if expected[0] != got[0]{
			t.Fatalf("Fail in Postpone")
		}
	}
}

func TestClacBasic(t*testing.T){
	example := []string{"1+2", "3-2", "5*3", "15/3"}
	expected := []float64{3, 1, 15, 5 }
	for i := 0; i <= 3; i++{
		got, err, _ := CalcBasic(example[i])
		if err != nil{
			panic(err)
		}
		if got != expected[i]{
			t.Fatalf("Fail by testing Calc Basic")
		}
	}
}

func TestClac(t*testing.T){
	example := []string{"(1+2)*2", "(4-3)*3", "(5*3)-4", "15/(3-2)"}
	expected := []float64{6, 3, 11, 15}
	for i := 0; i <= 3; i++{
		got, err, _ := Calc(example[i])
		if err != nil{
			panic(err)
		}
		if got != expected[i]{
			t.Fatalf("Fail by testing Calc Basic")
		}
	}
}