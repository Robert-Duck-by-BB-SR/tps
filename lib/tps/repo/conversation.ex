defmodule TPS.Repo.Conversation do
  def new_convo(conn, users) do
    {:ok, statement} =
      conn
      |> Exqlite.Sqlite3.prepare("insert into conversation values (?1, ?2)")

    id = UUID.uuid4()

    :ok = Exqlite.Sqlite3.bind(statement, [id, users])

    case Exqlite.Sqlite3.step(conn, statement) do
      {:error, reason} ->
        Exqlite.Sqlite3.release(conn, statement)
        {:error, reason}

      :done ->
        Exqlite.Sqlite3.release(conn, statement)
        {:ok, id}
    end
  end

  def conversations_by_username(conn, username) do
    {:ok, statement} =
      conn
      |> Exqlite.Sqlite3.prepare("select * from conversation where users like '%?1%'")

    :ok = Exqlite.Sqlite3.bind(statement, [username])

    case Exqlite.Sqlite3.multi_step(conn, statement) do
      {:error, reason} ->
        Exqlite.Sqlite3.release(conn, statement)
        {:error, reason}

      {:rows, result} ->
        Exqlite.Sqlite3.release(conn, statement)
        {:ok, result}
    end
  end

  def conversation_users(conn, id) do
    {:ok, statement} =
      conn
      |> Exqlite.Sqlite3.prepare("select users from conversation where id=?1")

    :ok = Exqlite.Sqlite3.bind(statement, [id])

    case Exqlite.Sqlite3.step(conn, statement) do
      {:error, reason} ->
        Exqlite.Sqlite3.release(conn, statement)
        {:error, reason}

      {:row, result} ->
        Exqlite.Sqlite3.release(conn, statement)
        {:ok, result}
    end
  end
end
