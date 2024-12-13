defmodule TPS.Repo do
  require Logger
  use GenServer
  require Exqlite

  @impl true
  def init(db_name) do
    Exqlite.Sqlite3.open(db_name)
  end

  @impl true
  def handle_cast({:exec, query, values}, conn) do
    Logger.info("hello????")

    IO.inspect("waht?")
    st = Exqlite.Sqlite3.prepare(conn, query)

    Logger.info(st)
    {:ok, statement} = st

    :ok = Exqlite.Sqlite3.bind(statement, values)
    Logger.info("yes we are good")

    result = Exqlite.Sqlite3.step(conn, statement)
    Logger.info(result)
    Exqlite.Sqlite3.release(conn, statement)
    {:noreply, conn}
  end

  @impl true
  def handle_cast(bro, what) do
    Logger.info(bro)
    Logger.info(what)
    {:noreply, what}
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

  def collect(type, query, values) do
    GenServer.call(__MODULE__, {type, query, values})
  end

  def exec(query, values) do
    Logger.info(query)
    Logger.info(values)
    GenServer.cast(__MODULE__, {:exec, query, values})
  end
end
