defmodule TPS.Repo do
  alias TPS.Repo
  require Logger
  use GenServer
  require Exqlite

  @impl true
  def init(db_name) do
    Exqlite.Sqlite3.open(db_name)
  end

  @impl true
  def handle_call({:select, query, values}, _from, conn) do
    {:ok, statement} =
      conn
      |> Exqlite.Sqlite3.prepare(query)

    :ok = Exqlite.Sqlite3.bind(statement, values)

    case Exqlite.Sqlite3.multi_step(conn, statement) do
      {:error, reason} ->
        Exqlite.Sqlite3.release(conn, statement)
        {:stop, :error, reason, conn}

      {_, result} ->
        Exqlite.Sqlite3.release(conn, statement)
        {:reply, result, conn}
    end
  end

  @impl true
  def handle_call({:get, query, values}, _from, conn) do
    {:ok, statement} =
      conn
      |> Exqlite.Sqlite3.prepare(query)

    :ok = Exqlite.Sqlite3.bind(statement, values)

    case Exqlite.Sqlite3.step(conn, statement) do
      {:error, reason} ->
        Exqlite.Sqlite3.release(conn, statement)
        {:stop, :error, reason, conn}

      {:row, result} ->
        Exqlite.Sqlite3.release(conn, statement)
        {:reply, result, conn}

      :done ->
        {:reply, :success, conn}
    end
  end

  @impl true
  def handle_call([["get", "conversation"], ["key", key] | _], _from, conn) do
    with {:ok, [username]} <- Repo.User.username(conn, key),
         {:ok, rows} <- Repo.Conversation.conversations_by_username(conn, username) do
      conversations =
        rows
        |> Enum.reduce("", fn row, acc ->
          [id, users] = row
          "#{acc};#{id}:#{users}"
        end)

      {:reply, {:ok, conversations}, conn}
    else
      {:error, reason} ->
        Logger.error(reason)
        {:reply, {:error, reason}, conn}
    end
  end

  @impl true
  def handle_call([["get", "users"] | _], _from, conn) do
    case Repo.User.users(conn) do
      {:error, reason} ->
        Logger.error(reason)
        {:reply, {:error, reason}, conn}

      {:ok, users} ->
        Logger.warning("printing users")
        IO.inspect(users)

        usernames =
          users
          |> Enum.join(";")

        IO.inspect(usernames)

        {:reply, {:ok, usernames}, conn}
    end
  end

  def start_link([name, opts]) do
    GenServer.start_link(__MODULE__, name, opts)
  end

  def query_raw(type, query_string, values) do
    GenServer.call(__MODULE__, {type, query_string, values})
  end

  def query(request) do
    GenServer.call(__MODULE__, request)
  end
end
