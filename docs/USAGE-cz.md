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