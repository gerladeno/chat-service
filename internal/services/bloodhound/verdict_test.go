package bloodhound_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/gerladeno/chat-service/internal/services/bloodhound"
)

func TestVerdict_ZeroValue(t *testing.T) {
	var v bloodhound.Verdict
	assert.Equal(t, 0, v.FactsNumber())
	assert.True(t, v.Equals(bloodhound.VerdictOK))
}

func TestVerdict_HasFact(t *testing.T) {
	t.Run("zero verdict", func(t *testing.T) {
		v := bloodhound.NewVerdict()
		require.Equal(t, 0, v.FactsNumber())
		require.True(t, v.Equals(bloodhound.VerdictOK))

		for _, f := range []bloodhound.Fact{
			bloodhound.FactContainsCardCVC,
			bloodhound.FactContainsCardNumber,
			bloodhound.FactContainsCardDate,
			bloodhound.FactContainsSuspiciousPhrases,
		} {
			assert.Falsef(t, v.HasFact(f), "fact %d", f)
		}
	})

	t.Run("nonzero verdict", func(t *testing.T) {
		v := bloodhound.NewVerdict(
			bloodhound.FactContainsCardCVC,
			bloodhound.FactContainsSuspiciousPhrases,
		)
		require.Equal(t, 2, v.FactsNumber())

		for _, f := range []bloodhound.Fact{
			bloodhound.FactContainsCardCVC,
			bloodhound.FactContainsSuspiciousPhrases,
		} {
			assert.Truef(t, v.HasFact(f), "fact %d", f)
		}

		for _, f := range []bloodhound.Fact{
			bloodhound.FactContainsCardNumber,
			bloodhound.FactContainsCardDate,
		} {
			assert.Falsef(t, v.HasFact(f), "fact %d", f)
		}
	})
}

func TestVerdict_Equals(t *testing.T) {
	var v1, v2 bloodhound.Verdict
	assert.True(t, v1.Equals(v2))

	v1.AddFact(bloodhound.FactContainsSMSCode)
	v1.AddFact(bloodhound.FactContainsCardCVC)
	assert.False(t, v1.Equals(v2))

	v2.AddFact(bloodhound.FactContainsCardCVC)
	v2.AddFact(bloodhound.FactContainsSMSCode)
	assert.True(t, v1.Equals(v2))
}
