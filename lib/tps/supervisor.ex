defmodule TPS.Supervisor do
  use Supervisor

  def start_link(opts) do
    Supervisor.start_link(__MODULE__, :ok, opts)
  end

  @impl true
  def init(:ok) do
    children = [
      {TPS.Chat, name: TPS.Chat},
      {TPS.Repo, "testing.db"}
    ]

    Supervisor.init(children, strategy: :one_for_one)
  end
end
