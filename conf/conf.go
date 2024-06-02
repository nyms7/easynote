package conf

var GlobalConf *NoteConf

func MaxContentSize() int {
	return GlobalConf.MaxContentSize
}

func MaxCodeSize() int {
	return GlobalConf.MaxCodeSize
}

func MaxTokenSize() int {
	return GlobalConf.MaxTokenSize
}

func AdminToken() string {
	return GlobalConf.AdminToken
}

type NoteConf struct {
	MaxCodeSize    int    `json:"max_code_size"`
	MaxContentSize int    `json:"max_content_size"`
	MaxTokenSize   int    `json:"max_token_size"`
	AdminToken     string `json:"-"`
}
