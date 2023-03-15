%{
package parser

func setResult(l yyLexer, v Result) {
  l.(*lex).result = v
}

%}

%union{
    result Result
    val string
    hh map[string]string
}

%type <result> request
%type <hh> headers
%type <hh> header
%type <hh> flags

%token <val> S
%token <val> Flag

%start request

%%

request: S S headers flags
    {
        $$ = Result{Method: $1, Url: $2, Headers: $3}
        setResult(yylex, $$)
    }

headers: headers header
    { $$ = merge($2, $1) }
| {}

header: S ':' S
    {
        $$ = map[string]string{
            $1: $3,
        }
    }

flags: flags Flag
    { $$ = merge(mapOf($2, ""), $1)}
| {}
%%