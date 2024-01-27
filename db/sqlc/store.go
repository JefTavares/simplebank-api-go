package db

import (
	"context"
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
)

/* Implementa as transaction - fornecerá todas as funções
   para executar consultas de banco de dados individualmente,
*/

// Store defines all functions to execute db queries and transactions
type Store struct {
	*Queries
	db *sql.DB
}

// NewStore creates a new store
func NewStore(db *sql.DB) *Store {
	//Exemplo utilizando o pgx
	//ctx := context.Background()
	//dbSource := "postgresql://root:secret@localhost:5432/simple_bank?sslmode=disable"
	//conn, err := pgx.Connect(ctx, dbSource)
	// if err != nil {
	// 	fmt.Println("erro NewStore >>> ", err)
	// }
	return &Store{
		db:      db,
		Queries: New(db),
	}
}

// execTx execute a function within a database transaction
// Execute a generic DB transaction
func (store *Store) execTx(ctx context.Context, fn func(*Queries) error) error {
	tx, err := store.db.BeginTx(ctx, &sql.TxOptions{})
	if err != nil {
		return err
	}

	q := New(tx)
	err = fn(q)
	if err != nil {
		if rbErr := tx.Rollback(); rbErr != nil {
			return fmt.Errorf("tx err: %v, rb err: %v", err, rbErr)
		}
		return err
	}

	return tx.Commit()
}

// TransferTxParams Esta estrutura contém todos os parâmetros de entrada necessários para transferir dinheiro entre 2 contas.
type TransferTxParams struct {
	FromAccountID int64 `json:"from_account_id"`
	ToAccountID   int64 `json:"to_account_id"`
	Amount        int64 `json:"amount"`
}

// A estrutura TransferTxResult contém o resultado da operação de transferência.
type TransferTxResult struct {
	Transfer    Transfer `json:"transfer"`
	FromAccount Account  `json:"from_account"`
	ToAccount   Account  `json:"to_account"`
	FromEntry   Entry    `json:"from_entry"`
	ToEntry     Entry    `json:"to_entry"`
}

//criado para testar transferencias passando o context.WithValue
//uma transação com nome tx+n
//var txKey = struct{}{}

// TransferTx performs a money transfer from one account to the other.
// It creates a transfer record, add account entries, and update accounts balance withina single database transaction
func (store *Store) TransferTx(ctx context.Context, arg TransferTxParams) (TransferTxResult, error) {
	var result TransferTxResult
	err := store.execTx(ctx, func(q *Queries) error {

		var err error

		//txName := ctx.Value(txKey)

		//fmt.Println(txName, "create transfer")
		//Registra uma transferencia. quem envia para quem recebe e valor
		result.Transfer, err = q.CreateTransfer(ctx, CreateTransferParams{
			FromAccountID: arg.FromAccountID,
			ToAccountID:   arg.ToAccountID,
			Amount:        arg.Amount, //Amount:        sql.NullInt64{Int64: arg.Amount
		})
		if err != nil {
			return err
		}

		//fmt.Println(txName, "create entry 1")
		//Registra uma entrada de quem esta enviando
		result.FromEntry, err = q.CreateEntry(ctx, CreateEntryParams{
			AccountID: arg.FromAccountID,
			Amount:    -arg.Amount,
		})
		if err != nil {
			return err
		}

		//fmt.Println(txName, "create entry 2")
		//Registra uma entrada para quem esta recebendo
		result.ToEntry, err = q.CreateEntry(ctx, CreateEntryParams{
			AccountID: arg.ToAccountID,
			Amount:    arg.Amount,
		})
		if err != nil {
			return err
		}

		//fmt.Println(txName, "get account 1")
		// get account -> update its balance (Atualiza as contas)
		// account1, err := q.GetAccountForUpdate(ctx, arg.FromAccountID)
		// if err != nil {
		// 	return err
		// }
		/* removido o exemplo de pegar a conta e fazer o calculo balance - amount */
		//fmt.Println(txName, "update account 1 balance")
		//Atualiza a conta de quem esta enviando e armazena no result.FromAccount
		// result.FromAccount, err = q.UpdateAccount(ctx, UpdateAccountParams{
		// 	ID:      arg.FromAccountID,
		// 	Balance: account1.Balance - arg.Amount,
		// })
		// if err != nil {
		// 	return err
		// }

		if arg.FromAccountID < arg.ToAccountID {
			result.FromAccount, result.ToAccount, err = addMoney(ctx, q, arg.FromAccountID, -arg.Amount, arg.ToAccountID, arg.Amount)
			/* trocado pela função addMoney*/
			// result.FromAccount, err = q.AddAccountBalance(ctx, AddAccountBalanceParams{
			// 	ID:     arg.FromAccountID,
			// 	Amount: -arg.Amount,
			// })
			if err != nil {
				return err
			}

			/* removido o exemplo de pegar a conta e fazer o calculo balance + amount */
			//fmt.Println(txName, "get account 2")
			// account2, err := q.GetAccountForUpdate(ctx, arg.ToAccountID)
			// if err != nil {
			// 	return err
			// }

			//fmt.Println(txName, "update account 2 balance")
			// result.ToAccount, err = q.UpdateAccount(ctx, UpdateAccountParams{
			// 	ID:      arg.ToAccountID,
			// 	Balance: account2.Balance + arg.Amount,
			// })
			// if err != nil {
			// 	return err
			// }

			/* trocado pela função addMoney*/
			// result.ToAccount, err = q.AddAccountBalance(ctx, AddAccountBalanceParams{
			// 	ID:     arg.ToAccountID,
			// 	Amount: arg.Amount,
			// })
			// if err != nil {
			// 	return err
			// }
		} else {
			result.FromAccount, result.ToAccount, err = addMoney(ctx, q, arg.ToAccountID, arg.Amount, arg.FromAccountID, -arg.Amount)
			/* trocado pela função addMoney*/
			// result.ToAccount, err = q.AddAccountBalance(ctx, AddAccountBalanceParams{
			// 	ID:     arg.ToAccountID,
			// 	Amount: arg.Amount,
			// })
			if err != nil {
				return err
			}

			// result.FromAccount, err = q.AddAccountBalance(ctx, AddAccountBalanceParams{
			// 	ID:     arg.FromAccountID,
			// 	Amount: -arg.Amount,
			// })
			// if err != nil {
			// 	return err
			// }

		}

		return nil
	})

	return result, err

}

func addMoney(
	ctx context.Context,
	q *Queries,
	accountID1 int64,
	amount1 int64,
	accountID2 int64,
	amount2 int64,
) (account1 Account, account2 Account, err error) {
	account1, err = q.AddAccountBalance(ctx, AddAccountBalanceParams{
		ID:     accountID1,
		Amount: amount1,
	})
	if err != nil {
		return //retorna null já que os retorno estão nomeados
	}

	account2, err = q.AddAccountBalance(ctx, AddAccountBalanceParams{
		ID:     accountID2,
		Amount: amount2,
	})
	if err != nil {
		return //retorna null já que os retorno estão nomeados
	}

	return
}
