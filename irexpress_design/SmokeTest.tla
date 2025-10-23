---- MODULE SmokeTest ----
EXTENDS Naturals

VARIABLES x

Init == x = 0
Next == x < 5 /\ x' = x + 1

Spec == Init /\ [][Next]_x

====