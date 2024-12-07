defmodule TPSTest do
  use ExUnit.Case
  doctest TPS

  test "greets the world" do
    assert TPS.hello() == :world
  end
end
