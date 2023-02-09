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

%token <val> S

%start request

%%

request: S S headers
    {
        $$ = Result{method: $1, url: $2, headers: $3}
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
%%