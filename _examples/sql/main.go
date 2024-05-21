// nolint: govet
package main

import (
	"github.com/alecthomas/kong"

	"github.com/alecthomas/repr"

	"github.com/alecthomas/participle/v2"
	"github.com/alecthomas/participle/v2/lexer"
)

type Boolean bool

func (b *Boolean) Capture(values []string) error {
	*b = values[0] == "TRUE"
	return nil
}

// Select based on http://www.h2database.com/html/grammar.html
type QueryParser struct {
	Pos    lexer.Position
	Select *SelectCmd `  "SELECT" @@ `
}

type SelectCmd struct {
	Pos           lexer.Position
	Distinct      *Distinct         `( "DISTINCT" @@ )?`
	All           *bool             `@("ALL")?`
	Expression    *SelectExpression `@@ ?`
	From          *From             `("FROM" (@@))?`
	SetOperations *SetOperations    `@@?`
}

type SetOperations struct {
	Pos        lexer.Position
	Operation  *string      `@SetOpeartions`
	All        *bool        `@("ALL")?`
	Distinct   *bool        `@("DISTINCT")?`
	OtherQuery *QueryParser `@@`
}

type From struct {
	Pos      lexer.Position
	SubQuery *QueryParser       `("("? @@ ")"? )?`
	Tables   []*TableExpression `  (@@ ( "," @@ )*)?`
	Join     []*Joins           `(@@ ( "" @@ )*)?`
	Where    *ConditionExpress  `( "WHERE" @@ )?`
}

type ConditionExpress struct {
	Pos lexer.Position
	OR  []*OrCondition ` @@ ( "OR" @@ )*`
}
type OrCondition struct {
	Pos lexer.Position
	And []*Condition ` @@ ( "AND" @@ )*  `
}

type Condition struct {
	Pos                 lexer.Position
	Not                 *bool        `(@"NOT")?`
	Openbracket         *string      `"("?`
	LHSFieldTableName   *string      `( @Ident ( ".")`
	LHSFieldName        *string      ` @Ident `
	Value               *Expression  `| "("? @@ ")"? )?`
	ValueIndex          *int         `("[" @Int "]")?`
	LHSTypecaste        *string      `("::" @Ident )?`
	Operator            *string      `(@( "<>" | "<=" | ">=" | "=" | "<" | ">" | "!=" | "<@"| "@>"| "&&"  )`
	Like                *bool        ` | @ "LIKE" `
	NotLike             *bool        ` | @ "NOT" "LIKE" `
	ILike               *bool        ` | @ "ILIKE" `
	NotILike            *bool        ` | @ "NOT" "ILIKE" `
	SimilarTo           *bool        ` | @ "SIMILAR" "TO" `
	NotSimilarTo        *bool        ` | @ "NOT" "SIMILAR" "TO" `
	Is                  *bool        `|  @ "IS" `
	Between             []*Value     `| "BETWEEN" ("("? @@ ( "AND"   @@ ) ")"?)`
	NotBetween          []*Value     `| "NOT" "BETWEEN" ("("? @@ ( "AND"   @@ ) ")"?)`
	In                  *bool        `| @ "IN" `
	NotIn               *bool        `| @ "NOT" "IN" )?`
	RHSFieldTableName   *string      `( @Ident ( ".")`
	RHSFieldName        *string      ` @Ident )?`
	Val                 []*Value     `("("? @@ ( ","   @@ )* ")"?)?`
	RHSTypecaste        *string      `("::" @Ident)?`
	AliasColumn         *string      `( "AS" @Ident )?`
	SubOperator         *string      `( (@SubQueryOperators)? `
	SubQuery            *QueryParser `"("? @@ ")"? )?`
	StuffUntilSemicolon []*string    ` (@Ignore ( @!(";") )* (";")?)?`
	Closebracket        *string      `")"?`
}

type Value struct {
	Pos          lexer.Position
	ArrayLHSCat  *int     `( @Int "||" )?`
	Array        []*Value ` "ARRAY" "[" @@ ("," @@)* "]"`
	ArrayRHSCat  *int     `("||" @Int )?`
	Int          *int     `| @Int`
	Float        *float64 `| @Float`
	String       *string  ` | @String`
	Null         *bool    ` | @"NULL" `
	Variable     *string  ` | @Ident `
	IndexofArray *int     `("[" @Int "]")?`
	NotNull      *bool    ` | @ "Not" "NULL" `
}

type TableExpression struct {
	Pos                 lexer.Position
	SchemaName          *string     `( @Ident (".") `
	SchemaTable         *string     `@Ident`
	Function            *string     `| @Functions`
	Functionparams      *Expression ` "(" @@ ")" `
	TableName           *string     `| @Ident) `
	AliasFunctionName   *string     `( "AS" @Ident  `
	AliasFunctionKeys   []*Value    `"(" @@ ( "," @@ )* ")" )?`
	AliasTableName      *string     `( ("AS")? @Ident )?`
	StuffUntilSemicolon []*string   ` (@Ignore ( @!(";") )* (";")?)?`
}

type Distinct struct {
	Pos       lexer.Position
	WithOn    *WithOn             `@@`
	WithoutOn []*SelectExpression `| @@`
}
type WithOn struct {
	Pos        lexer.Position
	On         *bool             `@("ON") `
	Columns    []string          ` "(" @Ident ( "," @Ident )* ")" `
	Expression *SelectExpression `@@ `
}

type SelectExpression struct {
	Pos                 lexer.Position
	Columns             []*AsExpression `"("? @@ ( "," @@ )* ")"?`
	StuffUntilSemicolon []*string       ` (@Ignore ( @!(";") )* (";")?)?`
}

type Aggregate struct {
	Pos       lexer.Position
	Aggregate *string `@Aggregate`
	Val       *string `("(" @Ident ")"`
	ValAll    *bool   `|"(" @"*"  ")")`
}
type AsExpression struct {
	Pos            lexer.Position
	FieldTableName *string       `(@Ident ( "." ) `
	ShortFiledName *string       `@Ident  `
	IndexofArray   *int          `("[" @Int "]")?`
	AsteriskExp    *bool         `| @"*"`
	Aggregate      *Aggregate    `| @@ `
	Function       *string       `| @Functions`
	Functionparams []*Expression ` "(" @@ ( "," @@ )* ")"  `
	Expression     *Expression   `| "("? @@ ")"? `
	SubQuery       *QueryParser  `| "("? @@ ")"? )`
	As             *string       `( ("AS")? (@Functions|@Ident))?`
}

type Expression struct {
	Pos        lexer.Position
	Value      *Value      `@@`
	Operator   *string     `(@("->>" | "->" | "#>>" | "#>" )`
	JsonColumn *Expression `@@ )?`
}

// JOINS
type Joins struct {
	Pos       lexer.Position
	JoinsType *string          `@Joins`
	Table     *TableExpression `@@`
	On        *bool            `@ "ON"`
	LHSTable  *string          `(@Ident ( "."  )`
	LHSField  *string          `@Ident)?`
	Operator  *string          `@( "<>" | "<=" | ">=" | "=" | "<" | ">" | "!=" | "<@"| "@>"| "&&"  ) ?`
	RHSTable  *string          `(@Ident ( ".")`
	RHSField  *string          `@Ident)?`
	Column    *Value           `(@@)?`
	// OtherJoin  *Joins           `(@@)?`
	Expression *Condition `(@@)?`
}

var (
	sqlLexer = lexer.MustSimple([]lexer.SimpleRule{
		{`Keyword`, `(?i)\b(SELECT|DISTINCT|AS|FROM|ON|WHERE|NOT|AND|OR|BETWEEN|ILIKE|SIMILAR|TO|BETWEEN|ARRAY|IS)\b`},
		{`SubQueryOperators`, `(?i)\b(ANY|ALL|EXISTS)\b`},
		{`Functions`, `(?i)\b(JSONB_ARRAY_ELEMENTS|JSONB_EACH|TRIM_ARRAY|ARRAY_CAT|ARRAY_APPEND|ARRAY_LENGTH|ARRAY_DIMS|ARRAY_NDIMS|ARRAY_LOWER|ARRAY_UPPER|ARRAY_TO_STRING|CARDINALITY|ARRAY_PREPEND|ARRAY_REPLACE|ARRAY_REMOVE|ARRAY_POSITION|ARRAY_POSITIONS|ARRAY_FILL)\b`},
		{`Ignore`, `(?i)\b(GROUP BY|HAVING|ORDER BY|LIMIT|FETCH|OFFSET|OVER)\b`},
		{`SetOpeartions`, `(?i)\b(UNION|INTERSECT|EXCEPT)\b`},
		{`Aggregate`, `(?i)\b(SUM|MIN|MAX|AVG|COUNT)\b`},
		{`Joins`, `(?i)\b(JOIN|INNER JOIN|FULL OUTER JOIN|LEFT JOIN|RIGHT JOIN|CROSS JOIN)\b`},
		{`Ident`, `[a-zA-Z_][a-zA-Z0-9_]*`},
		{`Float`, `-?\d*\.\d+`},
		{`Int`, `[0-9]+`},
		{`Operators`, `<>|!=|#>>|#>|<@|@>|->>|->|<=|>=|::|&&|[-+*/%,.()=<>]|\|\|`},
		{`String`, `'[^']*'|"[^"]*"`},
		{"Comment", `(?:#|--)[^\n]*\n?`},
		{"Semicolon", `(?:#|;)[^\n]*\n?`},
		{"Punct", `[-[!@#$%^&*()+_={}\|:;"'<,>.?/]|]`},
		{"Whitespace", `\s+`},
	})
	parser = participle.MustBuild(&QueryParser{},
		participle.Lexer(sqlLexer),
		participle.CaseInsensitive("Keyword", "Ignore", "Functions", "Ident", "Joins", "Aggregate", "SetOpeartions", "SubQueryOperators"),
		participle.Elide("Comment", "Semicolon", "Whitespace"),
		participle.UseLookahead(2),
		participle.Unquote("String"),
	)
)
func main() {
	ctx := kong.Parse(&cli)
	sql, err := parser.ParseString("", cli.SQL)
	repr.Println(sql, repr.Indent("  "), repr.OmitEmpty(true))
	ctx.FatalIfErrorf(err)
}
