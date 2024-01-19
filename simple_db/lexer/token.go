package lexer

type Tag uint32

const (
	//AND 对应SQL关键字
	AND Tag = iota + 256
	//BREAK
	//DO
	EQ
	FALSE
	GE
	ID
	//IF
	//ELSE

	LE

	FLOAT
	MINUS
	PLUS
	NE
	NUM
	//OR
	REAL
	//TRUE
	//WHILE
	LEFT_BRACE    // "{"
	RIGHT_BRACE   // "}"
	LEFT_BRACKET  //"("
	RIGHT_BRACKET //")"
	AND_OPERATOR
	OR_OPERATOR
	ASSIGN_OPERATOR
	NEGATE_OPERATOR
	LESS_OPERATOR
	GREATER_OPERATOR
	//BASIC //对应int , float, bool, char 等类型定义
	//TEMP  //对应中间代码的临时寄存器变量
	//SEMICOLON

	//新增SQL对应关键字
	SELECT
	FROM
	WHERE
	INSERT
	INTO
	VALUES
	DELETE
	UPDATE
	SET
	CREATE
	TABLE
	INT
	VARCHAR
	VIEW
	AS
	INDEX
	ON
	OR
	COMMA
	STRING
	//SQL关键字定义结束
	EOF

	ERROR
)

var token_map = make(map[Tag]string)

func init() {
	//初始化SQL关键字对应字符串
	token_map[AND] = "AND"
	token_map[SELECT] = "SELECT"
	token_map[WHERE] = "where"
	token_map[INSERT] = "INSERT"
	token_map[INTO] = "INTO"
	token_map[VALUES] = "VALUES"
	token_map[DELETE] = "DELETE"
	token_map[UPDATE] = "UPDATE"
	token_map[SET] = "SET"
	token_map[CREATE] = "CREATE"
	token_map[TABLE] = "TABLE"
	token_map[INT] = "INT"
	token_map[VARCHAR] = "VARCHAR"
	token_map[VIEW] = "VIEW"
	token_map[AS] = "AS"
	token_map[INDEX] = "INDEX"
	token_map[ON] = "ON"
	token_map[COMMA] = ","

	//token_map[BASIC] = "BASIC"
	//token_map[DO] = "do"
	//token_map[ELSE] = "else"
	token_map[EQ] = "EQ"
	token_map[FALSE] = "FALSE"
	token_map[GE] = "GE"
	token_map[ID] = "ID"
	//token_map[IF] = "if"
	token_map[INT] = "int"
	token_map[FLOAT] = "float"

	token_map[LE] = "<="
	token_map[MINUS] = "-"
	token_map[PLUS] = "+"
	token_map[NE] = "!="
	token_map[NUM] = "NUM"
	token_map[OR] = "OR"
	token_map[REAL] = "REAL"
	//token_map[TEMP] = "t"
	//token_map[TRUE] = "TRUE"
	//token_map[WHILE] = "while"
	//token_map[DO] = "do"
	//token_map[BREAK] = "break"
	token_map[AND_OPERATOR] = "&"
	token_map[OR_OPERATOR] = "|"
	token_map[ASSIGN_OPERATOR] = "="
	token_map[NEGATE_OPERATOR] = "!"
	token_map[LESS_OPERATOR] = "<"
	token_map[GREATER_OPERATOR] = ">"
	token_map[LEFT_BRACE] = "{"
	token_map[RIGHT_BRACE] = "}"
	token_map[LEFT_BRACKET] = "("
	token_map[RIGHT_BRACKET] = ")"
	token_map[EOF] = "EOF"
	token_map[ERROR] = "ERROR"
	//token_map[SEMICOLON] = ";"

}

type Token struct {
	lexeme string
	Tag    Tag
}

func (t *Token) ToString() string {
	if t.lexeme == "" {
		return token_map[t.Tag]
	}

	return t.lexeme
}

func NewToken(tag Tag) Token {
	return Token{
		lexeme: "",
		Tag:    tag,
	}
}

func NewTokenWithString(tag Tag, lexeme string) *Token {
	return &Token{
		lexeme: lexeme,
		Tag:    tag,
	}
}
