package tests

type MyStructAnon struct {
	MyInt    int     `json:"my_int"`
	MyFloat  float64 `json:"my_float"`
	MyString string  `json:"my_string"`
	MyNested struct {
		MyNestedInt int `json:"nested_int"`
	} `json:"my_nested"`
	MyList []string `json:"my_list"`
}

type MyStructAnonPtr struct {
	MyInt    int     `json:"my_int"`
	MyFloat  float64 `json:"my_float"`
	MyString string  `json:"my_string"`
	MyNested *struct {
		MyNestedInt int `json:"nested_int"`
	} `json:"my_nested"`
	MyList []string `json:"my_list"`
}

type MyStruct struct {
	MyInt    int      `json:"my_int"`
	MyFloat  float64  `json:"my_float"`
	MyString string   `json:"my_string"`
	MyNested MyNested `json:"my_nested"`
	MyList   []string `json:"my_list"`
}

type MyStructPtr struct {
	MyInt    int       `json:"my_int"`
	MyFloat  float64   `json:"my_float"`
	MyString string    `json:"my_string"`
	MyNested *MyNested `json:"my_nested"`
	MyList   []string  `json:"my_list"`
}

type MyNested struct {
	MyNestedInt int `json:"nested_int"`
}
