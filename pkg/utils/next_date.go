package utils

import (
	"errors"
	"strconv"
	"strings"
	"time"
)

func NextDate(now time.Time, date time.Time, repeat string) (string, error) {
	parts := strings.Split(strings.TrimSpace(repeat), " ")
	if len(parts) == 0 {
		return "", errors.New("пустое правило повторения")
	}

	rule := parts[0]

	switch rule {
	case "d":
		if len(parts) != 2 {
			return "", errors.New("неверный формат правила d: ожидается d <число>")
		}

		interval, err := strconv.Atoi(parts[1])
		if err != nil {
			return "", errors.New("неверное число дней в правиле d")
		}

		if interval < 1 || interval > 400 {
			return "", errors.New("интервал дней должен быть от 1 до 400")
		}

		date = date.AddDate(0, 0, interval)
		for !AfterNow(date, now) {
			date = date.AddDate(0, 0, interval)
		}

		return date.Format("20060102"), nil

	case "y":
		if len(parts) != 1 {
			return "", errors.New("правило y не требует дополнительных параметров")
		}

		date = date.AddDate(1, 0, 0)
		for !AfterNow(date, now) {
			date = date.AddDate(1, 0, 0)
		}

		return date.Format("20060102"), nil
	default:
		return "", errors.New("неизвестное правило повторения: " + rule)
	}
}

func AfterNow(date, now time.Time) bool {
	dateOnly := time.Date(date.Year(), date.Month(), date.Day(), 0, 0, 0, 0, time.UTC)
	nowOnly := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.UTC)

	return dateOnly.After(nowOnly)
}
