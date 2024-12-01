
defmodule EchoServer do
  def init(options) do
    {:ok, options}
  end

  def handle_in({"ping", [opcode: :text]}, state) do
    {:reply, :ok, {:text, "pong"}, state}
  end

  def handle_in({message, [opcode: :text]}, state) do
    IO.puts(message)
    {:reply, :ok, {:text, "message received"}, state}
  end

  def terminate(:timeout, state) do
    {:ok, state}
  end
end

require Logger
webserver = {Bandit, plug: Router, scheme: :http, port: 6969}
{:ok, _} = Supervisor.start_link([webserver], strategy: :one_for_one)
Logger.info("Plug now running on localhost:6969")
Process.sleep(:infinity)
