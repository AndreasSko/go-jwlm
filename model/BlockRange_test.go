package model

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"testing"
	"text/tabwriter"

	"github.com/stretchr/testify/assert"
)

func TestBlockRange_SetID(t *testing.T) {
	m1 := &BlockRange{BlockRangeID: 1}
	m1.SetID(10)
	assert.Equal(t, 10, m1.BlockRangeID)

	m2 := BlockRange{BlockRangeID: 2}
	m2.SetID(20)
	assert.Equal(t, 20, m2.BlockRangeID)
}

func TestBlockRange_UniqueKey(t *testing.T) {
	m1 := &BlockRange{
		BlockRangeID: 1,
		BlockType:    1,
		Identifier:   1,
		StartToken:   sql.NullInt32{Int32: 1, Valid: true},
		EndToken:     sql.NullInt32{Int32: 2, Valid: true},
		UserMarkID:   1,
	}
	assert.Equal(t, "1_1_1_2_1", m1.UniqueKey())

	m2 := &BlockRange{
		BlockRangeID: 2,
		BlockType:    1,
		Identifier:   20,
		StartToken:   sql.NullInt32{Int32: 15, Valid: true},
		EndToken:     sql.NullInt32{Int32: 25, Valid: true},
		UserMarkID:   3334,
	}
	assert.Equal(t, "1_20_15_25_3334", m2.UniqueKey())
}

func TestBlockRange_PrettyPrint(t *testing.T) {
	m1 := &BlockRange{
		BlockRangeID: 1,
		BlockType:    2,
		Identifier:   3,
		StartToken:   sql.NullInt32{Int32: 4, Valid: true},
		EndToken:     sql.NullInt32{Int32: 5, Valid: true},
		UserMarkID:   6,
	}

	buf := new(bytes.Buffer)
	w := tabwriter.NewWriter(buf, 0, 0, 1, ' ', 0)
	fmt.Fprint(w, "\nIdentifier:\t3")
	fmt.Fprint(w, "\nStartToken:\t4")
	fmt.Fprint(w, "\nEndToken:\t5")
	w.Flush()
	expectedResult := buf.String()

	assert.Equal(t, expectedResult, m1.PrettyPrint(nil))

	m1.EndToken.Valid = false

	buf.Reset()
	fmt.Fprint(w, "\nIdentifier:\t3")
	fmt.Fprint(w, "\nStartToken:\t4")
	w.Flush()
	expectedResult = buf.String()

	assert.Equal(t, expectedResult, m1.PrettyPrint(nil))
}

func TestBlockRange_Equals(t *testing.T) {
	m1 := &BlockRange{
		BlockRangeID: 1,
		BlockType:    1,
		Identifier:   1,
		StartToken:   sql.NullInt32{Int32: 1, Valid: true},
		EndToken:     sql.NullInt32{Int32: 2, Valid: true},
		UserMarkID:   1,
	}
	m1_1 := &BlockRange{
		BlockRangeID: 10000,
		BlockType:    1,
		Identifier:   1,
		StartToken:   sql.NullInt32{Int32: 1, Valid: true},
		EndToken:     sql.NullInt32{Int32: 2, Valid: true},
		UserMarkID:   1,
	}

	m2 := &BlockRange{
		BlockRangeID: 2,
		BlockType:    1,
		Identifier:   20,
		StartToken:   sql.NullInt32{Int32: 15, Valid: true},
		EndToken:     sql.NullInt32{Int32: 25, Valid: true},
		UserMarkID:   3334,
	}

	assert.True(t, m1.Equals(m1_1))
	assert.False(t, m1.Equals(m2))
}

func TestBlockRange_MarshalJSON(t *testing.T) {
	m1 := &BlockRange{
		BlockRangeID: 1,
		BlockType:    1,
		Identifier:   1,
		StartToken:   sql.NullInt32{Int32: 1, Valid: true},
		EndToken:     sql.NullInt32{Int32: 2, Valid: true},
		UserMarkID:   1,
	}
	result, err := json.Marshal(m1)
	assert.NoError(t, err)
	assert.Equal(t, `{"Type":"BlockRange","BlockRangeID":1,"BlockType":1,"Identifier":1,"StartToken":{"Int32":1,"Valid":true},"EndToken":{"Int32":2,"Valid":true},"UserMarkID":1}`, string(result))
}
