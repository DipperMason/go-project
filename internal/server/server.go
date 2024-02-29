package server

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"html/template"
	"net/http"
	"strconv"

	_ "github.com/go-sql-driver/mysql"

	"github.com/DipperMason/go_calculator/internal/agent"
)

var calcAgent = &agent.CalculatorAgent{}

var db *sql.DB

// Функция для установки соединения с базой данных
func SetupDatabase() error {
	// Открываем соединение с базой данных
	var err error
	db, err = sql.Open("mysql", "user:password@tcp(localhost:3306)/database")
	if err != nil {
		return err
	}
	// Проверяем, что соединение установлено успешно
	if err := db.Ping(); err != nil {
		return err
	}
	return nil
}

// Функция для сохранения выражения в базе данных и возврата его уникального идентификатора
func SaveExpression(expression string) (string, error) {
	// Подготавливаем SQL-запрос для вставки данных
	stmt, err := db.Prepare("INSERT INTO expressions (expression) VALUES (?)")
	if err != nil {
		return "", err
	}
	defer stmt.Close()

	// Выполняем SQL-запрос
	res, err := stmt.Exec(expression)
	if err != nil {
		return "", err
	}

	// Получаем уникальный идентификатор (ID) вставленной записи
	id, err := res.LastInsertId()
	if err != nil {
		return "", err
	}

	// Преобразуем ID в строковый формат и возвращаем его
	return strconv.FormatInt(id, 10), nil
}

// StartServer запускает сервер
func StartServer() {
	http.HandleFunc("/", handleRequest)
	port := 8080
	fmt.Printf("Server started on port %d\n", port)
	http.ListenAndServe(fmt.Sprintf(":%d", port), nil)
}

func handleRequest(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		CalculateHandler(w, r)
		return
	}

	// Парсим шаблон
	tmpl, err := template.ParseFiles("../html/template.html")
	if err != nil {
		http.Error(w, "Ошибка при загрузке шаблона", http.StatusInternalServerError)
		return
	}

	// Генерируем HTML на основе шаблона
	err = tmpl.Execute(w, nil)
	if err != nil {
		http.Error(w, "Ошибка при генерации HTML", http.StatusInternalServerError)
		return
	}
}

// CalculateHandler обрабатывает запросы на вычисление
func CalculateHandler(w http.ResponseWriter, r *http.Request) {
	// Проверяем метод запроса
	if r.Method != http.MethodPost {
		http.Error(w, "Метод не поддерживается", http.StatusMethodNotAllowed)
		return
	}

	// Декодируем JSON-тело запроса
	var requestData struct {
		Expression string `json:"expression"`
	}
	if err := json.NewDecoder(r.Body).Decode(&requestData); err != nil {
		http.Error(w, "Неверный формат JSON", http.StatusBadRequest)
		return
	}

	// Выполняем вычисления с использованием агента калькулятора
	result, err := calcAgent.Calculate(requestData.Expression)
	if err != nil {
		http.Error(w, "Ошибка при выполнении вычислений: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Формируем JSON-ответ
	responseData := struct {
		Result float64 `json:"result"`
	}{Result: result}

	// Отправляем ответ
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(responseData); err != nil {
		http.Error(w, "Ошибка при формировании JSON-ответа", http.StatusInternalServerError)
		return
	}
}
