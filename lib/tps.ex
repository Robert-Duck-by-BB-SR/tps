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
      :gen_tcp.listen(port, [
        :binary,
        packet: :line,
        active: false,
        reuseaddr: true
      ])

    Logger.info("Acception connections on port #{port}")
    loop_acceptor(socket)
  end

  defp loop_acceptor(socket) do
    {:ok, client} = :gen_tcp.accept(socket)
    Logger.info("accepting connection to")

    Task.start_link(fn -> read_line(client) end)
    loop_acceptor(socket)
  end

  defp read_line(socket) do
    case :gen_tcp.recv(socket, 0) do
      {:ok, data} ->
        Logger.info(data)
        message = String.trim(data)

        case Chat.Message.parse_incoming(message) do
          {:message, m} ->
            Chat.push_message(:push, m)

          {:req, r} ->
            "#{Chat.request(r)}\n" |> Chat.write_line(socket)

          {:connect, username} ->
            Chat.push_message(:connect, {username, socket})

          {:error, reason} ->
            Chat.write_line("BAD REQUEST: #{reason}", socket)
        end

        read_line(socket)

      {:error, :closed} ->
        Chat.remove_connection(socket)
        :ok
    end
  end
end
