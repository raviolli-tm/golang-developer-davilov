package hw09structvalidator

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

type UserRole string

type UserRole2 struct {
	Role           string `validate:"in:admin,stuff"`
	WorkExperience int
}

type (
	User struct {
		ID     string `json:"id" validate:"len:36"`
		Name   string
		Age    int       `validate:"min:18|max:50"`
		Email  string    `validate:"regexp:^\\w+@\\w+\\.\\w+$"`
		Role   UserRole2 `validate:"nested"`
		Phones []string  `validate:"len:11"`
		meta   string    //nolint:unused
	}

	App struct {
		Version string `validate:"len:5"`
	}

	Response struct {
		Code int    `validate:"in:200,404,500"`
		Body string `json:"omitempty"`
	}

	Ints struct {
		Values []int `validate:"min:10|max:20"`
	}

	NestedPointer struct {
		Role *UserRole2 `validate:"nested"`
	}
)

func validUser() User {
	return User{
		ID:     strings.Repeat("a", 36),
		Name:   "dima",
		Age:    20,
		Email:  "mail@mail.ru",
		Role:   UserRole2{Role: "admin", WorkExperience: 2},
		Phones: []string{"12345678901"},
	}
}

func TestValidate(t *testing.T) {
	t.Run("valid struct", func(t *testing.T) {
		require.NoError(t, Validate(validUser()))
	})

	t.Run("pointer to struct", func(t *testing.T) {
		user := validUser()
		require.NoError(t, Validate(&user))
	})

	t.Run("string length", func(t *testing.T) {
		user := validUser()
		user.ID = "short"

		err := Validate(user)
		require.ErrorIs(t, err, ErrStringLength)

		var validationErrors ValidationErrors
		require.ErrorAs(t, err, &validationErrors)
		require.Len(t, validationErrors, 1)
		require.Equal(t, "ID", validationErrors[0].Field)
	})

	t.Run("min and max", func(t *testing.T) {
		user := validUser()
		user.Age = 17
		require.ErrorIs(t, Validate(user), ErrIntMin)

		user.Age = 51
		require.ErrorIs(t, Validate(user), ErrIntMax)
	})

	t.Run("regexp", func(t *testing.T) {
		user := validUser()
		user.Email = "not-an-email"
		require.ErrorIs(t, Validate(user), ErrStringRegexp)
	})

	t.Run("string in", func(t *testing.T) {
		user := validUser()
		user.Role.Role = "user"

		err := Validate(user)
		require.ErrorIs(t, err, ErrStringIn)

		var validationErrors ValidationErrors
		require.ErrorAs(t, err, &validationErrors)
		require.Len(t, validationErrors, 1)
		require.Equal(t, "Role.Role", validationErrors[0].Field)
	})

	t.Run("slice elements", func(t *testing.T) {
		user := validUser()
		user.Phones = []string{"12345678901", "123"}

		err := Validate(user)
		require.ErrorIs(t, err, ErrStringLength)

		var validationErrors ValidationErrors
		require.ErrorAs(t, err, &validationErrors)
		require.Len(t, validationErrors, 1)
		require.Equal(t, "Phones[1]", validationErrors[0].Field)
	})

	t.Run("all errors are accumulated", func(t *testing.T) {
		user := validUser()
		user.ID = "short"
		user.Age = 17
		user.Email = "invalid"
		user.Role.Role = "user"
		user.Phones = []string{"123"}

		var validationErrors ValidationErrors
		require.ErrorAs(t, Validate(user), &validationErrors)
		require.Len(t, validationErrors, 5)
	})

	t.Run("int in", func(t *testing.T) {
		require.NoError(t, Validate(Response{Code: 200}))
		require.ErrorIs(t, Validate(Response{Code: 201}), ErrIntIn)
	})

	t.Run("slice of ints", func(t *testing.T) {
		require.NoError(t, Validate(Ints{Values: []int{10, 15, 20}}))
		require.ErrorIs(t, Validate(Ints{Values: []int{9}}), ErrIntMin)
		require.ErrorIs(t, Validate(Ints{Values: []int{21}}), ErrIntMax)
	})

	t.Run("nested pointer", func(t *testing.T) {
		require.NoError(t, Validate(NestedPointer{Role: &UserRole2{Role: "stuff"}}))
		require.ErrorIs(t, Validate(NestedPointer{Role: &UserRole2{Role: "user"}}), ErrStringIn)
	})

	t.Run("not a struct", func(t *testing.T) {
		require.ErrorIs(t, Validate(42), ErrNotStruct)
		require.ErrorIs(t, Validate(nil), ErrNotStruct)

		var user *User
		require.ErrorIs(t, Validate(user), ErrNotStruct)
	})

	t.Run("unknown validator", func(t *testing.T) {
		type Bad struct {
			Value string `validate:"unknown:1"`
		}

		require.ErrorIs(t, Validate(Bad{}), ErrUnknownValidator)
	})

	t.Run("invalid tag", func(t *testing.T) {
		type Bad struct {
			Value int `validate:"min"`
		}

		require.ErrorIs(t, Validate(Bad{}), ErrInvalidTag)
	})

	t.Run("ignores fields without tag", func(t *testing.T) {
		require.NoError(t, Validate(App{Version: "12345"}))
		require.ErrorIs(t, Validate(App{Version: "1234"}), ErrStringLength)
	})

	t.Run("unsupported type", func(t *testing.T) {
		type Bad struct {
			Value float64 `validate:"min:1"`
		}

		require.ErrorIs(t, Validate(Bad{}), ErrUnsupportedType)
	})
}
