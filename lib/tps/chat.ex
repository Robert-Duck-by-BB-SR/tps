defmodule TPS.Chat do
  alias TPS.Repo
  alias TPS.Chat.Message
  require Logger
  use GenServer

  @impl true
  def init(:ok) do
    {:ok, []}
  end

  @impl true
  def handle_cast({:push, message}, clients) do
    m = Message.parse_message(message)

    Logger.warning("#{inspect(m)}")

    datetime = Time.utc_now()

    {:ok, response} =
      Repo.push_message([m.type, m.key, m.convo, datetime, m.message])

    Logger.warning(response)

    clients
    |> Enum.each(fn socket ->
      write_line(response, socket)
    end)

    {:noreply, clients}
  end

  @impl true
  def handle_cast({:connect, socket}, clients) do
    Logger.info("connect")

    clients
    |> Enum.each(fn socket -> write_line("connected!\n", socket) end)

    {:noreply, [socket | clients]}
  end

  @impl true
  def handle_cast({:remove, socket}, clients) do
    list = Enum.filter(clients, &(&1 !== socket))
    {:noreply, list}
  end

  @impl true
  def handle_call(:get, _from, clients) do
    {:reply, :ok, clients}
  end

  @impl true
  def handle_call({:request, message}, _from, clients) do
    req = Message.parse_request(message)
    {:ok, result} = Repo.query(req)
    Logger.warning(result)
    {:reply, result, clients}
  end

  def start_link(opts) do
    GenServer.start_link(__MODULE__, :ok, opts)
  end

  def print_clients() do
    GenServer.call(__MODULE__, :get)
  end

  def push_message(type, opts) do
    GenServer.cast(__MODULE__, {type, opts})
  end

  def remove_connection(socket) do
    GenServer.cast(__MODULE__, {:remove, socket})
  end

  def request(message) do
    GenServer.call(__MODULE__, {:request, message})
  end

  def write_line(line, socket) do
    :gen_tcp.send(socket, line)
  end
end
