package spentcalories

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"
)

const (
	mInKm                      = 1000
	minInH                     = 60
	stepLengthCoefficient      = 0.45
	walkingCaloriesCoefficient = 0.5
)

func parseTraining(data string) (int, string, time.Duration, error) {
	parts := strings.Split(data, ",")
	if len(parts) != 3 {
		return 0, "", 0, errors.New("неверный формат данных: ожидается 3 значения через запятую")
	}

	stepsStr := strings.TrimSpace(parts[0])
	if stepsStr == "" {
		return 0, "", 0, errors.New("количество шагов не может быть пустым")
	}

	steps, err := strconv.Atoi(stepsStr)
	if err != nil {
		return 0, "", 0, fmt.Errorf("ошибка парсинга шагов: %v", err)
	}
	if steps <= 0 {
		return 0, "", 0, errors.New("количество шагов должно быть положительным")
	}

	activity := strings.TrimSpace(parts[1])
	if activity == "" {
		return 0, "", 0, errors.New("тип активности не может быть пустым")
	}

	durationStr := strings.TrimSpace(parts[2])
	if durationStr == "" {
		return 0, "", 0, errors.New("продолжительность не может быть пустой")
	}

	duration, err := time.ParseDuration(durationStr)
	if err != nil {
		return 0, "", 0, fmt.Errorf("ошибка парсинга продолжительности: %v", err)
	}
	if duration <= 0 {
		return 0, "", 0, errors.New("продолжительность должна быть положительной")
	}

	return steps, activity, duration, nil
}

func distance(steps int, height float64) float64 {
	stepLength := height * stepLengthCoefficient
	return float64(steps) * stepLength / mInKm
}

func meanSpeed(steps int, height float64, duration time.Duration) float64 {
	if duration <= 0 {
		return 0
	}
	return distance(steps, height) / duration.Hours()
}

func TrainingInfo(data string, weight, height float64) (string, error) {
	if weight <= 0 {
		return "", errors.New("вес должен быть положительным")
	}
	if height <= 0 {
		return "", errors.New("рост должен быть положительным")
	}

	steps, activity, duration, err := parseTraining(data)
	if err != nil {
		return "", err
	}

	var calories float64
	distanceKm := distance(steps, height)
	speed := meanSpeed(steps, height, duration)

	switch strings.ToLower(activity) {
	case "ходьба", "walking":
		calories, err = WalkingSpentCalories(steps, weight, height, duration)
	case "бег", "running":
		calories, err = RunningSpentCalories(steps, weight, height, duration)
	default:
		return "", errors.New("неизвестный тип тренировки")
	}

	if err != nil {
		return "", fmt.Errorf("ошибка расчета калорий: %v", err)
	}

	// Форматируем калории в зависимости от значения
	var caloriesStr string
	if calories == float64(int(calories)) {
		caloriesStr = fmt.Sprintf("%.0f", calories)
	} else {
		caloriesStr = fmt.Sprintf("%.2f", calories)
	}

	return fmt.Sprintf(
		"Тип тренировки: %s\nДлительность: %.2f ч.\nДистанция: %.2f км.\nСкорость: %.2f км/ч\nСожгли калорий: %s\n",
		activity,
		duration.Hours(),
		distanceKm,
		speed,
		caloriesStr,
	), nil
}

func RunningSpentCalories(steps int, weight, height float64, duration time.Duration) (float64, error) {
	if steps <= 0 {
		return 0, errors.New("количество шагов должно быть положительным")
	}
	if weight <= 0 {
		return 0, errors.New("вес должен быть положительным")
	}
	if height <= 0 {
		return 0, errors.New("рост должен быть положительным")
	}
	if duration <= 0 {
		return 0, errors.New("продолжительность должна быть положительной")
	}

	speed := meanSpeed(steps, height, duration)
	return speed * weight * duration.Minutes() / minInH, nil
}

func WalkingSpentCalories(steps int, weight, height float64, duration time.Duration) (float64, error) {
	if steps <= 0 {
		return 0, errors.New("количество шагов должно быть положительным")
	}
	if weight <= 0 {
		return 0, errors.New("вес должен быть положительным")
	}
	if height <= 0 {
		return 0, errors.New("рост должен быть положительным")
	}
	if duration <= 0 {
		return 0, errors.New("продолжительность должна быть положительной")
	}

	speed := meanSpeed(steps, height, duration)
	return walkingCaloriesCoefficient * speed * weight * duration.Minutes() / minInH, nil
}
