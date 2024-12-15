defmodule TPS.Repo.Message do
  def new_message(conn, [type, user, conversation, datetime, message]) do
    {:ok, statement} =
      conn
      |> Exqlite.Sqlite3.prepare("insert into message values (?1, ?2, ?3, ?4, ?5, ?6)")

    id = UUID.uuid4()
    :ok = Exqlite.Sqlite3.bind(statement, [id, type, user, conversation, datetime, message])

    case Exqlite.Sqlite3.step(conn, statement) do
      {:error, reason} ->
        Exqlite.Sqlite3.release(conn, statement)
        {:error, reason}

      :done ->
        Exqlite.Sqlite3.release(conn, statement)
        :ok
    end
  end

  def messages_in_convo(conn, conversation) do
    {:ok, statement} =
      conn
      |> Exqlite.Sqlite3.prepare("select * from message where conversation=?1")

    :ok = Exqlite.Sqlite3.bind(statement, [conversation])

    case Exqlite.Sqlite3.multi_step(conn, statement) do
      {:error, reason} ->
        Exqlite.Sqlite3.release(conn, statement)
        {:error, reason}

      {:done, result} ->
        Exqlite.Sqlite3.release(conn, statement)
        {:ok, result}
    end
  end
end
