# Migration no GO

### biblioteca utilizada.

Ferramente escrita em GO [Migrate](https://github.com/golang-migrate/migrate) é uma CLI completa para migração de banco de dados universal

### Instalação

[Link](https://github.com/golang-migrate/migrate/tree/master/cmd/migrate)

pode ser via CLI como esta na documentação mas preferi utilizar o brew.

`brew install golang-migrate`

## Commandos

O primeiro é o create.
para isso vamos criar o primeiro arquivo de migração
Para inicializar o esquema de banco de dados do nosso simple_bank.
na pasta do projeto criar as pasta `db/migration`

Rodar o comando:

```
migrate create -ext sql -dir db/migration -seq init_schema
```

> OBS: -seq para gerar um número de versão sequencial para o arquivo de migração init_schema e finalmente o nome da migração, que é “init_schema” neste caso.

São gerados 2 aqruivos. O script `up` é executado para fazer uma alteração direta no esquema. E o script `down` é executado se quisermos reverter a alteração feita pelo script `up`.

No exemplo foi criado os `000001_init_schema.down.sql` e `000001_init_schema.up.sql` os dois estão vazio

No `000001_init_schema.up.sql` Colar a criação dos banco.

> exemplo do arquivo gerado anteriormente go-simples-bank.sql que está na posta do projeto

No arquivo `000001_init_schema.down.sql` desfazer oq foi feito

> no anterior (os drop table) Exemplo:

```
DROP TABLE IF EXISTS entries;
DROP TABLE IF EXISTS transfers;
DROP TABLE IF EXISTS accounts;
```

> Não esquecer de por na ordem (como tenho chaves estrangeiras)

## Parte Docker o banco de dados

> aqui são só exemplos

Acessar o container postgres e criar um novo database

```
docker exec -it 614 /bin/sh
createdb --username=postgres --owner=postgres simple_bank
```

Ou de fora do container, passando o nome do container (NAMES = postgres_jeftavares)

```
docker exec -it postgres_jeftavares createdb --username=postgres --owner=postgres simple_bank
```

Alterações no banco via docker

```
docker exec -it postgres_jeftavares psql -U postgres simples_bank
```

mas qualquer alteração direto é a boa conectar no banco

```
docker exec -it 614 /bin/sh
ls -ls (visualizar arquivos e comandos dentro do sh)
```

## Neste projeto utilizo esses comando atraves do arquivo Makefile

no arquivo makefile adicionei a primeira linha:

```
postgres:
	docker run --name postgres12 -p 5432:5432 -e POSTGRES_USER=root POSTGRES_PASSWORD=secret -d postgres:12-alpine
```

> para seguir os moldes do curso, claro que pode ser feito via docker compose ou utlizar um banco de dados em produção na sequencia.

### Rodar o Make

no arquivo makefile pronto eu executo em dev
`make postgres`
e
`make createdb`

## Proximo passo "finalmente fazer a migration"

```
migrate -path db/migration -database "postgresql://root:secret@localhost:5432/simple_bank?sslmode=disable" -verbose up
```

Detalhe dos comandos:

- `-verbose` Usamos a opção `-verbose` para solicitar a migração para imprimir o registro detalhado
- `?sslmode=disable` Isso ocorre porque nosso contêiner postgres não habilita SSL por padrão. Portanto, devemos adicionar o parâmetro “sslmode=disable” à URL do banco de dados.

apos rodar ele criar as tabelas e faz a migração:

detalhe para a tabela `schema_migration` que traz os detalhes da migração e a versão a primeira coluna é a versão no caso a 1 o 000001_init_schema.up do arquivo

## Uteis:

```
migrate -version
migrate -help
```

## Dicas

historico de comando no cli:

```
history | grep "docker run"
history | grep "docker exec"
```

Comandos docker:

```
docker stop NAMES_CONTAINER
docker ps -a
docker rm NAMES_CONTAINER
```

criar um tipo enum no postgres

```
CREATE TYPE "Currency" AS ENUM (
  'USD',
  'EUR'
);

```

Configuração do dbeaver

![Config](img/dbaver%20conexao%20banco%20pg.png)
