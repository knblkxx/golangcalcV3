package calculator

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/joho/godotenv"
	_ "github.com/mattn/go-sqlite3"
	"github.com/knblkxx/golangcalcV3/structs"
)

type MoveType struct {
	Type      string
	Index     int
	Prioritet int
}
type ExpressionOutput struct {
	Num   string
	Index int
}
type ExpressionInput struct {
	Move MoveType
	Chan chan ExpressionOutput
}

var TIME_ADDITION_MS int = 0
var TIME_SUBTRACTION_MS int = 0
var TIME_MULTIPLICATIONS_MS int = 0
var TIME_DIVISIONS_MS int = 0

var TaskURL = "http://localhost:8080/internal/task"


func NoSpaces(nums string) string {
	var out []string
	for _, c := range nums {
		if c != ' ' {
			out = append(out, string(c))
		}
	}
	return strings.Join(out, "")
}
func FinalCalc(input ExpressionInput, expression []string) {
	a, _ := strconv.ParseFloat(expression[input.Move.Index-1], 64)
	b, _ := strconv.ParseFloat(expression[input.Move.Index+1], 64)
	var result structs.AgentResult
	url := TaskURL

	if input.Move.Type == "+" {

		data, _ := json.Marshal(structs.AgentResponse{a, b, input.Move.Type, TIME_ADDITION_MS})
		resp, _ := http.Post(url, "application/json", bytes.NewBuffer(data))
		body, _ := ioutil.ReadAll(resp.Body)
		resp.Body.Close()
		json.Unmarshal(body, &result)
	} else if input.Move.Type == "-" {
		data, _ := json.Marshal(structs.AgentResponse{a, b, input.Move.Type, TIME_SUBTRACTION_MS})
		resp, _ := http.Post(url, "application/json", bytes.NewBuffer(data))
		body, _ := ioutil.ReadAll(resp.Body)
		resp.Body.Close()
		json.Unmarshal(body, &result)
	} else if input.Move.Type == "*" {
		data, _ := json.Marshal(structs.AgentResponse{a, b, input.Move.Type, TIME_MULTIPLICATIONS_MS})
		resp, _ := http.Post(url, "application/json", bytes.NewBuffer(data))
		body, _ := ioutil.ReadAll(resp.Body)
		resp.Body.Close()
		json.Unmarshal(body, &result)
	} else if input.Move.Type == "/" {
		if b != 0 {
			data, _ := json.Marshal(structs.AgentResponse{a, b, input.Move.Type, TIME_DIVISIONS_MS})
			resp, _ := http.Post(url, "application/json", bytes.NewBuffer(data))
			body, _ := ioutil.ReadAll(resp.Body)
			resp.Body.Close()
			json.Unmarshal(body, &result)
		} else {
			input.Chan <- ExpressionOutput{
				"error",
				input.Move.Index,
			}
			return
		}
	}
	input.Chan <- ExpressionOutput{
		fmt.Sprintf("%f", result.Result),
		input.Move.Index,
	}
}
func CalcExpression(expression string) (string, error) {
	expression = NoSpaces(expression)
	expression = strings.Replace(expression, "/", " / ", -1)
	expression = strings.Replace(expression, "*", " * ", -1)
	expression = strings.Replace(expression, "+", " + ", -1)
	expression = strings.Replace(expression, "-", " - ", -1)
	nums := strings.Split(expression, " ")
	if strings.Contains("+-/*", nums[0]) {
		return "", errors.New("Невалидное выражение")
	}
	if strings.Contains("+-/*", nums[len(nums)-1]) {
		return "", errors.New("Невалидное выражение")
	}
	if len(nums) == 1 {
		return expression, nil
	}
	//fmt.Println(nums)
	n := 0
	moves := 0
	prioritet := -1
	proshloe := -1
	movesSlice := make([]MoveType, 0)
	for i, el := range nums {
		//fmt.Println(el)
		if strings.Contains("+-/*", el) {
			if proshloe != -1 {
				if proshloe == 1 {
					fmt.Println(nums)
					return "", errors.New("Выражение начианется  с действия")
				}
			}
			moves++
			proshloe = 1
			if prioritet == -1 && strings.Contains("/*", el) {
				prioritet = i
			}
			movesSlice = append(movesSlice, MoveType{
				Type:      el,
				Index:     i,
				Prioritet: prioritet,
			})
		} else {
			for _, c := range el {
				if !strings.Contains("1234567890.", string(c)) {
					return "", errors.New("Невалидные символы")
				}
			}
			if proshloe != -1 {
				if proshloe == 0 {
					return "", errors.New("Действия и цифры в странном порядке")
				}
			}
			proshloe = 0
			n++
		}
	}
	if len(movesSlice) > 1 {
		//fmt.Println(movesSlice)
		parralelMoves := make([]MoveType, 0)
		i := 0
		for {
			if i >= len(movesSlice) {
				break
			}
			if i == 0 {
				if movesSlice[i].Prioritet >= movesSlice[i+1].Prioritet {
					parralelMoves = append(parralelMoves, movesSlice[i])
					i += 2
				} else {
					i++
				}
			} else if i == len(movesSlice)-1 {
				if movesSlice[i].Prioritet >= movesSlice[i-1].Prioritet {
					parralelMoves = append(parralelMoves, movesSlice[i])
					i += 2
				} else {
					i++
				}
			} else {
				if movesSlice[i-1].Prioritet <= movesSlice[i].Prioritet && movesSlice[i].Prioritet >= movesSlice[i+1].Prioritet {
					parralelMoves = append(parralelMoves, movesSlice[i])
					i += 2
				} else {
					i++
				}
			}
		}
		//fmt.Println(parralelMoves)
		ch := make(chan ExpressionOutput)
		for _, mov := range parralelMoves {
			go FinalCalc(ExpressionInput{
				mov,
				ch,
			}, nums)
		}
		for i := 0; i < len(parralelMoves); i++ {
			select {
			case x, ok := <-ch:
				if ok {
					if x.Num == "error" {
						return "", errors.New("Деление на ноль")
					}
					nums[x.Index] = x.Num
				}
			}
		}
		new_nums := make([]string, 0)
		if strings.Contains("+-/*", nums[1]) {
			new_nums = append(new_nums, nums[0])
		}
		for i := 1; i+1 < len(nums); i++ {
			if !strings.Contains("+-/*", nums[i-1]) && !strings.Contains("+-/*", nums[i+1]) {
				new_nums = append(new_nums, nums[i])
			} else if strings.Contains("+-/*", nums[i-1]) && strings.Contains("+-/*", nums[i+1]) {
				new_nums = append(new_nums, nums[i])
			} else if !strings.Contains("+-/*", nums[i]) && !strings.Contains("+-/*", nums[i-1]) && !strings.Contains("+-/*", nums[i+1]) {
				new_nums = append(new_nums, nums[i])
			}
		}
		if strings.Contains("+-/*", nums[len(nums)-2]) {
			new_nums = append(new_nums, nums[len(nums)-1])
		}
		//fmt.Println(nums, new_nums)
		return CalcExpression(strings.Join(new_nums, " "))
	}
	if n-moves != 1 {
		return "", errors.New("Чёт число символов не то")
	}
	out := 0.0
	if prioritet != -1 {
		a, _ := strconv.ParseFloat(nums[prioritet-1], 64)
		b, _ := strconv.ParseFloat(nums[prioritet+1], 64)
		if nums[prioritet] == "*" {
			timer := time.NewTimer(time.Duration(TIME_MULTIPLICATIONS_MS) * time.Millisecond)
			out = a * b
			<-timer.C
		} else {
			if b != 0 {
				timer := time.NewTimer(time.Duration(TIME_DIVISIONS_MS) * time.Millisecond)
				out = a / b
				<-timer.C
			} else {
				return "", errors.New("ДЕЛЕНИЕ НА НОЛЬ ХАХАХХА")
			}
		}
		if len(nums)-2 != 1 {
			return CalcExpression(fmt.Sprintf("%s%f%s", strings.Join(nums[:prioritet-1], ""), out, strings.Join(nums[prioritet+2:], "")))
		}
	} else {
		a, _ := strconv.ParseFloat(nums[0], 64)
		b, _ := strconv.ParseFloat(nums[2], 64)
		if nums[1] == "+" {
			timer := time.NewTimer(time.Duration(TIME_ADDITION_MS) * time.Millisecond)
			out = a + b
			<-timer.C
		} else {
			timer := time.NewTimer(time.Duration(TIME_SUBTRACTION_MS) * time.Millisecond)
			out = a - b
			<-timer.C
		}
		if len(nums)-2 != 1 {
			return CalcExpression(fmt.Sprintf("%f%s", out, strings.Join(nums[3:], "")))
		}
	}
	return fmt.Sprintf("%f", out), nil
}

func Calc(expression string, id int) (float64, error) {
	open := 0
	begin := -1
	end := -1
	conn, err := sql.Open("sqlite3", "database/database.sql")
	if err != nil {
		return 0, err
	}
	defer conn.Close()
	for i, c := range expression {
		if c == '(' {
			open++
			begin = i
		} else if c == ')' {
			open--
			end = i
			if open == -1 {
				_, err = conn.Exec(`
					UPDATE expressions
					SET status = "error", result = 0
					WHERE id = ?
				`, fmt.Sprintf("%d", id))
				if err != nil {
					return 0, err
				}
			}
			if end-begin == 1 {
				_, err = conn.Exec(`
					UPDATE expressions
					SET status = "error", result = 0
					WHERE id = ?
				`, fmt.Sprintf("%d", id))
				if err != nil {
					return 0, err
				}
			}
			res, err := CalcExpression(expression[begin+1 : end])
			if err != nil {
				_, err = conn.Exec(`
					UPDATE expressions
					SET status = "error", result = 0
					WHERE id = ?
				`, fmt.Sprintf("%d", id))
				if err != nil {
					return 0, err
				}
			}
			return Calc(expression[:begin]+res+expression[end+1:], id)
		}
	}

	if open > 0 {
		_, err = conn.Exec(`
			UPDATE expressions
			SET status = "error", result = 0
			WHERE id = ?
		`, fmt.Sprintf("%d", id))
		if err != nil {
			return 0, err
		}
	}
	out, err := CalcExpression(expression)
	if err != nil {
		_, err = conn.Exec(`
			UPDATE expressions
			SET status = "error", result = 0
			WHERE id = ?
		`, fmt.Sprintf("%d", id))
		if err != nil {
			return 0, err
		}
	}
	out1, _ := strconv.ParseFloat(out, 64)
	_, err = conn.Exec(`
			UPDATE expressions
			SET status = "ok", result = ?
			WHERE id = ?
		`, fmt.Sprintf("%f", out1), fmt.Sprintf("%d", id))
	if err != nil {
		return 0, err
	}
	return out1, nil
}

func Initial() {
	godotenv.Load(".env")
	value := os.Getenv("TIME_ADDITION_MS")
	if value != "" {
		intvalue, err := strconv.Atoi(value)
		if err != nil {
			fmt.Println("Ошибка в environment variable TIME_ADDITION_MS")
			os.Exit(0)
		}
		TIME_ADDITION_MS = intvalue
	} else {
		TIME_ADDITION_MS = 0
	}

	value = os.Getenv("TIME_SUBTRACTION_MS")
	if value != "" {
		intvalue, err := strconv.Atoi(value)
		if err != nil {
			fmt.Println("Ошибка в environment variable TIME_SUBTRACTION_MS")
			os.Exit(0)
		}
		TIME_SUBTRACTION_MS = intvalue
	} else {
		TIME_SUBTRACTION_MS = 0
	}

	value = os.Getenv("TIME_MULTIPLICATIONS_MS")
	if value != "" {
		intvalue, err := strconv.Atoi(value)
		if err != nil {
			fmt.Println("Ошибка в environment variable TIME_MULTIPLICATIONS_MS")
			os.Exit(0)
		}
		TIME_MULTIPLICATIONS_MS = intvalue
	} else {
		TIME_MULTIPLICATIONS_MS = 0
	}

	value = os.Getenv("TIME_DIVISIONS_MS")
	if value != "" {
		intvalue, err := strconv.Atoi(value)
		if err != nil {
			fmt.Println("Ошибка в environment variable TIME_DIVISIONS_MS")
			os.Exit(0)
		}
		TIME_DIVISIONS_MS = intvalue
	} else {
		TIME_DIVISIONS_MS = 0
	}
	fmt.Printf("TIME_ADDITION_MS: %d\n", TIME_ADDITION_MS)
	fmt.Printf("TIME_SUBTRACTION_MS: %d\n", TIME_SUBTRACTION_MS)
	fmt.Printf("TIME_MULTIPLICATIONS_MS: %d\n", TIME_MULTIPLICATIONS_MS)
	fmt.Printf("TIME_DIVISIONS_MS: %d\n", TIME_DIVISIONS_MS)

}

func test() {
	fmt.Println(Calc("2 + 2 + 2 + 2 + 2 + 2 + (2 + (2 + (2 + 2)))", 0))
	fmt.Println(Calc("1+1", 0))
	fmt.Println(Calc("(2+2)*2", 0))
	fmt.Println(Calc("2+2*2", 0))
	fmt.Println(Calc("1+1*", 0))
}
