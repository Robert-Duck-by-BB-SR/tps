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
  def handle_call([["create", "conversation"], ["key", key], ["users", users]], _from, conn) do
    {:ok, decoded_key} = Base.decode16(key)

    with {:ok, [username]} <- Repo.User.username(conn, decoded_key),
         {:ok, id} <- Repo.Conversation.new_convo(conn, "#{username}|#{users}") do
      {:reply, {:ok, id}, conn}
    else
      {:error, reason} ->
        Logger.error(reason)
        {:reply, {:error, reason}, conn}
    end
  end

  @impl true
  def handle_call(
        [
          ["get", "messages"],
          ["key", key],
          ["conversation", conversation]
        ],
        _from,
        conn
      ) do
    Logger.warning(key)
    Logger.warning(conversation)
    {:ok, decoded_key} = Base.decode16(key)

    with {:ok, _} <- Repo.User.username(conn, decoded_key),
         {:ok, rows} <-
           Repo.Message.messages_in_convo(conn, conversation) do
      messages =
        rows
        |> Enum.reduce("", fn row, acc ->
          [_id, type, user, conversation, datetime, message] = row
          "#{acc};#{type}|#{user}|#{conversation}|#{datetime}|#{message}"
        end)

      {:reply, {:ok, messages}, conn}
    else
      {:error, reason} ->
        Logger.error(reason)
        {:reply, {:error, reason}, conn}
    end
  end

  @impl true
  def handle_call([["get", "conversation"], ["key", key] | _], _from, conn) do
    {:ok, decoded_key} = Base.decode16(key)

    with {:ok, [username]} <- Repo.User.username(conn, decoded_key),
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

  @impl true
  def handle_call([type, key, conversation, datetime, message], _from, conn) do
    {:ok, decoded_key} = Base.decode16(key)

    with {:ok, [username]} <- Repo.User.username(conn, decoded_key),
         :ok <- Repo.Message.new_message(conn, [type, username, conversation, datetime, message]) do
      {:reply, {:ok, "#{type}|#{username}|#{conversation}|#{datetime}|#{message}\n"}, conn}
    else
      {:error, reason} ->
        Logger.error(reason)
        {:reply, {:error, reason}, conn}
    end
  end

  def start_link([name, opts]) do
    GenServer.start_link(__MODULE__, name, opts)
  end

  def query_raw(type, query_string, values) do
    GenServer.call(__MODULE__, {type, query_string, values})
  end

  def push_message([type, key, conversation, datetime, message]) do
    GenServer.call(__MODULE__, [type, key, conversation, datetime, message])
  end

  def query(request) do
    GenServer.call(__MODULE__, request)
  end
end
