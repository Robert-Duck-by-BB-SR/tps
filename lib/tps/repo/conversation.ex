defmodule TPS.Repo.Conversation do
  def new_convo, do: "insert into conversation values (?1, ?2)"

  def(conversations_by_username(conn, username)) do
    {:ok, statement} =
      conn
      |> Exqlite.Sqlite3.prepare("select * from conversation where users like '%?1%'")

    :ok = Exqlite.Sqlite3.bind(statement, username)

    case Exqlite.Sqlite3.multi_step(conn, statement) do
      {:error, reason} ->
        Exqlite.Sqlite3.release(conn, statement)
        {:error, reason}

      {:rows, result} ->
        Exqlite.Sqlite3.release(conn, statement)
        {:ok, result}
    end
  end
end
