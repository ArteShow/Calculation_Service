package calculate

import (
	"errors"
	"log"
	"os"
	"strconv"
	"strings"
	"unicode"
)

func createLogger(folderPath, fileName string) (*log.Logger, error) {
	err := os.MkdirAll(folderPath, os.ModePerm)
	if err != nil {
		return nil, errors.New("500")
	}

	filePath := folderPath + "/" + fileName

	file, err := os.OpenFile(filePath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		return nil, errors.New("500")
	}

	logger := log.New(file, "", log.LstdFlags)
	return logger, nil
}


func Postpone(nums []float64, index int) []float64 {
	return append(nums[:index], nums[index+1:]...)
}

func PostponeStringSlice(slice []string, index int) []string {
	return append(slice[:index], slice[index+1:]...)
}

func CalcBasic(expression string) (float64, error) {
	Logger, err := createLogger("./log", "CaclBasicLog.txt")
	if err != nil{
		Logger.Println("[ERROR]:", err)
		return 0.0, errors.New("500")
	}
	Logger.Println("Calculation started")
	var RuneWithout string
	MathOperators := make([]string, 0)
	Numbers := make([]float64, 0)
	var num string

	expression = strings.ReplaceAll(expression, ",", ".")
	for _, letter := range expression {
		if letter == ' ' {
			continue
		}
		RuneWithout += string(letter)
	}
	Logger.Println("Replaced all , to . and backspace to nothing")
	String := string(RuneWithout)
	if len(String) == 0 {
		Logger.Println("[ERROR]: error by converting")
		return 0.0, errors.New("422")
	}

	for _, letter := range String {
		if unicode.IsDigit(letter) || letter == '.' {
			num += string(letter)
		} else {
			if num != "" {
				numb, err := strconv.ParseFloat(num, 64)
				if err != nil {
					return 0.0, err
				}
				Numbers = append(Numbers, numb)
				num = ""
			}
		}
	}
	if num != "" {
		numb, err := strconv.ParseFloat(num, 64)
		if err != nil {
			return 0.0, err
		}
		Numbers = append(Numbers, numb)
	}
	Logger.Println("Did the clean up of the expresion")
	var lastOperator string
	for _, letter := range String {
		newLetter := string(letter)
		if newLetter == "+" || newLetter == "-" || newLetter == "*" || newLetter == "/" {
			if lastOperator == newLetter {
				return 0.0, errors.New("244")
				Logger.Println("[ERROR]: error by checking the mathmarkups")

			}
			lastOperator = newLetter
			MathOperators = append(MathOperators, newLetter)
		}
	}
	Logger.Println("Count all of the mathsmarksup.")
	for i := 0; i < len(MathOperators); i++ {
		if MathOperators[i] == "*" || MathOperators[i] == "/" {
			if MathOperators[i] == "*" {
				if i+1 < len(Numbers) {
					Numbers[i] = Numbers[i] * Numbers[i+1]
				} else {
					Logger.Println("[ERROR]: error by multiply")
					return 0.0, errors.New("244")
				}
			} else if MathOperators[i] == "/" {
				if Numbers[i+1] == 0 {
					Logger.Println("[ERROR]: error divide by zero")
					return 0.0, errors.New("244")
				}
				if i+1 < len(Numbers) {
					Numbers[i] = Numbers[i] / Numbers[i+1]
				} else {
					Logger.Println("[ERROR]: error by divide")
					return 0.0, errors.New("244")
				}
			}
			Numbers = Postpone(Numbers, i+1)
			MathOperators = PostponeStringSlice(MathOperators, i)
			i--
		}
	}
	Logger.Println("Did the / and *")
	for i := 0; i < len(MathOperators); i++ {
		if MathOperators[i] == "-" {
			Numbers[i] = Numbers[i] - Numbers[i+1]
			Numbers = Postpone(Numbers, i+1)
			MathOperators = PostponeStringSlice(MathOperators, i)
			i--
		} else if MathOperators[i] == "+" {
			Numbers[i] = Numbers[i] + Numbers[i+1]
			Numbers = Postpone(Numbers, i+1)
			MathOperators = PostponeStringSlice(MathOperators, i)
			i--
		}
	}
	Logger.Println("Did the + and -")
	if len(Numbers) == 0 {
		Logger.Println("[ERROR]: error by plus and minus")
		return 0.0, errors.New("244")
	}
	Logger.Println("End of the function CalcBasic")
	return Numbers[0], nil
}

func Calc(expression string) (float64, error) {
	Logger, err := createLogger("./log", "CalcLog.txt")
	if err != nil{
		Logger.Println("[ERROR]:", err)
		return 0.0, errors.New("500")
	}
	Logger.Println("Check for ( and )")
	expression = strings.ReplaceAll(expression, ",", ".")
	var ResultString string
	stack := []string{}
	for _, letter := range expression {
		if letter == ' ' {
			continue
		}
		ResultString += string(letter)
	}
	Logger.Println("Replace all , to .")
	ResultString = "(" + ResultString + ")"
	for i := 0; i < len(ResultString); i++ {
		if ResultString[i] == '(' {
			stack = append(stack, ResultString[i:i+1])
		} else if ResultString[i] == ')' {
			if len(stack) == 0 {
				Logger.Println("[ERROR]: error by ( and )")
				return 0.0, errors.New("500")
			}
			Logger.Println("Check for wrong ( )")
			stack = stack[:len(stack)-1]
			startIndex := strings.LastIndex(ResultString[:i], "(")
			innerExpression := ResultString[startIndex+1 : i]
			Result, err := CalcBasic(innerExpression)
			if err != nil {
				return 0.0, err
			}
			ResultString = ResultString[:startIndex] + strconv.FormatFloat(Result, 'f', 2, 64) + ResultString[i+1:]
			i = startIndex + len(strconv.FormatFloat(Result, 'f', 2, 64)) - 1
		}
	}
	Logger.Println("Calculated it right")
	if len(stack) != 0 {
		Logger.Println("[ERROR]: to many ( )")
		return 0.0, errors.New("500")
	}

	result, err := CalcBasic(ResultString)
	if err != nil {
		Logger.Println("[ERROR]: error in CalcBasic")
		return 0.0, err
	}
	Logger.Println("End of the function Calc")
	return result, nil
}