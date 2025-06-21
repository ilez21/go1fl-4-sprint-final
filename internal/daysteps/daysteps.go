package daysteps

import (
	"errors"
	"fmt"
	"github.com/Yandex-Practicum/tracker/internal/spentcalories"
	"log"
	"strconv"
	"strings"
	"time"
)

const (
	// Длина одного шага в метрах
	stepLength = 0.65
	// Количество метров в одном километре
	mInKm = 1000
)

func parsePackage(data string) (int, time.Duration, error) {
	parts := strings.Split(data, ",")
	if len(parts) != 2 {
		return 0, 0, errors.New("неверный формат данных: ожидается 2 значения через запятую")
	}

	// Удаляем все пробелы из строки с шагами
	stepsStr := strings.ReplaceAll(parts[0], " ", "")
	if stepsStr == "" {
		return 0, 0, errors.New("количество шагов не может быть пустым")
	}

	steps, err := strconv.Atoi(stepsStr)
	if err != nil {
		return 0, 0, fmt.Errorf("ошибка парсинга шагов: %v", err)
	}
	if steps <= 0 {
		return 0, 0, errors.New("количество шагов должно быть положительным")
	}

	// Удаляем все пробелы из строки с продолжительностью
	durationStr := strings.ReplaceAll(parts[1], " ", "")
	if durationStr == "" {
		return 0, 0, errors.New("продолжительность не может быть пустой")
	}

	duration, err := time.ParseDuration(durationStr)
	if err != nil {
		return 0, 0, fmt.Errorf("ошибка парсинга продолжительности: %v", err)
	}
	if duration <= 0 {
		return 0, 0, errors.New("продолжительность должна быть положительной")
	}

	return steps, duration, nil
}
func DayActionInfo(data string, weight, height float64) string {
	steps, duration, err := parsePackage(data)
	if err != nil {
		log.Println(err)
		return ""
	}

	distanceMeters := stepLength * float64(steps)
	distanceKm := distanceMeters / mInKm

	calories, err := spentcalories.WalkingSpentCalories(steps, weight, height, duration)
	if err != nil {
		log.Println(err)
		return ""
	}

	return fmt.Sprintf(
		"Количество шагов: %d.\nДистанция составила %.2f км.\nВы сожгли %.2f ккал.\n",
		steps, distanceKm, calories,
	)
}
