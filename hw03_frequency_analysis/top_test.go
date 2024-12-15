package hw03frequencyanalysis

import (
	"testing"

	"github.com/stretchr/testify/require"
)

// Change to true if needed.
var taskWithAsteriskIsCompleted = false

var text = `Как видите, он  спускается  по  лестнице  вслед  за  своим
	другом   Кристофером   Робином,   головой   вниз,  пересчитывая
	ступеньки собственным затылком:  бум-бум-бум.  Другого  способа
	сходить  с  лестницы  он  пока  не  знает.  Иногда ему, правда,
		кажется, что можно бы найти какой-то другой способ, если бы  он
	только   мог   на  минутку  перестать  бумкать  и  как  следует
	сосредоточиться. Но увы - сосредоточиться-то ему и некогда.
		Как бы то ни было, вот он уже спустился  и  готов  с  вами
	познакомиться.
	- Винни-Пух. Очень приятно!
		Вас,  вероятно,  удивляет, почему его так странно зовут, а
	если вы знаете английский, то вы удивитесь еще больше.
		Это необыкновенное имя подарил ему Кристофер  Робин.  Надо
	вам  сказать,  что  когда-то Кристофер Робин был знаком с одним
	лебедем на пруду, которого он звал Пухом. Для лебедя  это  было
	очень   подходящее  имя,  потому  что  если  ты  зовешь  лебедя
	громко: "Пу-ух! Пу-ух!"- а он  не  откликается,  то  ты  всегда
	можешь  сделать вид, что ты просто понарошку стрелял; а если ты
	звал его тихо, то все подумают, что ты  просто  подул  себе  на
	нос.  Лебедь  потом  куда-то делся, а имя осталось, и Кристофер
	Робин решил отдать его своему медвежонку, чтобы оно не  пропало
	зря.
		А  Винни - так звали самую лучшую, самую добрую медведицу
	в  зоологическом  саду,  которую  очень-очень  любил  Кристофер
	Робин.  А  она  очень-очень  любила  его. Ее ли назвали Винни в
	честь Пуха, или Пуха назвали в ее честь - теперь уже никто  не
	знает,  даже папа Кристофера Робина. Когда-то он знал, а теперь
	забыл.
		Словом, теперь мишку зовут Винни-Пух, и вы знаете почему.
		Иногда Винни-Пух любит вечерком во что-нибудь поиграть,  а
	иногда,  особенно  когда  папа  дома,  он больше любит тихонько
	посидеть у огня и послушать какую-нибудь интересную сказку.
		В этот вечер...`

func TestTop10(t *testing.T) {
	t.Run("no words in empty string", func(t *testing.T) {
		require.Len(t, Top10(""), 0)
	})

	t.Run("single word", func(t *testing.T) {
		text := "Жили"
		expected := []string{text}
		require.Equal(t, expected, Top10(text))
	})

	t.Run("same word multiple times", func(t *testing.T) {
		text := "опа опа опа"
		expected := []string{"опа"}
		require.Equal(t, expected, Top10(text))
	})

	t.Run("multiple spaces between words", func(t *testing.T) {
		input := "слово1    слово2     слово3      слово1"
		expected := []string{"слово1", "слово2", "слово3"}
		result := Top10(input)
		require.Equal(t, expected[:len(result)], result)
	})

	t.Run("case sensitive words", func(t *testing.T) {
		input := "Слово слово СЛОВО СлоВо слово"
		expected := []string{"слово", "СЛОВО", "СлоВо", "Слово"}
		result := Top10(input)
		require.Equal(t, expected[:len(result)], result)
	})

	t.Run("words with punctuation", func(t *testing.T) {
		input := "слово1, слово3! слово3? слово1... слово2,"
		expected := []string{"слово1,", "слово1...", "слово2,", "слово3!", "слово3?"}
		result := Top10(input)
		require.Equal(t, expected[:len(result)], result)
	})

	t.Run("more than 10 different words", func(t *testing.T) {
		input := "а я вовсе не колдунья, я любила и люблю, это мне судьба послала..."
		result := Top10(input)
		require.Len(t, result, 10)
	})

	t.Run("words with special characters", func(t *testing.T) {
		input := "слово-1 слово_2 слово#3 слово-1 слово_2 слово-1"
		expected := []string{"слово-1", "слово_2", "слово#3"}
		result := Top10(input)
		require.Equal(t, expected[:len(result)], result)
	})

	t.Run("latin characters", func(t *testing.T) {
		input := "unique New York, New York unique"
		expected := []string{"New", "unique", "York", "York,"}
		result := Top10(input)
		require.Equal(t, expected[:len(result)], result)
	})

	t.Run("mixed newlines and tabs", func(t *testing.T) {
		input := "слово1\nслово2\tслово3\nслово1\tслово2\nслово1"
		expected := []string{"слово1", "слово2", "слово3"}
		result := Top10(input)
		require.Equal(t, expected[:len(result)], result)
	})

	t.Run("positive test", func(t *testing.T) {
		if taskWithAsteriskIsCompleted {
			expected := []string{
				"а",         // 8
				"он",        // 8
				"и",         // 6
				"ты",        // 5
				"что",       // 5
				"в",         // 4
				"его",       // 4
				"если",      // 4
				"кристофер", // 4
				"не",        // 4
			}
			require.Equal(t, expected, Top10(text))
		} else {
			expected := []string{
				"он",        // 8
				"а",         // 6
				"и",         // 6
				"ты",        // 5
				"что",       // 5
				"-",         // 4
				"Кристофер", // 4
				"если",      // 4
				"не",        // 4
				"то",        // 4
			}
			require.Equal(t, expected, Top10(text))
		}
	})
}
