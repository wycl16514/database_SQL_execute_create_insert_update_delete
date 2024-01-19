package lexer

type Word struct {
	lexeme string
	Tag    Token
}

func NewWordToken(s string, tag Tag) Word {
	return Word{
		lexeme: s,
		Tag:    NewToken(tag),
	}
}

func (w *Word) ToString() string {
	return w.lexeme
}

func GetKeyWords() []Word {
	key_words := []Word{}
	key_words = append(key_words, NewWordToken("||", OR))
	key_words = append(key_words, NewWordToken("==", EQ))
	key_words = append(key_words, NewWordToken("!=", NE))
	key_words = append(key_words, NewWordToken("<=", LE))
	key_words = append(key_words, NewWordToken(">=", GE))
	//增加SQL语言对应关键字
	key_words = append(key_words, NewWordToken("AND", AND))
	key_words = append(key_words, NewWordToken("SELECT", SELECT))
	key_words = append(key_words, NewWordToken("WHERE", WHERE))
	key_words = append(key_words, NewWordToken("FROM", FROM))
	key_words = append(key_words, NewWordToken("INSERT", INSERT))
	key_words = append(key_words, NewWordToken("INTO", INTO))
	key_words = append(key_words, NewWordToken("VALUES", VALUES))
	key_words = append(key_words, NewWordToken("DELETE", DELETE))
	key_words = append(key_words, NewWordToken("UPDATE", UPDATE))
	key_words = append(key_words, NewWordToken("SET", SET))
	key_words = append(key_words, NewWordToken("CREATE", CREATE))
	key_words = append(key_words, NewWordToken("TABLE", TABLE))
	key_words = append(key_words, NewWordToken("INT", INT))
	key_words = append(key_words, NewWordToken("VARCHAR", VARCHAR))
	key_words = append(key_words, NewWordToken("VIEW", VIEW))
	key_words = append(key_words, NewWordToken("AS", AS))
	key_words = append(key_words, NewWordToken("INDEX", INDEX))
	key_words = append(key_words, NewWordToken("ON", ON))

	//key_words = append(key_words, NewWordToken("minus", MINUS))
	//key_words = append(key_words, NewWordToken("true", TRUE))
	//key_words = append(key_words, NewWordToken("false", FALSE))
	//key_words = append(key_words, NewWordToken("if", IF))
	//key_words = append(key_words, NewWordToken("else", ELSE))
	//增加while, do关键字
	//key_words = append(key_words, NewWordToken("while", WHILE))
	//key_words = append(key_words, NewWordToken("do", DO))
	//key_words = append(key_words, NewWordToken("break", BREAK))
	//添加类型定义
	//key_words = append(key_words, NewWordToken("int", BASIC))
	//key_words = append(key_words, NewWordToken("float", BASIC))
	//key_words = append(key_words, NewWordToken("bool", BASIC))
	//key_words = append(key_words, NewWordToken("char", BASIC))

	return key_words
}
