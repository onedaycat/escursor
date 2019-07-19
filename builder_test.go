package escursor

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestBuild(t *testing.T) {
	b := NewBuilder().
		Sort("a", map[string]interface{}{
			"order": "desc",
		}).
		Sort("b", map[string]interface{}{
			"order": "asc",
		}).
		Sort("c", map[string]interface{}{
			"order": "asc",
		}).
		Sort("d", map[string]interface{}{
			"order": "asc",
		}).
		Size(2)

	query := b.Build()
	require.Equal(t, map[string]interface{}{
		"size": 3,
		"sort": []map[string]interface{}{
			{
				"a": map[string]interface{}{
					"order": "desc",
				},
			},
			{
				"b": map[string]interface{}{
					"order": "asc",
				},
			},
			{
				"c": map[string]interface{}{
					"order": "asc",
				},
			},
			{
				"d": map[string]interface{}{
					"order": "asc",
				},
			},
		},
	}, query)

	nextToken, err := CreateNextToken(2, 3,
		func(index int) []interface{} {
			require.Equal(t, index, 1)
			return []interface{}{"string", 1, false, nil}
		},
		func(index int) {
			require.Equal(t, index, 1)
		})

	require.NoError(t, err)
	require.NotEmpty(t, nextToken)

	b = NewBuilder().
		Token(nextToken).
		Sort("a", map[string]interface{}{
			"order": "desc",
		}).
		Sort("b", map[string]interface{}{
			"order": "asc",
		}).
		Sort("c", map[string]interface{}{
			"order": "asc",
		}).
		Sort("d", map[string]interface{}{
			"order": "asc",
		}).
		Size(2)

	require.Equal(t, map[string]interface{}{
		"size": 3,
		"sort": []map[string]interface{}{
			{
				"a": map[string]interface{}{
					"order": "desc",
				},
			},
			{
				"b": map[string]interface{}{
					"order": "asc",
				},
			},
			{
				"c": map[string]interface{}{
					"order": "asc",
				},
			},
			{
				"d": map[string]interface{}{
					"order": "asc",
				},
			},
		},
		"search_after": []interface{}{"string", int64(1), false, nil},
	}, b.Build())
}

func TestBuildDataLessThanSize(t *testing.T) {
	nextToken, err := CreateNextToken(2, 1,
		func(index int) []interface{} {
			require.Equal(t, index, 0)
			return []interface{}{"string", 1, false, nil}
		},
		func(index int) {
			require.Equal(t, index, 0)
		})

	require.NoError(t, err)
	require.Empty(t, nextToken)

}

func TestBuildDataEqualThanSize(t *testing.T) {
	nextToken, err := CreateNextToken(2, 2,
		func(index int) []interface{} {
			require.Equal(t, index, 0)
			return []interface{}{"string", 1, false, nil}
		},
		func(index int) {
			require.Equal(t, index, 0)
		})

	require.NoError(t, err)
	require.Empty(t, nextToken)

}
