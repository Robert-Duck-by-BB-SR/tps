defmodule TPS.Repo do
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

  def start_link([name, opts]) do
    GenServer.start_link(__MODULE__, name, opts)
  end

  def query(type, query_string, values) do
    GenServer.call(__MODULE__, {type, query_string, values})
  end
end
