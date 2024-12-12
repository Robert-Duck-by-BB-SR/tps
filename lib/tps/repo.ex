defmodule TPS.Repo do
  use GenServer
  require Exqlite

  @impl true
  def init(db_name) do
    Exqlite.Sqlite3.open(db_name)
  end

  @impl true
  def handle_cast({:exec, query, values}, conn) do
    {:ok, statement} =
      conn
      |> Exqlite.Sqlite3.prepare(query)

    :ok = Exqlite.Sqlite3.bind(conn, statement, values)

    Exqlite.Sqlite3.step(conn, statement)
    Exqlite.Sqlite3.release(conn, statement)
    {:noreply, conn}
  end

  @impl true
  def handle_call({:select, query, values}, _from, conn) do
    {:ok, statement} =
      conn
      |> Exqlite.Sqlite3.prepare(query)

    :ok = Exqlite.Sqlite3.bind(conn, statement, values)

    case Exqlite.Sqlite3.multi_step(conn, statement) do
      {:busy} ->
        Exqlite.Sqlite3.release(conn, statement)
        {:stop, :busy, "database is busy", conn}

      {:error, reason} ->
        Exqlite.Sqlite3.release(conn, statement)
        {:stop, :error, reason, conn}

      {_, result} ->
        Exqlite.Sqlite3.release(conn, statement)
        {:reply, result, conn}
    end
  end

  @impl true
  def handle_call({:select, query, values}, _from, conn) do
    {:ok, statement} =
      conn
      |> Exqlite.Sqlite3.prepare(query)

    :ok = Exqlite.Sqlite3.bind(conn, statement, values)

    case Exqlite.Sqlite3.step(conn, statement) do
      {:busy} ->
        Exqlite.Sqlite3.release(conn, statement)
        {:stop, :busy, "database is busy", conn}

      {:error, reason} ->
        Exqlite.Sqlite3.release(conn, statement)
        {:stop, :error, reason, conn}

      {result} ->
        Exqlite.Sqlite3.release(conn, statement)
        {:reply, result, conn}
    end
  end

  def start_link(opts) do
    GenServer.start_link(__MODULE__, :ok, opts)
  end

  def collect(type, query, values) do
    GenServer.call(__MODULE__, {type, query, values})
  end

  def exec(query, values) do
    GenServer.cast(__MODULE__, {:exec, query, values})
  end
end
