# Použití

## Výstup č.1

- Návod na spuštění na windows `./dist/castle.exe -c ./examples/bemdas.cst`
- Návod na spuštění na linux-based systémech `./dist/castle -c ./examples/bemdas.cst`
- Očekávaný výstup:

```
Loaded file contents:
6/2*(1+2) + 3
Label: 6, Type: LT_NUMBER
Label: /, Type: LT_DIVIDE
Label: 2, Type: LT_NUMBER
Label: *, Type: LT_MULTIPLY
Label: (, Type: LT_LPAREN
Label: 1, Type: LT_NUMBER
Label: +, Type: LT_PLUS
Label: 2, Type: LT_NUMBER
Label: ), Type: LT_RPAREN
Label: +, Type: LT_PLUS
Label: 3, Type: LT_NUMBER
```

- Co se děje:
  - Compiler pomocí flagu `-c` načte přepřipravený soubor `./examples/bemdas.cst`
  - Compiler zpracuje obsah souboru na jednotlivé lexémy
  - Compiler vypíše text jednotlivých lexémů a jejich typ

## Výstup č.2

- Návod na spuštění na windows `./dist/castle.exe -c ./examples/bemdas.cst`
- Návod na spuštění na linux-based systémech `./dist/castle -c ./examples/bemdas.cst`
- Očekávaný výstup:

```
Loaded file contents:
//Abrakadabra
6/2*(1+2) + 3
Label: 6, Type: LT_NUMBER
Label: /, Type: LT_DIVIDE
Label: 2, Type: LT_NUMBER
Label: *, Type: LT_MULTIPLY
Label: (, Type: LT_LPAREN
Label: 1, Type: LT_NUMBER
Label: +, Type: LT_PLUS
Label: 2, Type: LT_NUMBER
Label: ), Type: LT_RPAREN
Label: +, Type: LT_PLUS
Label: 3, Type: LT_NUMBER
Label: , Type: LT_END
//Abrakadabra
6/2*(1+2) + 3
[ ADD ]
≫ [ MULTIPLY ]
≫ ≫ [ DIVIDE ]
≫ ≫ ≫ [ VALUE ]
≫ ≫ ≫ ≫ 6
≫ ≫ ≫ [ VALUE ]
≫ ≫ ≫ ≫ 2
≫ ≫ [ GROUP ]
≫ ≫ ≫ [ ADD ]
≫ ≫ ≫ ≫ [ VALUE ]
≫ ≫ ≫ ≫ ≫ 1
≫ ≫ ≫ ≫ [ VALUE ]
≫ ≫ ≫ ≫ ≫ 2
≫ [ VALUE ]
≫ ≫ 3

```

- Co se děje:
  - Compiler pomocí flagu `-c` načte přepřipravený soubor `./examples/bemdas.cst`
  - Compiler zpracuje obsah souboru na jednotlivé lexémy
  - Compiler vypíše text jednotlivých lexémů a jejich typ
  - Compiler vynechá komentáře a logicky zpracuje výraz a vypíše jeho strukturu

## Výstup č.3

- Návod na spuštění na windows `./dist/castle.exe -c ./examples/bemdas.cst`
- Návod na spuštění na linux-based systémech `./dist/castle -c ./examples/bemdas.cst`

- Program vypisuje jeho průběh při schvalování a zpracovávání vstupu a počet tvrzení v vstupním textu

- Co probíhá v programu:

  - Načte text
  - Zpracuje text na lexémy
  - Zpracuje lexémy do abstraktní syntaktické stromové struktury pomocí parseru dle dané gramatiky
  - Program spadne, když objeví syntaktickou chybu
  - Program umí zpracovávat výrazy podle priority početních operací

    1. číslo, desetinné číslo, pravdivost (true, false), závorky, volání funkce, identifikátor
    2. operátor !, -, nebo + (převrácení hodnoty)
    3. násobení, dělení
    4. sčítání, odčítání,
    5. porovnávání >=, <=, >, <
    6. porovnávání ==, !=
    7. logická operace AND,
    8. logická operace OR,
    9. výraz
