package service

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/google/uuid"
)

type MathCaptcha struct {
	Question string
	Answer   int
}

var mathCaptchaStore = map[string]MathCaptcha{}

func GenerateMathCaptcha() (string, string) {
	rand.Seed(time.Now().UnixNano())

	a := rand.Intn(20) + 1
	b := rand.Intn(20) + 1

	id := uuid.New().String()

	mathCaptchaStore[id] = MathCaptcha{
		Question: fmt.Sprintf("%d + %d", a, b),
		Answer:   a + b,
	}

	return id, mathCaptchaStore[id].Question
}

func GetMathCaptchaQuestion(id string) string {
	if c, ok := mathCaptchaStore[id]; ok {
		return c.Question
	}
	return ""
}

func VerifyMathCaptcha(id string, answer int) bool {
	c, ok := mathCaptchaStore[id]
	if !ok {
		return false
	}
	delete(mathCaptchaStore, id)
	return c.Answer == answer
}
