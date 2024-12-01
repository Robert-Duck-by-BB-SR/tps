defmodule TpsTest do
  use ExUnit.Case
  doctest Tps

  test "greets the world" do
    assert Tps.hello() == :world
  end
end
