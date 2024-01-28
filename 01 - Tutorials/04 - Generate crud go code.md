# Gerando código sql e integração com o banco em Go

Existem várias maneiras de implementar operações CRUD em golang.

## A primeira é usar o pacote banco de dados/sql da biblioteca padrão de baixo nível.

[Documentação da GOlang](https://pkg.go.dev/database/sql#DB.QueryContext)

```
package main

import(
    "context"
    "database/sql"
    "log"
    "time"
)
var (
    ctx Context.Context
    db *sql.DB
)
func main() {
    id := 123
    var username string
    var created time.Time
    err := db.QueryRowContext(ctx, "select username, created_at from user where id=?", id).Scan(&username, &created)
    switch {
        case err == sql.ErrNoRows:
                log.Printf("no user with id %d\n", id)
        case err != nil:
                log.Fatalf("query error: %v\n", err)
        default:
                log.Printf("Username is %q, account created on %s\n", username, created)
    }
}
```

> A principal vantagem desta abordagem é ele roda muito rápido, e escrever códigos é bastante simples.
> No entanto, sua desvantagem é termos que mapear manualmente os campos SQL para variáveis, o que é muito chato e fácil de cometer erros.
> Se de alguma forma a ordem das variáveis ​​não corresponder, ou se esquecermos de passar alguns argumentos para a chamada de função, os erros só aparecerão em tempo de execução.

## GORM - que é uma biblioteca de mapeamento objeto-relacional de alto nível para golang.

É super conveniente de usar

porque todas as operações CRUD já estão implementadas.

Portanto, nosso código de produção será muito curto,

pois só precisamos declarar os modelos

e chame as funções fornecidas pelo gorm.

[Docs GORM](https://gorm.io/docs/)

Exemplos na documentação:

### Create Record

```
user := User{Name: "Jinzhu", Age: 18, Birthday: time.Now()}

result := db.Create(&user) // pass pointer of data to Create

user.ID             // returns inserted data's primary key
result.Error        // returns error
result.RowsAffected // returns inserted records count
```

Na documentação do [create](https://gorm.io/docs/create.html) tem varios exemplos de create como

- multiple records with Create()
- Create Record With Selected Fields
- Batch Insert
- Create Hooks
- Create From Map
- Create From SQL Expression/Context Valuer
- Upsert / On Conflict

Exemplos de [Query](https://gorm.io/docs/query.html) (Consultas)

> mas o problema é:
> devemos aprender como escrever consultas usando as funções fornecidas pelo gorm.
> Será irritante se não soubermos quais funções usar.
> Especialmente quando temos algumas consultas complexas que exigem a união de tabelas,
> Temos que aprender como declarar tags de associação
> para fazer o gorm entender as relações entre as tabelas,
> Para que possa gerar a consulta SQL correta.

> OBS: Uma grande preocupação ao usar gorm é que
> Ele funciona muito lentamente quando o tráfego está alto.
> Existem alguns benchmarks na internet
> o que mostra que
> gorm pode ser executado de 3 a 5 vezes mais lento que a biblioteca padrão.

## SQLx uma abordagem intermediária

[Doc e git](https://github.com/jmoiron/sqlx)

- Funciona quase tão rápido quanto a biblioteca padrão.
- E é muito fácil de usar.

Ele fornece algumas funções como Select() ou StructScan(), que irá fazer o "scan" automaticamente o resultado nos campos struct, então não precisamos fazer o mapeamento manualmente como no pacote banco de dados/sql.
Isso ajudará a encurtar os códigos, e reduzir possíveis erros.

No entanto, o código que temos que escrever ainda é bastante longo,
E quaisquer erros na consulta só serão detectados em tempo de execução.

Então, existe alguma maneira melhor?

A resposta é SQLc

## SQLc

- Ele roda muito rápido, assim como o banco de dados/sql
- E é super fácil de usar
- O mais emocionante é, só precisamos escrever consultas SQL.

então os códigos golang CRUD serão gerados automaticamente para nós.

Como você pode ver neste exemplo,
Simplesmente passamos o esquema db e as consultas SQL para o sqlc,
Cada consulta tem 1 comentário em cima dela
para dizer ao SQLC para gerar a assinatura de função correta.

```
CREATE TABLE authors(
    id BIGSERIAL PRIMARY KEY,
    name text not null,
    bio text
)

-- name: GetAuthor :one
Select * from authors
order by name;

-- name: CreateAuthor :one
INSERT INTO authors (
    name, bio
) values (
    $1, $2
)
RETURNING *;

-- name: DeleteAuthor :exec
DELETE FROM author where id = $1;
```

Então o sqlc irá gerar códigos Golang idiomáticos,
que usa a biblioteca padrão de banco de dados/sql.

E porque o SQLC analisa as consultas SQL para entender o que faz
para gerar os códigos para nós,
portanto, quaisquer erros serão detectados e relatados imediatamente.
Parece incrível, certo?

> O único problema que encontrei no SQLC é que, No momento, Ele oferece suporte total apenas ao Postgres.

[Site](https://sqlc.dev/)

[Documentação SQLc](https://docs.sqlc.dev/en/stable/overview/install.html)

### instalação

[Documentação de instalação](https://docs.sqlc.dev/en/stable/overview/install.html)

```
brew install sqlc
```

### Passo a passo

1.  Executar o init `sqlc init` vai gerar o arquivo `sqlc.yaml`
2.  Configurar o arquivo:
    Posso seguir a doc [Getting started with PostgreSQL](https://docs.sqlc.dev/en/stable/tutorials/getting-started-postgresql.html)
    - Para configurar corretamente é necessarios seguie as seguintes configurações
    1. Criar o arquivo `sqlc.yaml`
    2. Configurar lo da seguinte maneira:

```
version: "2"
sql:
- engine: "postgresql"
queries: "./db/query"
schema: "./db/migration"
gen:
    go:
    package: "db"
    out: "./db/sqlc"
    sql_package: "pgx/v5"
```

Removendo as linhas `cloud project` e `database managed` conforme descrito na [documentação](https://docs.sqlc.dev/en/stable/tutorials/getting-started-postgresql.html#setting-up)

Também posso utilizar algumas configurações extras

![configs](img/config%20sqlc%20params.png)

> pode se utilizar a biblioteca pgx/v5. Basta incluir a linha sql_package: "pgx/v5" na config do sqlc.yaml linha abaixo do out:
> [DocRef - pgx](https://docs.sqlc.dev/en/stable/guides/using-go-and-pgx.html)

3. Executar o comando `sqlc generate`

   vai gerar a saida:

![saida](/img/sqlc%20out.png)

> OBS não preciso da pasta schema, ao configurar o sqlc.yaml eu aponto para a pasta db/migration ele já sabe que lá existe os create table

Os arquivo são:

- models.go as minhas struct com os models já incluindo as tags json `emit_json_tags: true`
- `db.go` Este arquivo contém a interface dbtx.
  Ele define 4 métodos comuns que os objetos sql.DB e sql.Tx possuem. Isso nos permite usar livremente um banco de dados ou uma transação para executar uma consulta.a função New() recebe um DBTX como entrada e retorna um objeto Queries. Então podemos passar um objeto sql.DB ou sql.Tx depende se queremos executar apenas uma única consulta, ou um conjunto de múltiplas consultas dentro de uma transação. Há também um método WithTx, que permite que uma instância de Queries seja associada a uma transação. (Falarei mais disso quando chegar em transações)

- `account.sql.go` O nome do pacote é “db” conforme definimos no arquivo sqlc.yaml. No começo já temos a consulta no `const createAccount` igual a que escrevemos no arquivo `account.sql` exceto pelo return. Então temos a struct `CreateAccountParams`, que contém todas as colunas que queremos definir.

4. Na sequencia é só dar um `go mod init` e `go mod tidy` para carregar os driver do Postgres

5. Configurações adicionais:
   1.

Comandos uteis:

- `sqlc version` versão do sqlc instalado
- `sqlc help` helper do cli
- `sqlc init` Create an empty sqlc.yaml settings file
- `sqlc complile` Statically check SQL for syntax and type errors
- `sqlc completion` Generate the autocompletion script for the specified shell
- `sqlc generate` Generate source code from SQL

Um pouco das [configurações](https://docs.sqlc.dev/en/stable/reference/config.html#version-2)

### Alternativas de sql

- Alternativa 2 de update:

```
-- name: UpdateAccount :exec
UPDATE account SET
balance = $2
where id = $1;
```
