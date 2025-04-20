package spentcalories

import (
	"fmt"
	"strconv"
	"strings"
	"time"
)

// Основные константы, необходимые для расчетов.
const (
	lenStep                    = 0.65 // средняя длина шага.
	mInKm                      = 1000 // количество метров в километре.
	minInH                     = 60   // количество минут в часе.
	stepLengthCoefficient      = 0.45 // коэффициент для расчета длины шага на основе роста.
	walkingCaloriesCoefficient = 0.5  // коэффициент для расчета калорий при ходьбе
)

func parseTraining(data string) (int, string, time.Duration, error) {
	// Разделяем строку на части
	parts := strings.Split(data, ",")
	if len(parts) != 3 {
		return 0, "", 0, fmt.Errorf("invalid data format")
	}

	// Преобразуем количество шагов в int
	steps, err := strconv.Atoi(parts[0])
	if err != nil || steps <= 0 {
		return 0, "", 0, fmt.Errorf("invalid steps value")
	}

	activity := parts[1]

	duration, err := time.ParseDuration(strings.ReplaceAll(parts[2], "h", "h") + "m0s")
	if err != nil {
		return 0, "", 0, fmt.Errorf("invalid duration format")
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
	dist := distance(steps, height)
	hours := duration.Hours()
	return dist / hours
}

func TrainingInfo(data string, weight, height float64) (string, error) {
	steps, activity, duration, err := parseTraining(data)
	if err != nil {
		return "", err
	}

	var calories float64
	var dist float64
	var speed float64

	switch activity {
	case "Бег":
		calories, err = RunningSpentCalories(steps, weight, duration)
		if err != nil {
			return "", err
		}
		dist = distance(steps, height)
		speed = meanSpeed(steps, height, duration)
	case "Ходьба":
		calories, err = WalkingSpentCalories(steps, weight, height, duration)
		if err != nil {
			return "", err
		}
		dist = distance(steps, height)
		speed = meanSpeed(steps, height, duration)
	default:
		return "", fmt.Errorf("неизвестный тип тренировки")
	}

	return fmt.Sprintf("Тип тренировки: %s\nДлительность: %.2f ч.\nДистанция: %.2f км.\nСкорость: %.2f км/ч\nСожгли калорий: %.2f",
		activity, duration.Hours(), dist, speed, calories), nil
}

func RunningSpentCalories(steps int, weight float64, duration time.Duration) (float64, error) {
	if steps <= 0 || weight <= 0 || duration <= 0 {
		return 0, fmt.Errorf("invalid input parameters")
	}

	meanSpd := meanSpeed(steps, 0, duration)
	durationMinutes := duration.Minutes()
	calories := (weight * meanSpd * durationMinutes) / minInH
	return calories, nil
}

func WalkingSpentCalories(steps int, weight, height float64, duration time.Duration) (float64, error) {
	if steps <= 0 || weight <= 0 || height <= 0 || duration <= 0 {
		return 0, fmt.Errorf("invalid input parameters")
	}

	meanSpd := meanSpeed(steps, height, duration)
	durationMinutes := duration.Minutes()
	calories := (weight * meanSpd * durationMinutes) / minInH * walkingCaloriesCoefficient
	return calories, nil
}
