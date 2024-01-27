# Transaction em GO - operações de vários tabelas.

O que é basicamente é uma unica unidade de trabalho `Unit of Work`

Exemplo: Em nosso banco simples, queremos transferir 10 USD da conta 1 para a conta 2. Esta transação compreende 5 operações.

1. Criar um registro de transferencia com amount = 10
2. Criar unm registro de entrada na conta 1 com mount = -10 (Já que o dinheiro esta saindo da conta)
3. Criar um registro de entrada de conta para a conta 2 com amount = +10
4. Subtrai 10 da balance da conta 1
5. Adiciona 10 ao balance da conta 2
