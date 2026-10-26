# Empiryczne wskaźniki niezawodności

Program wyznaczający empiryczne wskaźniki niezawodności F\*(t), R\*(t), f\*(t), λ\*(t)
oraz E\*T na podstawie próby losowej lub ręcznie wprowadzonych danych.

Zadanie laboratoryjne z przedmiotu *Niezawodność systemów komputerowych*
(Politechnika Świętokrzyska, Wydział Elektrotechniki, Automatyki i Informatyki, grupa 2ID21A).
Autorzy: Filip Stępień, Bartłomiej Karkoszka, Rafał Grot.

## Uruchomienie

Program działa w przeglądarce pod adresem
**https://filip-stepien.github.io/niezawodnosc-systemow-komputerowych/**.

Można go też uruchomić bez dostępu do internetu: to jeden plik HTML, który nie wymaga
instalacji. Pobierz `empirical-reliability.html` z [najnowszego wydania](https://github.com/filip-stepien/niezawodnosc-systemow-komputerowych/releases/latest)
i otwórz go w przeglądarce.

Strona i wydanie aktualizują się automatycznie po każdej zmianie na gałęzi `main`.

## Budowanie

Potrzebny jest Node.js 24 (albo `nix develop`).

```
npm install
npm run dev     # podgląd na żywo podczas pracy
npm run build   # gotowy plik: dist/index.html
```

## Sprawozdanie

**[Sprawozdanie z zadania laboratoryjnego](report/report.pdf)** - opis teoretyczny,
algorytm działania, kod programu, wyniki testowania i instrukcja obsługi.
