defmodule TPS.NinjaHash do
  def encode(left, right) do
    lbytes = :binary.bin_to_list(left)
    rbytes = :binary.bin_to_list(right)

    Enum.zip(lbytes, rbytes)
    |> Enum.map(fn {l, r} -> rem(l - r, 256) end)
  end

  def decode(left, right) do
    rbytes = :binary.bin_to_list(right)

    Enum.zip(left, rbytes)
    |> Enum.map(fn {l, r} -> rem(l + r, 256) end)
  end
end
