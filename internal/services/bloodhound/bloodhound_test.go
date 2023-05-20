package bloodhound_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/gerladeno/chat-service/internal/services/bloodhound"
)

func TestSearch(t *testing.T) {
	cases := []struct {
		in      string
		verdict bloodhound.Verdict
		facts   int
	}{
		// "Clean".
		{
			in:      "Здравствуйте! Не могу зайти в мобильное приложение, что делать?",
			verdict: bloodhound.VerdictOK,
			facts:   0,
		},
		{
			in:      "Нет, номер не менял, всё тот же 880005553535",
			verdict: bloodhound.VerdictOK,
			facts:   0,
		},
		{
			in:      "Родился я 29 ноября 1996, но какое это имеет значение?",
			verdict: bloodhound.VerdictOK,
			facts:   0,
		},

		// "Dirty".
		{
			in: `Скидываю данные карты, как вы и просили:
2200 1234 5678 9010
12/30
961`,
			verdict: bloodhound.NewVerdict(
				bloodhound.FactContainsCardNumber,
				bloodhound.FactContainsCardCVC,
				bloodhound.FactContainsCardDate,
			),
			facts: 3,
		},
		{
			in: `Какой у меня был порядок действий:
- я ввела данные карты и трёхзначный код 769, как просил сайт;
- потом мне пришла смс с номером 10-8929;
- я ввела код и сайт завис;`,
			verdict: bloodhound.NewVerdict(
				bloodhound.FactContainsCardCVC,
				bloodhound.FactContainsSMSCode,
			),
			facts: 2,
		},
		{
			in: `Подскажите, пожалуйста, меня почему-то перестало пускать по моему паролю bwT4qR.
В чём может быть дело? Я уже оставлял заявку, её номер 22-1911`,
			verdict: bloodhound.NewVerdict(bloodhound.FactContainsSuspiciousPhrases),
			facts:   1,
		},
		{
			in: `Мой логин такой же как на карте 5179-4279-0625-0126 (IVAN KUZYAKIN) – ivankuzyakin.
Кстати у карты скоро истечёт срок действия, она до 07/22, может быть в этом проблема?`,
			verdict: bloodhound.NewVerdict(
				bloodhound.FactContainsCardNumber,
				bloodhound.FactContainsCardDate,
				bloodhound.FactContainsSuspiciousPhrases,
			),
			facts: 3,
		},
		{
			in:      "Деньги я пересылал с карты 4279012606251579, прошло уже 3 дня, а они ещё не пришли.",
			verdict: bloodhound.NewVerdict(bloodhound.FactContainsCardNumber),
			facts:   1,
		},
	}

	for _, tt := range cases {
		t.Run("", func(t *testing.T) {
			v, err := bloodhound.Search(tt.in)
			require.NoError(t, err)

			if !assert.True(t, tt.verdict.Equals(v)) {
				t.Logf("input: `%v`\nverdict: %v", tt.in, v)
			}
			assert.Equal(t, tt.facts, v.FactsNumber())
		})
	}
}
