function Register()
  return "48 89 5C 24 10 48 89 74 24 18 55 57 41 54 41 56 41 57 48 8D AC 24 30 FE FF FF 48 81 EC D0 02 00 00 48 8B ? ? ? ? ? 48 33 C4 48 89 85 C0 01 00 00 48 8B 31 48 8B D9 4C 8B 61 08 44 8B 79 18 F7 86 D4 00 00 00 80 00 00 10 74 7D 4C 8B 71 28 48 8D 79 28 4D 85 F6 74 2E 83 79 10 00 75 06 83 79 14 00 74 64 80 79 21 00 75 1C 4C 8B 41 10 45 8B CF 48 8B CE C6 44 24 20 00 49 8B D4 E8 ? ? ? ? 4C 3B F0 75 42 41 8B C7 C1 E8 12 F6 D0 A8 01 75 32 48 8D 44 24 40 C6 44 24 40 00 48 89 44 24 50"
end

function OnMatchFound(MatchAddress)
  return MatchAddress
end