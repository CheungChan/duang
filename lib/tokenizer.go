package lib

/**
 * 简化的词法分析器
 * 语法分析器从这里获取Token。
 */
type Tokenizer struct {
	tokens []Token
	pos    int
}

func NewTokenizer(tokens []Token) *Tokenizer {
	return &Tokenizer{tokens: tokens}
}
func (a *Tokenizer) Next() Token {
	//如果已经到了末尾，总是返回EOF
	token := a.tokens[a.pos]
	if a.pos <= len(a.tokens) {
		a.pos += 1
	}
	return token
}

func (a *Tokenizer) Position() int {
	return a.pos
}
func (a *Tokenizer) TraceBack(newPos int) {
	a.pos = newPos
}
