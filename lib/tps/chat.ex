defmodule TPS.Chat do
  alias Tps.Chat.Message
  require Logger
  use GenServer

  @impl true
  def init(:ok) do
    Logger.info("init?")
    {:ok, []}
  end

  @impl true
  def handle_cast({:push, message}, clients) do
    m = Message.parse_message(message)

    Logger.warning("#{inspect(m)}")

    clients
    |> Enum.each(fn socket ->
      write_line("#{m.key}: #{m.message} -> #{Time.utc_now()}\n", socket)
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

  def start_link(opts) do
    Logger.info("start link?")
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

  defp write_line(line, socket) do
    :gen_tcp.send(socket, line)
  end
end
