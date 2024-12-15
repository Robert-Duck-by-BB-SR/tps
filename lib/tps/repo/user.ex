defmodule TPS.Repo.User do
  def new, do: "insert into user values (?1, ?2, ?3)"

  def username(conn, key) do
    {:ok, statement} =
      conn
      |> Exqlite.Sqlite3.prepare("select username from user where key=?1")

    :ok = Exqlite.Sqlite3.bind(statement, key)

    case Exqlite.Sqlite3.step(conn, statement) do
      {:error, reason} ->
        Exqlite.Sqlite3.release(conn, statement)
        {:error, reason}

      {:row, result} ->
        Exqlite.Sqlite3.release(conn, statement)
        {:ok, result}
    end
  end

  def users(conn) do
    {:ok, statement} =
      conn
      |> Exqlite.Sqlite3.prepare("select username from user")

    :ok = Exqlite.Sqlite3.bind(statement, [])

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
