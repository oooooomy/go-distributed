package user

type User struct {
	ID       string
	Username string
	Nickname string
	Password string
	Address  string
}

var users []User = []User{
	{ID: "1",
		Username: "618736123",
		Nickname: "Jack",
		Password: "123456",
		Address:  "xxx",
	},
	{ID: "2",
		Username: "127512676",
		Nickname: "Tom",
		Password: "0.00.00.0",
		Address:  "xxx",
	},
	{ID: "3",
		Username: "787837811",
		Nickname: "Mary",
		Password: "abc145845",
		Address:  "xxx",
	},
}
