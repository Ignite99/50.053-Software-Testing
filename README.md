# 50.053-Software-Testing

### Usage

Makefile below for Windows users:
```
CC=g++
CFLAGS=-c -Wall
LDFLAGS=
SOURCES=./src/main.cpp ./src/fuzzer/fuzzer.cpp # Add your cpp files as you see fit
OBJECTS=$(SOURCES:.cpp=.o)
EXECUTABLE=Software_Testing_Project.exe

all: $(EXECUTABLE)

$(EXECUTABLE): $(OBJECTS)
	$(CC) $(LDFLAGS) $(OBJECTS) -o $@

.cpp.o:
	$(CC) $(CFLAGS) $< -o $@

clean:
	del .\src\*.o
	del .\src\fuzzer\*.o
	del Software_Testing_Project.exe
```

Makefile for Linux/MACOS users (note I have not tested this in wsl or mac environments):
```
CC=g++
CFLAGS=-c -Wall
LDFLAGS=
SOURCES=main.cpp
OBJECTS=$(SOURCES:.cpp=.o)
EXECUTABLE=Software_Testing_Project.exe

all: $(SOURCES) $(EXECUTABLE)

$(EXECUTABLE): $(OBJECTS)
	$(CC) $(LDFLAGS) $(OBJECTS) -o $@

.cpp.o:
	$(CC) $(CFLAGS) $< -o $@

clean:
	rm -rf $(OBJECTS) $(EXECUTABLE)
```