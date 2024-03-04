# 50.053-Software-Testing

### Usage

Makefile below for Windows users:
```
CC=g++
CFLAGS=-c -Wall
LDFLAGS=
SOURCES=./src/main.cpp ./src/fuzzer/fuzzer.cpp ./src/BLE_Zephyr/ble_zephyr.cpp ./src/CoAP_Protocol/coap_protocol.cpp ./src/Django_Web/django_web.cpp # Add your cpp files as you see fit
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
SOURCES=./src/main.cpp ./src/fuzzer/fuzzer.cpp ./src/BLE_Zephyr/ble_zephyr.cpp ./src/CoAP_Protocol/coap_protocol.cpp ./src/Django_Web/django_web.cpp
OBJECTS=$(SOURCES:.cpp=.o)
EXECUTABLE=Software_Testing_Project.exe

all: $(EXECUTABLE)

$(EXECUTABLE): $(OBJECTS)
	$(CC) $(LDFLAGS) $(OBJECTS) -o $@

.cpp.o:
	$(CC) $(CFLAGS) $< -o $@

clean:
	rm -rf ./src/*.o
	rm -rf ./src/fuzzer/*.o
	rm -rf ./src/BLE_Zephyr/*.o
	rm -rf ./src/CoAP_Protocol/*.o
	rm -rf ./src/Django_Web/*.o
	rm -rf Software_Testing_Project.exe
```

### Executing file

These are the arguments that need to be put in:
```
./<compiled cpp file> <project type> <url> <request_type> <input_file_path>
```
For example:
```
./Software_Testing_Project.exe django test get ./src/inputs/test1.txt
```