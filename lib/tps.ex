defmodule TPS do
  alias TPS.Chat
  alias TPS.Supervisor
  use Application
  require Logger

  def start(_type, _args) do
    Supervisor.start_link([])
    accept(6969)
  end

  def accept(port) do
    {:ok, socket} =
      :gen_tcp.listen(port, [:binary, packet: :line, active: false, reuseaddr: true])

    Logger.info("Acception connections on port #{port}")
    loop_acceptor(socket)
  end

  defp loop_acceptor(socket) do
    {:ok, client} = :gen_tcp.accept(socket)
    Logger.info("accepting connection to")

    Task.start_link(fn ->
      Chat.push_message(:connect, client)
      read_line(client)
    end)

    loop_acceptor(socket)
  end

  defp read_line(socket) do
    case :gen_tcp.recv(socket, 0) do
      {:ok, data} ->
        Logger.info("printing lcients")
        IO.inspect(Chat.print_clients())
        message = String.trim(data)
        Chat.push_message(:push, message)
        read_line(socket)

      {:error, :close} ->
        Chat.remove_connection(socket)
        :ok
    end
  end
end
