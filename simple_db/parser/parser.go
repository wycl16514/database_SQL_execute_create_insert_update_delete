package parser

import (
	"fmt"
	"lexer"
	"query"
	"record_manager"
	"strconv"
	"strings"
)

type SQLParser struct {
	sqlLexer lexer.Lexer
}

func NewSQLParser(s string) *SQLParser {
	return &SQLParser{
		sqlLexer: lexer.NewLexer(s),
	}
}

func (p *SQLParser) UpdateCmd() interface{} {
	tok, err := p.sqlLexer.Scan()
	if err != nil {
		panic(err)
	}

	if tok.Tag == lexer.INSERT {
		p.sqlLexer.ReverseScan()
		return p.Insert()
	} else if tok.Tag == lexer.DELETE {
		p.sqlLexer.ReverseScan()
		return p.Delete()
	} else if tok.Tag == lexer.UPDATE {
		p.sqlLexer.ReverseScan()
		return p.Modify()
	} else {
		p.sqlLexer.ReverseScan()
		return p.Create()
	}
}

func (p *SQLParser) checkWordTag(wordTag lexer.Tag) {
	tok, err := p.sqlLexer.Scan()
	if err != nil {
		panic(err)
	}
	if tok.Tag != wordTag {
		panic("token is not match")
	}
}

func (p *SQLParser) isMatchTag(wordTag lexer.Tag) bool {
	tok, err := p.sqlLexer.Scan()
	if err != nil {
		panic(err)
	}

	if tok.Tag == wordTag {
		return true
	} else {
		p.sqlLexer.ReverseScan()
		return false
	}
}

func (p *SQLParser) Create() interface{} {
	tok, err := p.sqlLexer.Scan()
	if err != nil {
		panic(err)
	}
	if tok.Tag != lexer.CREATE {
		panic("token is not create")
	}

	tok, err = p.sqlLexer.Scan()
	if err != nil {
		panic(err)
	}

	if tok.Tag == lexer.TABLE {
		return p.CreateTable()
	} else if tok.Tag == lexer.VIEW {
		return p.CreateView()
	} else if tok.Tag == lexer.INDEX {
		return p.CreateIndex()
	}

	panic("sql string with create should not end here")
}

func (p *SQLParser) CreateView() interface{} {
	p.checkWordTag(lexer.ID)
	viewName := p.sqlLexer.Lexeme
	p.checkWordTag(lexer.AS)
	qd := p.Query()

	vd := NewViewData(viewName, qd)
	vdDef := fmt.Sprintf("vd def: %s", vd.ToString())
	fmt.Println(vdDef)
	return vd
}

func (p *SQLParser) CreateIndex() interface{} {
	p.checkWordTag(lexer.ID)
	idexName := p.sqlLexer.Lexeme
	p.checkWordTag(lexer.ON)
	p.checkWordTag(lexer.ID)
	tableName := p.sqlLexer.Lexeme
	p.checkWordTag(lexer.LEFT_BRACKET)
	_, fldName := p.Field()
	p.checkWordTag(lexer.RIGHT_BRACKET)

	idxData := NewIndexData(idexName, tableName, fldName)
	fmt.Printf("create index result: %s", idxData.ToString())
	return idxData
}

func (p *SQLParser) CreateTable() interface{} {
	tok, err := p.sqlLexer.Scan()
	if err != nil {
		panic(err)
	}
	if tok.Tag != lexer.ID {
		panic("token should be ID for table name")
	}

	tblName := p.sqlLexer.Lexeme
	tok, err = p.sqlLexer.Scan()
	if err != nil {
		panic(err)
	}
	if tok.Tag != lexer.LEFT_BRACKET {
		panic("missing left bracket")
	}
	sch := p.FieldDefs()
	tok, err = p.sqlLexer.Scan()
	if err != nil {
		panic(err)
	}
	if tok.Tag != lexer.RIGHT_BRACKET {
		panic("missing right bracket")
	}

	return NewCreateTableData(tblName, sch)
}

func (p *SQLParser) FieldDefs() *record_manager.Schema {
	schema := p.FieldDef()
	tok, err := p.sqlLexer.Scan()
	if err != nil {
		panic(err)
	}
	if tok.Tag == lexer.COMMA {
		schema2 := p.FieldDefs()
		schema.AddAll(schema2)
	} else {
		p.sqlLexer.ReverseScan()
	}

	return schema
}

func (p *SQLParser) FieldDef() *record_manager.Schema {
	_, fldName := p.Field()
	return p.FieldType(fldName)
}

func (p *SQLParser) FieldType(fldName string) *record_manager.Schema {
	schema := record_manager.NewSchema()
	tok, err := p.sqlLexer.Scan()
	if err != nil {
		panic(err)
	}

	if tok.Tag == lexer.INT {
		schema.AddIntField(fldName)
	} else if tok.Tag == lexer.VARCHAR {
		tok, err := p.sqlLexer.Scan()
		if err != nil {
			panic(err)
		}
		if tok.Tag != lexer.LEFT_BRACKET {
			panic("missing left bracket")
		}

		tok, err = p.sqlLexer.Scan()
		if err != nil {
			panic(err)
		}

		if tok.Tag != lexer.NUM {
			panic("it is not a number for varchar")
		}

		num := p.sqlLexer.Lexeme
		fldLen, err := strconv.Atoi(num)
		if err != nil {
			panic(err)
		}
		schema.AddStringField(fldName, fldLen)

		tok, err = p.sqlLexer.Scan()
		if err != nil {
			panic(err)
		}
		if tok.Tag != lexer.RIGHT_BRACKET {
			panic("missing right bracket")
		}
	}

	return schema
}

func (p *SQLParser) fieldList() []string {
	L := make([]string, 0)
	_, field := p.Field()
	L = append(L, field)
	if p.isMatchTag(lexer.COMMA) {
		fields := p.fieldList()
		L = append(L, fields...)
	}

	return L
}

func (p *SQLParser) constList() []*query.Constant {
	L := make([]*query.Constant, 0)
	L = append(L, p.Constant())
	if p.isMatchTag(lexer.COMMA) {
		consts := p.constList()
		L = append(L, consts...)
	}

	return L
}

func (p *SQLParser) Insert() interface{} {
	/*
		根据语法规则：Insert -> INSERT INTO ID LEFT_PARAS FieldList RIGHT_PARAS VALUES LEFT_PARS ConstList RIGHT_PARAS
		我们首先要匹配四个关键字，分别为insert, into, id, 左括号,
		然后就是一系列由逗号隔开的field,
		接着就是右括号，然后是关键字values
		接着是常量序列，最后以右括号结尾
	*/
	p.checkWordTag(lexer.INSERT)
	p.checkWordTag(lexer.INTO)
	p.checkWordTag(lexer.ID)
	tblName := p.sqlLexer.Lexeme
	p.checkWordTag(lexer.LEFT_BRACKET)
	flds := p.fieldList()
	p.checkWordTag(lexer.RIGHT_BRACKET)
	p.checkWordTag(lexer.VALUES)
	p.checkWordTag(lexer.LEFT_BRACKET)
	vals := p.constList()
	p.checkWordTag(lexer.RIGHT_BRACKET)

	return NewInsertData(tblName, flds, vals)
}

func (p *SQLParser) Delete() interface{} {
	/*
		第一个关键字 delete,第二个关键字必须 from
	*/
	p.checkWordTag(lexer.DELETE)
	p.checkWordTag(lexer.FROM)
	p.checkWordTag(lexer.ID)
	tblName := p.sqlLexer.Lexeme
	pred := query.NewPredicate()
	if p.isMatchTag(lexer.WHERE) {
		pred = p.Predicate()
	}
	return NewDeleteData(tblName, pred)
}

func (p *SQLParser) Modify() interface{} {
	p.checkWordTag(lexer.UPDATE)
	p.checkWordTag(lexer.ID)
	//获得表名
	tblName := p.sqlLexer.Lexeme
	p.checkWordTag(lexer.SET)
	_, fldName := p.Field()
	p.checkWordTag(lexer.ASSIGN_OPERATOR)
	newVal := p.Expression()
	pred := query.NewPredicate()
	if p.isMatchTag(lexer.WHERE) {
		pred = p.Predicate()
	}
	return NewModifyData(tblName, fldName, newVal, pred)
}

func (p *SQLParser) Field() (lexer.Token, string) {
	tok, err := p.sqlLexer.Scan()
	if err != nil {
		panic(err)
	}
	if tok.Tag != lexer.ID {
		panic("Tag of FIELD is no ID")
	}

	return tok, p.sqlLexer.Lexeme
}

func (p *SQLParser) Constant() *query.Constant {
	tok, err := p.sqlLexer.Scan()
	if err != nil {
		panic(err)
	}

	switch tok.Tag {
	case lexer.STRING:
		s := strings.Clone(p.sqlLexer.Lexeme)
		return query.NewConstantWithString(&s)
		break
	case lexer.NUM:
		//注意堆栈变量在函数执行后是否会变得无效
		v, err := strconv.Atoi(p.sqlLexer.Lexeme)
		if err != nil {
			panic("string is not a number")
		}
		return query.NewConstantWithInt(&v)
		break
	default:
		panic("token is not string or num when parsing constant")
	}

	return nil
}

func (p *SQLParser) Expression() *query.Expression {
	tok, err := p.sqlLexer.Scan()
	if err != nil {
		panic(err)
	}

	if tok.Tag == lexer.ID {
		p.sqlLexer.ReverseScan()
		_, str := p.Field()
		return query.NewExpressionWithString(str)
	} else {
		p.sqlLexer.ReverseScan()
		constant := p.Constant()
		return query.NewExpressionWithConstant(constant)
	}
}

func (p *SQLParser) Term() *query.Term {
	lhs := p.Expression()
	tok, err := p.sqlLexer.Scan()
	if err != nil {
		panic(err)
	}
	if tok.Tag != lexer.ASSIGN_OPERATOR {
		panic("should have = in middle of term")
	}

	rhs := p.Expression()
	return query.NewTerm(lhs, rhs)
}

func (p *SQLParser) Predicate() *query.Predicate {
	//predicate 对应where 语句后面的判断部分，例如where a > b and c < b
	//这里的a > b and c < b就是predicate
	pred := query.NewPredicateWithTerms(p.Term())
	tok, err := p.sqlLexer.Scan()
	// 如果语句已经读取完则直接返回
	if err != nil && fmt.Sprint(err) != "EOF" {
		panic(err)
	}

	if tok.Tag == lexer.AND {
		pred.ConjoinWith(p.Predicate())
	} else {
		p.sqlLexer.ReverseScan()
	}

	return pred
}

func (p *SQLParser) Query() *QueryData {
	//query 解析select 语句
	tok, err := p.sqlLexer.Scan()
	if err != nil {
		panic(err)
	}

	if tok.Tag != lexer.SELECT {
		panic("token is not select")
	}

	fields := p.SelectList()
	tok, err = p.sqlLexer.Scan()
	if err != nil {
		panic(err)
	}

	if tok.Tag != lexer.FROM {
		panic("token is not from")
	}

	//获取select语句作用的表名
	tables := p.TableList()
	//判断select语句是否有where子句
	tok, err = p.sqlLexer.Scan()
	pred := query.NewPredicate()
	if err == nil && tok.Tag == lexer.WHERE {
		pred = p.Predicate()
	} else {
		p.sqlLexer.ReverseScan()
	}

	return NewQueryData(fields, tables, pred)
}

func (p *SQLParser) SelectList() []string {
	//SELECT_LIST 对应select关键字后面的列名称
	l := make([]string, 0)
	_, field := p.Field()
	l = append(l, field)

	tok, err := p.sqlLexer.Scan()
	if err != nil {
		panic(err)
	}
	if tok.Tag == lexer.COMMA {
		//selct 多个列，每个列由逗号隔开
		selectList := p.SelectList()
		l = append(l, selectList...)
	} else {
		p.sqlLexer.ReverseScan()
	}

	return l
}

func (p *SQLParser) TableList() []string {
	//TBALE_LSIT对应from后面的表名
	l := make([]string, 0)
	tok, err := p.sqlLexer.Scan()
	if err != nil {
		panic(err)
	}
	if tok.Tag != lexer.ID {
		panic("token is not id")
	}

	l = append(l, p.sqlLexer.Lexeme)
	tok, err = p.sqlLexer.Scan()
	//change here
	if err == nil && tok.Tag == lexer.COMMA {
		tableList := p.TableList()
		l = append(l, tableList...)
	} else {
		p.sqlLexer.ReverseScan()
	}

	return l
}
